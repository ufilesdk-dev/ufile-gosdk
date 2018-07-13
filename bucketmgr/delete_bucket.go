package bucketmgr

import (
	"net/url"
	"ufile-gosdk/common"
	"encoding/json"
	"ufile-gosdk/uflog"
)

type DeleteBucketResponse struct {
	Action     string
	Message    string
	RetCode    int
	BucketName string
	BucketId   string
}

func DeleteBucket(bucketname string, projectid string, publickey string, privatekey string) (DeleteBucketResponse, error) {
	deletebucketresponse := DeleteBucketResponse{}

	data := url.Values{}
	data.Add("Action", "DeleteBucket")
	data.Add("BucketName", bucketname)
	if projectid != "" {
		data.Add("ProjectId", projectid)
	}
	data.Add("PublicKey", publickey)

	singnaturedata := common.Signature(privatekey, data)

	res, err := common.ApiPost(singnaturedata)
	if err != nil {
		return deletebucketresponse, err
	}

	err = json.Unmarshal(res, &deletebucketresponse)
	if err != nil {
		return deletebucketresponse, err
	}

	uflog.INFO(deletebucketresponse)

	return deletebucketresponse, nil
}