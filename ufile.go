package ufsdk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//UploadHit 文件秒传，它的原理是计算出文件的 etag 值与远端服务器进行对比，如果文件存在就快速返回。
func (u *UFileRequest) UploadHit(filePath, keyName string) (err error) {
	file, err := openFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fsize := getFileSize(file)
	etag := calculateEtag(file)

	query := &url.Values{}
	query.Add("Hash", etag)
	query.Add("FileName", keyName)
	query.Add("FileSize", strconv.FormatInt(fsize, 10))
	reqURL := u.genFileURL("uploadhit") + "?" + query.Encode()
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return err
	}
	authorization := u.Auth.Authorization("POST", u.BucketName, keyName, req.Header)
	req.Header.Add("authorization", authorization)

	return u.request(req)
}

//PostFile 使用 HTTP Form 的方式上传一个文件。
//注意：使用本接口上传文件后，调用 UploadHit 接口会返回 404，因为经过 form 包装的文件，etag 值会不一样，所以会调用失败。
//mimeType 如果为空的话，会调用 net/http 里面的 DetectContentType 进行检测。
//小于 100M 的文件推荐使用本接口上传。
func (u *UFileRequest) PostFile(filePath, keyName, mimeType string) (err error) {
	file, err := openFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	h := make(http.Header)
	if mimeType == "" {
		mimeType = getMimeType(file)
	}
	h.Add("Content-Type", mimeType)
	authorization := u.Auth.Authorization("POST", u.BucketName, keyName, h)

	boundry := makeBoundry()
	body := makeFormBody(authorization, boundry, keyName, mimeType, file)
	//lastline 一定要写，否则后端解析不到。
	lastline := fmt.Sprintf("\r\n--%s--\r\n", boundry)
	body.Write([]byte(lastline))

	reqURL := u.genFileURL("")
	req, err := http.NewRequest("POST", reqURL, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundry)
	contentLength := body.Len()
	req.Header.Add("Content-Length", strconv.Itoa(contentLength))

	return u.request(req)
}

//PutFile 把文件直接放到 HTTP Body 里面上传，相对 PostFile 接口，这个要更简单，速度会更快（因为不用包装 form）。
//mimeType 如果为空的，会调用 net/http 里面的 DetectContentType 进行检测。
//小于 100M 的文件推荐使用本接口上传。
func (u *UFileRequest) PutFile(filePath, keyName, mimeType string) error {
	reqURL := u.genFileURL(keyName)
	file, err := openFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", reqURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Go-http-client/1.1")
	if mimeType == "" {
		mimeType = getMimeType(file)
	}
	req.Header.Add("Content-Type", mimeType)

	authorization := u.Auth.Authorization("PUT", u.BucketName, keyName, req.Header)
	req.Header.Add("authorization", authorization)
	fileSize := getFileSize(file)
	req.Header.Add("Content-Length", strconv.FormatInt(fileSize, 10))

	return u.request(req)
}

//DeleteFile 删除一个文件，如果删除成功 statuscode 会返回 204，否则会返回 404 表示文件不存在。
func (u *UFileRequest) DeleteFile(keyName string) error {
	reqURL := u.genFileURL(keyName)
	req, err := http.NewRequest("DELETE", reqURL, nil)
	if err != nil {
		return err
	}
	authorization := u.Auth.Authorization("DELETE", u.BucketName, keyName, req.Header)
	req.Header.Add("authorization", authorization)
	return u.request(req)
}

//HeadFile 获取一个文件的基本信息，返回的信息全在 header 里面。包含 mimeType, content-length（文件大小）, etag, Last-Modified:。
func (u *UFileRequest) HeadFile(keyName string) error {
	reqURL := u.genFileURL(keyName)
	req, err := http.NewRequest("HEAD", reqURL, nil)
	if err != nil {
		return err
	}
	authorization := u.Auth.Authorization("HEAD", u.BucketName, keyName, req.Header)
	req.Header.Add("authorization", authorization)
	return u.request(req)
}

//PrefixFileList 获取文件列表。
//prefix 表示匹配文件前缀。
//marker 标志字符串
//limit 列表数量限制，传 0 会默认设置为 20.
func (u *UFileRequest) PrefixFileList(prefix, marker string, limit int) error {
	query := &url.Values{}
	query.Add("prefix", prefix)
	query.Add("marker", marker)
	if limit == 0 {
		limit = 20
	}
	query.Add("limit", strconv.Itoa(limit))
	reqURL := u.genFileURL("") + "?list&" + query.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}

	authorization := u.Auth.Authorization("GET", u.BucketName, "", req.Header)
	req.Header.Add("authorization", authorization)

	return u.request(req)
}

//GetPublicURL 获取公有空间的文件下载 URL
func (u *UFileRequest) GetPublicURL(keyName string) string {
	return u.genFileURL(keyName)
}

//GetPrivateURL 获取私有空间的文件下载 URL。
//expiresDuation 表示下载链接的过期时间，从现在算起，24 * time.Hour 表示过期时间为一天。
func (u *UFileRequest) GetPrivateURL(keyName string, expiresDuation time.Duration) string {
	t := time.Now()
	t.Add(expiresDuation)
	expires := strconv.Itoa(int(t.Unix()))
	signature, publicKey := u.Auth.AuthorizationPrivateURL("GET", u.BucketName, keyName, expires, http.Header{})
	query := url.Values{}
	query.Add("UCloudPublicKey", publicKey)
	query.Add("Signature", signature)
	query.Add("Expires", expires)
	reqURL := u.genFileURL(keyName)
	return reqURL + "?" + query.Encode()
}

//Download 把文件下载到 HTTP Body 里面。
func (u *UFileRequest) Download(reqURL string) error {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}
	return u.request(req)
}

func (u *UFileRequest) genFileURL(keyName string) string {
	return fmt.Sprintf("http://%s.%s/%s", u.BucketName, u.Host, keyName)
}
