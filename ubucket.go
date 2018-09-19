package ufsdk

import (
	"net/http"
	"net/url"
	"strconv"
)

//CreateBucket 创建一个 bucket, bucketName 必须全部为小写字母，不能带符号和特殊字符。
//
//region 可以填如下字段：
//
//cn-bj 北京
//
//hk 香港
//
//cn-gd 广州
//
//cn-sh2 上海二
//
//dn-jakarta 雅加达
//
//us-ca 洛杉矶
//
//bucketType 可以填 public（公有空间） 和 private（私有空间）
//projectID bucket 所在的项目 ID，可为空。
func (u *UFileRequest) CreateBucket(bucketName, region, bucketType, projectID string) error {
	query := url.Values{}
	query.Add("Action", "CreateBucket")
	query.Add("BucketName", bucketName)
	query.Add("Type", bucketType)
	query.Add("Region", region)

	if projectID != "" {
		query.Add("ProjectId", projectID)
	}

	return u.bucketRequest(query)
}

//DeleteBucket 删除一个 bucket，如果成功，status code 会返回 204 no-content
func (u *UFileRequest) DeleteBucket(bucketName, projectID string) error {
	query := url.Values{}
	query.Add("Action", "DeleteBucket")
	query.Add("BucketName", bucketName)
	if projectID != "" {
		query.Add("ProjectId", projectID)
	}

	return u.bucketRequest(query)
}

//UpdateBucket 更新一个 bucket，你可以改 bucket 的类型（私有或公有）和 项目 ID。
//bucketType 填公有（public）或私有（private）。
//projectID 没有可以填空（""）。
func (u *UFileRequest) UpdateBucket(bucketName, bucketType, projectID string) error {
	query := url.Values{}
	query.Add("Action", "UpdateBucket")
	query.Add("BucketName", bucketName)
	query.Add("Type", bucketType)
	if projectID != "" {
		query.Add("ProjectId", projectID)
	}

	return u.bucketRequest(query)
}

//DescribeBucket 获取 bucket 的详细信息，如果 bucketName 为空，返回当前账号下所有的 bucket。
//limit 是限制返回的 bucket 列表数量。
//offset 是列表的偏移量，默认为 0。
//projectID 可为空。
func (u *UFileRequest) DescribeBucket(bucketName string, offset, limit int, projectID string) error {
	query := url.Values{}
	query.Add("Action", "DescribeBucket")
	if bucketName != "" {
		query.Add("BucketName", bucketName)
	}
	//offset default is 0
	query.Add("Offset", strconv.Itoa(offset))

	if limit == 0 {
		limit = 20
	}
	query.Add("Limit", strconv.Itoa(limit))

	if projectID != "" {
		query.Add("ProjectId", projectID)
	}
	return u.bucketRequest(query)
}

func (u *UFileRequest) bucketRequest(query url.Values) error {
	reqURL := u.genBucketURL(query)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := u.Client.Do(req)
	if err != nil {
		return err
	}
	return u.resposneParse(resp)
}

func (u *UFileRequest) genBucketURL(query url.Values) string {
	var scheme = "http://"
	return scheme + u.Host + "/?" + u.Auth.AuthorizationBucketMgr(query)
}
