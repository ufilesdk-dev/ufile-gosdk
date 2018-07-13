package filemgr

import (
	"ufile-gosdk/common"
	"io/ioutil"
	"encoding/json"

	"io"
	"os"
	"net/url"
	"strconv"
	//"net/http"
	"net/http"
	"time"
	"fmt"

	"ufile-gosdk/uflog"
)

type MultipartUploadResponse struct {
	RetCode    int
	ErrMsg     string
	XSessionId string

	UploadId   string
	BlkSize    int
	Bucket     string
	Key        string
}

func InitiateMultipartUpload(bucketname string, filename string, publickey string, privatekey string, proxyhost string) (MultipartUploadResponse, error) {
	multipartuploadresponse := MultipartUploadResponse{}

	authorization := common.Authorization(bucketname, filename, publickey, privatekey, "POST", "", "", "", "")

	if len(proxyhost) <= 4 {
		return multipartuploadresponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	res, err := common.FileApiMultipartUploadInit(bucketname, suffix, filename, authorization)
	if err != nil {
		return multipartuploadresponse, err
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return multipartuploadresponse, err
	}

	if res.StatusCode != 200 && res.StatusCode != 204 {
		multipartuploadresponse.XSessionId = res.Header.Get("X-Sessionid")

		uflog.ERROR(string(result))

		err = json.Unmarshal(result, &multipartuploadresponse)
		if err != nil {
			return multipartuploadresponse, err
		}

		uflog.ERROR(multipartuploadresponse)

		return multipartuploadresponse, fmt.Errorf(multipartuploadresponse.ErrMsg)
	}

	multipartuploadresponse.RetCode = 0

	err = json.Unmarshal(result, &multipartuploadresponse)
	if err != nil {
		return multipartuploadresponse, err
	}

	uflog.INFO(multipartuploadresponse)

	return multipartuploadresponse, nil
}

type UploadPartResponse struct {
	RetCode    int
	ErrMsg     string
	ETag       []string
	XSessionId string

	PartNumber int64
}

type SingleUploadPartResponse struct {
	RetCode    int
	ErrMsg     string
	ETag       string
	XSessionId string

	PartNumber int64
}

func UploadPartStream(bucketname string, f *os.File, filename string, uploadid string, partnumber int64, retrycount int, retryinterval int, mimetype string, publickey string, privatekey string, proxyhost string, blksize int64, connectionTimeout int) (UploadPartResponse, error) {
	uploadpartresponse := UploadPartResponse{}

	var err error

	if mimetype == "" {
		mimetype = "application/octec-stream"
	}

	authorization := common.Authorization(bucketname, filename, publickey, privatekey, "PUT", "", mimetype, "", "")

	if len(proxyhost) <= 4 {
		return uploadpartresponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	var fi os.FileInfo

	if fi, err = f.Stat(); err != nil {
		return uploadpartresponse, err
	}

	f.Seek(blksize * partnumber, 0)

	for ; fi.Size() > blksize * partnumber; partnumber++ {

		var cl int64
		if fi.Size() - blksize * partnumber < blksize {
			cl = fi.Size() - blksize * partnumber
		} else {
			cl = blksize
		}

		var rs io.Reader = f
		reader := io.LimitReader(rs, cl)

		uri := "?uploadId=" + uploadid + "&partNumber=" + strconv.FormatInt(partnumber, 10)
		uri = url.QueryEscape(uri)

		res, err := singleuploadpartstream(bucketname, suffix, filename, cl, reader, uri, mimetype, authorization, retrycount, retryinterval, connectionTimeout)
		//res, err := common.FileApiMultipartUploadStream(bucketname, suffix, filename, cl, reader, uri, mimetype, authorization)
		if err != nil {
			return uploadpartresponse, err
		}

		etag := res.Header.Get("Etag")

		result, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return uploadpartresponse, err
		}

		if res.StatusCode != 200 && res.StatusCode != 204  && res.StatusCode != 206 {

			xsessionid := res.Header.Get("X-Sessionid")

			singleuploadpartresponse := SingleUploadPartResponse{}

			err = json.Unmarshal(result, &singleuploadpartresponse)
			if err != nil {
				return uploadpartresponse, err
			}

			uflog.ERROR(singleuploadpartresponse)

			uploadpartresponse.RetCode = singleuploadpartresponse.RetCode
			uploadpartresponse.ErrMsg = singleuploadpartresponse.ErrMsg
			uploadpartresponse.XSessionId = xsessionid
			uploadpartresponse.PartNumber = partnumber

			return uploadpartresponse, fmt.Errorf("UploadPart Error")
		}

		uflog.INFO(string(result))

		singleuploadpartresponse := SingleUploadPartResponse{}

		err = json.Unmarshal(result, &singleuploadpartresponse)
		if err != nil {
			return uploadpartresponse, err
		}

		uploadpartresponse.ETag = append(uploadpartresponse.ETag, etag)
		uploadpartresponse.PartNumber = singleuploadpartresponse.PartNumber

		uflog.INFO(uploadpartresponse)
	}

	return uploadpartresponse, nil
}

func singleuploadpartstream(bucketname string, suffix string, filename string, cl int64, reader io.Reader, uri string, mimetype string, authorization string, retrycount int, retryinterval int, connectionTimeout int) (*http.Response, error) {
	var res *http.Response
	var err error

	for i := 0; i < retrycount; i++ {
		res, err = common.FileApiMultipartUploadStream(bucketname, suffix, filename, cl, reader, uri, mimetype, authorization, connectionTimeout)
		if err != nil {
			return res, err
		}

		if res.StatusCode == 200 || res.StatusCode == 204 || res.StatusCode == 206 {
			return res, err
		}

		time.Sleep(time.Second * time.Duration(retryinterval))
	}

	return res, err
}

func UploadPartFile(bucketname string, path string, filename string, uploadid string, partnumber int64, retrycount int, retryinterval int, publickey string, privatekey string, proxyhost string, blksize int64, connectionTimeout int) (UploadPartResponse, error) {
	uploadpartresponse := UploadPartResponse{}

	mimeType, err := common.MimeType(path)
	if err != nil {
		return uploadpartresponse, err
	}

	uflog.INFO(mimeType)

	f, err := os.Open(path)
	if err != nil {
		return uploadpartresponse, err
	}
	defer f.Close()

	uploadpartresponse, err = UploadPartStream(bucketname, f, filename, uploadid, partnumber, retrycount, retryinterval, mimeType, publickey, privatekey, proxyhost, blksize, connectionTimeout)

	return uploadpartresponse, err

}

func ResumeUploadPartStream(bucketname string, f *os.File, filename string, uploadid string, partnumber int64, etags []string, retrycount int, retryinterval int, mimeType string, publickey string, privatekey string, proxyhost string, blksize int64, connectionTimeout int) (UploadPartResponse, error) {
	uploadpartresponse := UploadPartResponse{}

	uploadpartstream, err := UploadPartStream(bucketname, f, filename, uploadid, partnumber, retrycount, retryinterval, mimeType, publickey, privatekey, proxyhost, blksize, connectionTimeout)

	uploadpartresponse.RetCode = uploadpartstream.RetCode
	uploadpartresponse.ETag = append(uploadpartresponse.ETag, etags...)
	uploadpartresponse.ETag = append(uploadpartresponse.ETag, uploadpartstream.ETag...)
	uploadpartresponse.PartNumber = uploadpartstream.PartNumber
	uploadpartresponse.ErrMsg = uploadpartstream.ErrMsg
	uploadpartresponse.XSessionId = uploadpartstream.XSessionId

	return uploadpartresponse, err
}

func ResumeUploadPartFile(bucketname string, path string, filename string, uploadid string, partnumber int64, etags []string, retrycount int, retryinterval int, publickey string, privatekey string, proxyhost string, blksize int64, connectionTimeout int) (UploadPartResponse, error) {
	uploadpartresponse := UploadPartResponse{}

	mimeType, err := common.MimeType(path)
	if err != nil {
		return uploadpartresponse, err
	}
	uflog.INFO(mimeType)

	f, err := os.Open(path)
	if err != nil {
		return uploadpartresponse, err
	}
	defer f.Close()

	uploadpartresponse, err = ResumeUploadPartStream(bucketname, f, filename, uploadid, partnumber, etags, retrycount, retryinterval, mimeType, publickey, privatekey, proxyhost, blksize, connectionTimeout)

	return uploadpartresponse, err

}

type FinishMultipartUploadResponse struct {
	RetCode    int
	ErrMsg     string
	Bucket     string
	Key        string
	FileSize   int
}

func FinishMultipartUpload(bucketname string, filename string, etags []string, uploadid string, publickey string, privatekey string, proxyhost string, connectionTimeout int) (FinishMultipartUploadResponse, error) {
	finishmultipartuploadresponse := FinishMultipartUploadResponse{}

	//mimeType, err := common.MimeType(path)
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//	return finishmultipartuploadresponse, err
	//}
	//fmt.Println(mimeType)

	//data := url.Values{}
	//data.Add("uploadId", uploadid)

	//singaturedata := common.Signature(privatekey, data)

	authorization := common.Authorization(bucketname, filename, publickey, privatekey, "POST", "", "text/plain", "", "")

	if len(proxyhost) <= 4 {
		return finishmultipartuploadresponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	res, err := common.FileApiFinishMultipartUpload(bucketname, suffix, filename, etags, uploadid, "text/plain", authorization, connectionTimeout)
	if err != nil {
		return finishmultipartuploadresponse, err
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return finishmultipartuploadresponse, err
	}

	err = json.Unmarshal(result, &finishmultipartuploadresponse)
	if err != nil {
		return finishmultipartuploadresponse, err
	}

	uflog.INFO(string(result))

	return finishmultipartuploadresponse, nil
}

type GetMultiUploadPartResponse struct {

}

func GetMultiUploadPart(bucketname string, uploadid string, publickey string, privatekey string, proxyhost string, connectionTimeout int) (GetMultiUploadPartResponse, error) {
	getmultiuploadpartresponse := GetMultiUploadPartResponse{}

	data := url.Values{}
	data.Add("uploadId", uploadid)

	singaturedata := common.Signature(privatekey, data)

	authorization := common.Authorization(bucketname, "王牌贱谍：格林斯比.HD1280高清中英双字.mp4", publickey, privatekey, "GET", "", "", "", "")

	if len(proxyhost) <= 4 {
		return getmultiuploadpartresponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	res, err := common.FileApiGetMultiUploadPart(bucketname, suffix, uploadid, singaturedata, authorization, connectionTimeout)
	if err != nil {
		return getmultiuploadpartresponse, err
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return getmultiuploadpartresponse, err
	}

	fmt.Println(string(result))
	uflog.INFO(string(result))

	return getmultiuploadpartresponse, nil

}
