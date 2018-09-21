package ufsdk

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

//MultipartState 用于保存分片上传的中间状态
type MultipartState struct {
	BlkSize  int //服务器返回的分片大小
	uploadID string
	keyName  string
	etags    []string
	mux      sync.Mutex
}

//UnmarshalJSON custom unmarshal json
func (m *MultipartState) UnmarshalJSON(bytes []byte) error {
	tmp := struct {
		BlkSize  int    `json:"BlkSize"`
		UploadID string `json:"UploadId"`
	}{}
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	m.BlkSize = tmp.BlkSize
	m.uploadID = tmp.UploadID
	return nil
}

type uploadChan struct {
	etag string
	err  error
}

//MPut 分片上传一个文件，filePath 是本地文件所在的路径，内部会自动对文件进行分片上传，上传的方式是同步一片一片的上传。
//mimeType 如果为空的话，会调用 net/http 里面的 DetectContentType 进行检测。
//大于 100M 的文件推荐使用本接口上传。
func (u *UFileRequest) MPut(filePath, keyName, mimeType string) error {
	state, err := u.InitiateMultipartUpload(keyName, mimeType)
	if err != nil {
		return err
	}
	file, err := openFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	chunk := make([]byte, state.BlkSize)
	var pos int
	for {
		_, fileErr := file.Read(chunk)
		if fileErr == io.EOF {
			break
		}
		buf := bytes.NewBuffer(chunk)
		err := u.UploadPart(buf, state, pos)
		if err != nil {
			u.AbortMultipartUpload(state)
			return err
		}
		pos++
	}

	return u.FinishMultipartUpload(state)
}

//AsyncMPut 异步分片上传一个文件，filePath 是本地文件所在的路径，内部会自动对文件进行分片上传，上传的方式是使用异步的方式同时传多个分片的块。
//mimeType 如果为空的话，会调用 net/http 里面的 DetectContentType 进行检测。
//大于 100M 的文件推荐使用本接口上传。
func (u *UFileRequest) AsyncMPut(filePath, keyName, mimeType string) error {
	state, err := u.InitiateMultipartUpload(keyName, mimeType)
	if err != nil {
		return err
	}
	file, err := openFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fsize := getFileSize(file)

	chunkCount := int(fsize / int64(state.BlkSize))
	maxJobRunning := 10 //最多允许 10 个线程同时跑
	concurrentChan := make(chan error, maxJobRunning)
	for i := 0; i != maxJobRunning; i++ {
		concurrentChan <- nil
	}

	wg := &sync.WaitGroup{}
	for i := 0; i != chunkCount; i++ {
		uploadErr := <-concurrentChan //最初允许启动 10 个 goroutine，超出10个后，有分片返回才会开新的goroutine.
		if uploadErr != nil {
			err = uploadErr
			break // 中间如果出现错误立即停止继续上传
		}
		wg.Add(1)
		go func(pos int) {
			defer wg.Done()
			offset := int64(state.BlkSize * pos)
			chunk := make([]byte, state.BlkSize)
			file.ReadAt(chunk, offset)
			e := u.UploadPart(bytes.NewBuffer(chunk), state, pos)
			concurrentChan <- e //跑完一个 goroutine 后，发信号表示可以开启新的 goroutine。
		}(i)
	}
	wg.Wait()       //等待所有任务返回
	if err == nil { //再次检查剩余上传完的分片是否有错误
		for e := range concurrentChan {
			if e != nil {
				err = e
				break
			}
		}
	}
	close(concurrentChan)
	if err != nil {
		u.AbortMultipartUpload(state)
		return err
	}

	return u.FinishMultipartUpload(state)
}

//AbortMultipartUpload 取消分片上传，如果掉用 UploadPart 出现错误，可以调用本函数取消分片上传。
//state 参数是 InitiateMultipartUpload 返回的
func (u *UFileRequest) AbortMultipartUpload(state *MultipartState) error {
	query := &url.Values{}
	query.Add("uploadId", state.uploadID)
	reqURL := u.genFileURL(state.keyName) + "?" + query.Encode()

	req, err := http.NewRequest("DELETE", reqURL, nil)
	if err != nil {
		return err
	}
	authorization := u.Auth.Authorization("DELETE", u.BucketName, state.keyName, req.Header)
	req.Header.Add("authorization", authorization)
	return u.request(req)
}

//InitiateMultipartUpload 初始化分片上传，返回一个 state 用于后续的 UploadPart, FinishMultipartUpload, AbortMultipartUpload 的接口。
//keyName 表示传到 ufile 的文件名。
//mimeType 表示文件的 mimeType, 传空会调用 net/http 里面的 DetectContentType 进行检测。
//state 参数是 InitiateMultipartUpload 返回的
func (u *UFileRequest) InitiateMultipartUpload(keyName, mimeType string) (*MultipartState, error) {
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
	response := new(MultipartState)
	err = json.Unmarshal(u.LastResponseBody, response)
	if err != nil {
		return nil, err
	}
	response.keyName = keyName
	response.etags = make([]string, 0)

	return response, err
}

//UploadPart 上传一个分片，buf 就是分片数据，buf 的数据块大小必须为 state.BlkSize，否则会报错。
//pardNumber 表示第几个分片，从 0 开始。例如一个文件按 state.BlkSize 分为 5 块，那么分片分别是 0,1,2,3,4。
//state 参数是 InitiateMultipartUpload 返回的
func (u *UFileRequest) UploadPart(buf *bytes.Buffer, state *MultipartState, partNumber int) error {
	query := &url.Values{}
	query.Add("uploadId", state.uploadID)
	query.Add("partNumber", strconv.Itoa(partNumber))

	reqURL := u.genFileURL(state.keyName) + "?" + query.Encode()
	req, err := http.NewRequest("PUT", reqURL, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "")

	if u.verifyUploadMD5 {
		md5Str := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
		req.Header.Add("Content-MD5", md5Str)
	}

	authorization := u.Auth.Authorization("PUT", u.BucketName, state.keyName, req.Header)
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Length", strconv.Itoa(buf.Len()))

	err = u.request(req)
	if err != nil {
		return err
	}

	etag := u.LastResponseHeader.Get("Etag")
	state.mux.Lock()
	state.etags = append(state.etags, etag)
	state.mux.Unlock()
	return nil
}

//FinishMultipartUpload 完成分片上传。分片上传必须要调用的接口。
//state 参数是 InitiateMultipartUpload 返回的
func (u *UFileRequest) FinishMultipartUpload(state *MultipartState) error {
	query := &url.Values{}
	query.Add("uploadId", state.uploadID)
	reqURL := u.genFileURL(state.keyName) + "?" + query.Encode()
	etagsStr := strings.Join(state.etags, ",")

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(etagsStr))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	authorization := u.Auth.Authorization("POST", u.BucketName, state.keyName, req.Header)
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Length", strconv.Itoa(len(etagsStr)))

	return u.request(req)
}
