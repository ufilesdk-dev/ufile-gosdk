package bucketmgr

import (
	"net/url"
	"ufile-gosdk/common"
	"encoding/json"
	"ufile-gosdk/uflog"
)

type UpdateBucketResponse struct {
	Action     string
	RetCode    int
	Message    string
	BucketName string
	BucketId   string
}

func UpdateBucket(bucketName string, bucketType string, projectId string, publicKey string, privateKey string) (UpdateBucketResponse, error) {
	updateBucketResponse := UpdateBucketResponse{}

	data := url.Values{}
	data.Add("Action", "UpdateBucket")
	data.Add("BucketName", bucketName)
	data.Add("Type", bucketType)
	if projectId != "" {
		data.Add("ProjectId", projectId)
	}
	data.Add("PublicKey", publicKey)

	singnaturedata := common.Signature(privateKey, data)

	res, err := common.ApiPost(singnaturedata)
	if err != nil {
		return updateBucketResponse, err
	}

	err = json.Unmarshal(res, &updateBucketResponse)
	if err != nil {
		return updateBucketResponse, err
	}

	uflog.INFO(updateBucketResponse)


	return updateBucketResponse, nil
}