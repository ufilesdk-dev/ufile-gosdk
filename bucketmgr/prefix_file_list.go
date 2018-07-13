package bucketmgr

import (
	"net/url"
	"ufile-gosdk/common"
	"encoding/json"
	"strconv"
	"ufile-gosdk/uflog"
	"fmt"
)

type PrefixFileListResponse struct {
	Action     string
	RetCode    int
	ErrMsg     string
	BucketName string
	BucketId   string
	NextMarker string
	DataSet    []DataSet
}

type DataSet struct {
	BucketName string
	FileName   string
	Hash       string
	MimeType   string
	Size       int
	CreateTime int64
	ModifyTime int64
}

func PrefixFileList(bucketname string, prefix string, marker string, limit int, publickey string, privatekey string, proxyhost string) (PrefixFileListResponse, error) {

	prefixFileListResponse := PrefixFileListResponse{}

	data := url.Values{}
	data.Add("prefix", prefix)
	data.Add("marker", marker)
	if limit != 0 {
		data.Add("limit", strconv.Itoa(limit))
	}

	authorization := common.Authorization(bucketname, "", publickey, privatekey, "GET", "", "", "", "")

	// www.***
	if len(proxyhost) <= 4 {
		return prefixFileListResponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	res, err := common.FileListApiGet(bucketname, suffix, data, authorization)
	if err != nil {
		return prefixFileListResponse, err
	}

	err = json.Unmarshal(res, &prefixFileListResponse)
	if err != nil {
		return prefixFileListResponse, err
	}

	uflog.INFO(prefixFileListResponse)

	return prefixFileListResponse, fmt.Errorf(prefixFileListResponse.ErrMsg)
}
