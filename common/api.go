package common

import (
	"io/ioutil"
	"net/http"
	"strings"
	"net/url"
	"os"
	"bufio"
	"strconv"
	"time"
	"io"
	"bytes"
	"ufile-gosdk/uflog"
	"fmt"
)

func ApiPost(data url.Values) ([]byte, error) {
	uflog.INFO(data)

	body := strings.NewReader(data.Encode())

	res, err := http.Post("http://api.ucloud.cn", "application/x-www-form-urlencoded;param=value", body)
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	uflog.INFO(string(result))

	return result, err
}

func FileListApiGet(bucketname string, suffix string, data url.Values, authorization string) ([]byte, error) {
	uflog.INFO(data)

	client := &http.Client{}

	applyreq, err := http.NewRequest("GET", "http://" + bucketname + suffix + "/?list&" + data.Encode(), nil)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)

	res, err := client.Do(applyreq)

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}
	uflog.INFO(string(result))

	return result, err
}

func FileApiPut(bucketname string, suffix string, path string, filename string, data url.Values, authorization string, contentType string) (*http.Response, error) {
	var fi os.FileInfo
	var buf io.Reader
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	} else {
		buf = bufio.NewReader(f)
		if fi, err = f.Stat(); err != nil {
			return nil, err
		}
	}
	defer f.Close()
	var fileSize = fi.Size()

	uflog.INFO(data)

	client := &http.Client{}

	applyreq, err := http.NewRequest("PUT", "http://" + bucketname + suffix + "/" + filename, buf)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)
	applyreq.Header.Add("Content-Length", strconv.Itoa(int(fileSize)))
	applyreq.Header.Add("Content-Type", contentType)

	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err
}

func FileApiGetPublic(bucketname string, suffix string, filename string, byterange string, authorization string, connectionTimeout int) (*http.Response, error) {
	//client := &http.Client{}

	client, err := NewTimeoutClient(time.Duration(connectionTimeout) * time.Second, -1, 300 * time.Second, false)
	if err != nil {
		return nil, err
	}

	applyreq, err := http.NewRequest("GET", "http://" + bucketname + suffix + "/" + filename, nil)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)
	applyreq.Header.Add("Range", byterange)
	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err
}

func FileApiGetPrivate(bucketname string, suffix string, filename string, byterange string, authorization string, data url.Values, connectionTimeout int) (*http.Response, error) {
	//client := &http.Client{}

	client, err := NewTimeoutClient(time.Duration(connectionTimeout) * time.Second, -1, 300 * time.Second, false)
	if err != nil {
		return nil, err
	}

	uflog.INFO(data)

	body := strings.NewReader(data.Encode())

	applyreq, err := http.NewRequest("GET", "http://" + bucketname + suffix + "/" + filename, body)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)
	applyreq.Header.Add("Range", byterange)
	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err
}

func FileApiUploadhit(bucketname string, suffix string, localfile string, key string, data url.Values, authorization string, contentType string, connectionTimeout int) (*http.Response, error) {
	var fi os.FileInfo
	f, err := os.Open(localfile)
	if err != nil {
		return nil, err
	} else {
		if fi, err = f.Stat(); err != nil {
			return nil, err
		}
	}
	defer f.Close()

	etag, err := ETag(localfile)
	if err != nil {
		return nil, err
	}
	uri := "?FileName=" + key + "&Hash=" + etag + "&FileSize=" + strconv.FormatInt(fi.Size(), 10)
	uri = url.QueryEscape(uri)

	uflog.INFO(data)


	//client := &http.Client{}

	applyreq, err := http.NewRequest("POST", "http://" + bucketname + suffix + "/uploadhit" + uri, nil)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)
	applyreq.Header.Add("Content-Length", strconv.Itoa(int(fi.Size())))
	applyreq.Header.Add("Content-Type", contentType)

	client, err := NewTimeoutClient(time.Duration(connectionTimeout) * time.Second, -1, 300 * time.Second, false)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err
}

func FileApiDelete(bucketname string, suffix string, filename string, authorization string) (*http.Response, error) {
	client := &http.Client{}

	applyreq, err := http.NewRequest("DELETE", "http://" + bucketname + suffix + "/" + filename, nil)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)

	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err

}

func FileApiMultipartUploadInit(bucketname string, suffix string, filename string, authorization string) (*http.Response, error) {
	client := &http.Client{}

	applyreq, err := http.NewRequest("POST", "http://" + bucketname + suffix + "/" + filename + "?uploads", nil)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)

	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err
}

func FileApiMultipartUploadStream(bucketname string, suffix string, key string, contentlength int64, reader io.Reader, uri string, mimetype string, authorization string, connectionTimeout int) (*http.Response, error) {

	applyreq, err := http.NewRequest("PUT", "http://" + bucketname + suffix + "/" + key + uri, reader)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)
	applyreq.Header.Add("Content-Length", strconv.Itoa(int(contentlength)))
	applyreq.Header.Add("Content-Type", mimetype)

	client, err := NewTimeoutClient(time.Duration(connectionTimeout) * time.Second, -1, 300 * time.Second, false)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err

}

func FileApiGetMultiUploadPart(bucketname string, suffix string, uploadid string, data url.Values, authorization string, connectionTimeout int) (*http.Response, error) {
	body := strings.NewReader(data.Encode())
	uflog.INFO(data.Encode())

	uri := "?muploadpart&uploadId=" + uploadid
	uri = url.QueryEscape(uri)

	//applyreq, err := http.NewRequest("GET", "http://" + bucketname + suffix + "/"+uri, body)
	applyreq, err := http.NewRequest("GET", "http://" + bucketname + suffix + "/muploadpart", body)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)
	//applyreq.Header.Add("Content-Length", strconv.Itoa(int(contentlength)))
	//applyreq.Header.Add("Content-Type", "application/octec-stream")

	client, err := NewTimeoutClient(time.Duration(connectionTimeout) * time.Second, -1, 300 * time.Second, false)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err
}

func FileApiFinishMultipartUpload(bucketname string, suffix string, key string, etags []string, uploadid string, contenttype string, authorization string, connectionTimeout int) (*http.Response, error) {
	//var fi os.FileInfo
	//f, err := os.Open(path)
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//	return nil, err
	//} else {
	//	if fi, err = f.Stat(); err != nil {
	//		fmt.Errorf(err.Error())
	//		return nil, err
	//	}
	//}
	//defer f.Close()

	buffer := new(bytes.Buffer)
	for i, etag := range etags {
		etag = strings.Trim(etag, "\"")
		if i != len(etags) - 1 {
			buffer.WriteString(fmt.Sprintf("%s,", etag))
		} else {
			buffer.WriteString(fmt.Sprintf("%s", etag))
		}
	}

	uri := "?uploadId=" + uploadid
	uri = url.QueryEscape(uri)

	applyreq, err := http.NewRequest("POST", "http://" + bucketname + suffix + "/" + key + uri, buffer)
	if err != nil {
		return nil, err
	}

	applyreq.Header.Add("Authorization", authorization)
	applyreq.Header.Add("Content-Length", strconv.Itoa(buffer.Len()))
	applyreq.Header.Add("Content-Type", contenttype)

	client, err := NewTimeoutClient(time.Duration(connectionTimeout) * time.Second, -1, 300 * time.Second, false)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(applyreq)
	if err != nil {
		return nil, err
	}

	uflog.INFO(res)

	return res, err

}