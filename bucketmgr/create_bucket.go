package bucketmgr

import (
	"net/url"
	"ufile-gosdk/common"
	"encoding/json"

	"ufile-gosdk/uflog"
)

type CreateBucketResponse struct {
	Action     string
	Message    string
	RetCode    int
	BucketName string
	BucketId   string
}

func CreateBucket(bucketname string, buckettype string, region string, projectid string, publickey string, privatekey string) (CreateBucketResponse, error) {
	createbucketresponse := CreateBucketResponse{}

	data := url.Values{}
	data.Add("Action", "CreateBucket")
	data.Add("BucketName", bucketname)
	data.Add("Type", buckettype)
	data.Add("Region", region)
	if projectid != "" {
		data.Add("ProjectId", projectid)
	}
	data.Add("PublicKey", publickey)

	singnaturedata := common.Signature(privatekey, data)

	res, err := common.ApiPost(singnaturedata)
	if err != nil {
		return createbucketresponse, err
	}

	err = json.Unmarshal(res, &createbucketresponse)
	if err != nil {
		return createbucketresponse, err
	}

	uflog.INFO(createbucketresponse)

	return createbucketresponse, nil
}