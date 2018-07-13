package filemgr

import (
	"net/url"

	"ufile-gosdk/common"
	"encoding/json"
	"io/ioutil"
	"ufile-gosdk/uflog"
	"fmt"
)

type UploadHitResponse struct {
	RetCode int
	ErrMsg  string
}

func UploadHit(bucketname string, filename string, key string, publickey string, privatekey string, proxyhost string, connectionTimeout int) (UploadHitResponse, error) {
	uploadHitResponse := UploadHitResponse{}

	mimeType, err := common.MimeType(filename)
	if err != nil {
		return uploadHitResponse, err
	}

	data := url.Values{}

	authorization := common.Authorization(bucketname, key, publickey, privatekey, "POST", "", mimeType, "", "")

	if len(proxyhost) <= 4 {
		return uploadHitResponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	resp, err := common.FileApiUploadhit(bucketname, suffix, filename, key, data, authorization, mimeType, connectionTimeout)
	if err != nil {
		return uploadHitResponse, err
	}

	if resp.StatusCode == 200 {
		uflog.INFO("File[", key, "] Hit Success")

		uploadHitResponse.RetCode = 0
		return uploadHitResponse, nil
	}

	uflog.INFO("X-SessionId:", resp.Header.Get("X-SessionId"))

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return uploadHitResponse, err
	}

	err = json.Unmarshal(result, &uploadHitResponse)
	if err != nil {
		return uploadHitResponse, err
	}

	uflog.ERROR(uploadHitResponse)
	uploadHitResponse.RetCode = resp.StatusCode

	return uploadHitResponse, fmt.Errorf("Not Hit")
}