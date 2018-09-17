package ufsdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type shardingResponse struct {
	UploadID string `json:"UploadId"`
	BlkSize  int    `json:"BlkSize"`
}

type uploadChan struct {
	etag string
	err  error
}

//ShardingUpload 分片上传一个文件，filePath 是本地文件所在的路径，内部会自动对文件进行分片上传，上传的方式是同步一片一片的上传。
//mimeType 如果为空的话，会调用 net/http 里面的 DetectContentType 进行检测。
//返回的第一个值是 uploadID, 如果上传中间出现失败，可以调用 AbortShardingUpload 来取消分片上传（需要传入 uploadID）。
//大于 100M 的文件推荐使用本接口上传。
func (u *UFileRequest) ShardingUpload(filePath, keyName, mimeType string) (string, error) {
	prepare, err := u.initSharedingUpload(keyName, mimeType)
	if err != nil {
		return "", err
	}
	file, err := openFile(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	chunk := make([]byte, prepare.BlkSize)
	var tags []string
	var pos int
	for {
		bytesread, fileErr := file.Read(chunk)
		if fileErr == io.EOF {
			break
		}
		tag, err := u.chunkSharedingUpload(chunk[:bytesread], keyName, prepare.UploadID, pos)
		if err != nil {
			return prepare.UploadID, err
		}
		tags = append(tags, tag)
		pos++
	}

	err = u.finishSharedingUpload(prepare.UploadID, keyName, tags)
	return "", err
}

//AsyncShardingUpload 异步分片上传一个文件，filePath 是本地文件所在的路径，内部会自动对文件进行分片上传，上传的方式是使用异步的方式同时传多个分片的块。
//mimeType 如果为空的话，会调用 net/http 里面的 DetectContentType 进行检测。
//返回的第一个值是 uploadID, 如果上传中间出现失败，可以调用 AbortShardingUpload 来取消分片上传（需要传入 uploadID）。
//大于 100M 的文件推荐使用本接口上传。
func (u *UFileRequest) AsyncShardingUpload(filePath, keyName, mimeType string) (string, error) {
	prepare, err := u.initSharedingUpload(keyName, mimeType)
	if err != nil {
		return "", err
	}
	file, err := openFile(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	chunk := make([]byte, prepare.BlkSize)
	var tags []string
	var pos int
	var errOccur bool
	wg := &sync.WaitGroup{}
	resultChan := make(chan *uploadChan)
	for {
		bytesread, fileErr := file.Read(chunk)
		if fileErr == io.EOF {
			break
		}
		go func(data []byte, filePos int) {
			if errOccur == true {
				return
			}
			wg.Add(1)
			defer wg.Done()
			etag, chunkErr := u.chunkSharedingUpload(data, keyName, prepare.UploadID, filePos)
			resultChan <- &uploadChan{etag, chunkErr}
		}(chunk[:bytesread], pos)
		pos++
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for r := range resultChan {
		if r.err != nil {
			errOccur = true
			err = r.err
			break
		} else {
			tags = append(tags, r.etag)
		}
	}
	if err != nil {
		return prepare.UploadID, err
	}

	return prepare.UploadID, u.finishSharedingUpload(prepare.UploadID, keyName, tags)
}

//AbortShardingUpload 取消分片上传，如果调用 ShardingUpload 或 AsyncShardingUpload 出现错误，可以调用本函数取消分片上传。
//uploadID 就是 ShardingUpload 或 AsyncShardingUpload 返回的第一个参数。
func (u *UFileRequest) AbortShardingUpload(keyName, uploadID string) error {
	query := &url.Values{}
	query.Add("uploadId", uploadID)
	reqURL := u.genFileURL(keyName) + "?" + query.Encode()

	req, err := http.NewRequest("DELETE", reqURL, nil)
	if err != nil {
		return err
	}
	authorization := u.Auth.Authorization("DELETE", u.BucketName, keyName, req.Header)
	req.Header.Add("authorization", authorization)
	return u.request(req)
}

func (u *UFileRequest) initSharedingUpload(keyName, mimeType string) (*shardingResponse, error) {
	reqURL := u.genFileURL(keyName) + "?uploads"
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", mimeType)
	authorization := u.Auth.Authorization("POST", u.BucketName, keyName, req.Header)
	req.Header.Add("authorization", authorization)

	err = u.request(req)
	if err != nil {
		return nil, err
	}
	response := new(shardingResponse)
	err = json.Unmarshal(u.LastResponseBody, response)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (u *UFileRequest) chunkSharedingUpload(chunk []byte, keyName, uploadID string, partNumber int) (string, error) {
	query := &url.Values{}
	query.Add("uploadId", uploadID)
	query.Add("partNumber", strconv.Itoa(partNumber))

	reqURL := u.genFileURL(keyName) + "?" + query.Encode()
	req, err := http.NewRequest("PUT", reqURL, bytes.NewBuffer(chunk))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "")
	authorization := u.Auth.Authorization("PUT", u.BucketName, keyName, req.Header)
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Length", strconv.Itoa(len(chunk)))

	err = u.request(req)
	if err != nil {
		return "", err
	}

	etag := u.LastResponseHeader.Get("Etag")
	return strings.Trim(etag, "\""), nil
}

func (u *UFileRequest) finishSharedingUpload(uploadID, keyName string, etags []string) error {
	query := &url.Values{}
	query.Add("uploadId", uploadID)
	reqURL := u.genFileURL(keyName) + "?" + query.Encode()
	etagsStr := strings.Join(etags, ",")

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(etagsStr))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	authorization := u.Auth.Authorization("POST", u.BucketName, keyName, req.Header)
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Length", strconv.Itoa(len(etagsStr)))

	return u.request(req)
}
