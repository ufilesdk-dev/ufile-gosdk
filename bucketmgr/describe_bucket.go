package bucketmgr

import (
	"net/url"
	"ufile-gosdk/common"
	"encoding/json"
	"ufile-gosdk/uflog"
)

type DescribeBucketResponse struct {
	Action     string
	RetCode    int
	Message    string
	DataSet    []UFileBucket
}

type UFileBucket struct {
	BucketName    string
	BucketId      string
	Domain        UFileDomain
	CdnDomainId   []string
	Type          string
	CreateTime    int
	ModifyTime    int
	Biz           string
	Region        string
	Tag           string
	HasUserDomain int
}

type UFileDomain struct {
	Src       []string
	Cdn	      []string
	CustomSrc []string
	CustomCdn []string
}

func DescribeBucket(bucketName string, projectId string, publicKey string, privateKey string) (DescribeBucketResponse, error) {
	describeBucketResponse := DescribeBucketResponse{}

	data := url.Values{}
	data.Add("Action", "DescribeBucket")
	if bucketName != "" {
		data.Add("BucketName", bucketName)
	}
	if projectId != "" {
		data.Add("ProjectId", projectId)
	}
	data.Add("PublicKey", publicKey)

	singnaturedata := common.Signature(privateKey, data)

	res, err := common.ApiPost(singnaturedata)
	if err != nil {
		return describeBucketResponse, err
	}

	err = json.Unmarshal(res, &describeBucketResponse)
	if err != nil {
		return describeBucketResponse, err
	}

	uflog.INFO(describeBucketResponse)


	return describeBucketResponse, nil
}