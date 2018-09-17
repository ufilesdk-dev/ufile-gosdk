package ufsdk

import (
	"encoding/json"
	"fmt"
)

//DomainSet BucketDataSet 里面的 Domain 字段
type DomainSet struct {
	Src       []string `json:"Src,omitempty"`
	Cdn       []string `json:"Cdn,omitempty"`
	CustomSrc []string `json:"CustomSrc,omitempty"`
	CustomCdn []string `json:"CustomCdn,omitempty"`
}

//BucketDataSet BucketListResponse 里面的 DataSet 字段
type BucketDataSet struct {
	BucketName    string    `json:"BucketName,omitempty"`
	BucketID      string    `json:"BucketId,omitempty"`
	Domain        DomainSet `json:"Domain,omitempty"`
	Type          string    `json:"Type,omitempty"`
	CreateTime    int       `json:"CreateTime,omitempty"`
	ModifyTime    int       `json:"ModifyTime,omitempty"`
	CdnDomainID   []string  `json:"CdnDomainId,omitempty"`
	Biz           string    `json:"Biz,omitempty"`
	Region        string    `json:"Region,omitempty"`
	HasUserDomain int       `json:"HasUserDomain,omitempty"`
}

//BucketListResponse DescribeBucket 返回的数据填充。
type BucketListResponse struct {
	RetCode int             `json:"RetCode,omitempty"`
	Action  string          `json:"Action,omitempty"`
	Message string          `json:"Message,omitempty"`
	DataSet []BucketDataSet `json:"DataSet,omitempty"`
}

//AdaptedAPI 用来指示 BucketListResponse 可以用作哪些接口的返回 body 序列化。
//这里是用来序列化 DescribeBucket 返回 body。
func (b *BucketListResponse) AdaptedAPI() string {
	return "DescribeBucket"
}

//String 输出 BucketListResponse 里面的字段详情。
func (b *BucketListResponse) String() string {
	bytes, err := json.MarshalIndent(b, "", " ")
	if err == nil {
		return fmt.Sprintf("%s\n", bytes)
	}
	return ""
}

//FileDataSet  FileListResponse 里面的 DataSet 字段。
type FileDataSet struct {
	BucketName  string `json:"BucketName,omitempty"`
	FileName    string `json:"FileName,omitempty"`
	Hash        string `json:"Hash,omitempty"`
	MimeType    string `json:"MimeType,omitempty"`
	FirstObject string `json:"first_object,omitempty"`
	Size        int    `json:"Size,omitempty"`
	CreateTime  int    `json:"CreateTime,omitempty"`
	ModifyTime  int    `json:"ModifyTime,omitempty"`
}

//FileListResponse PrefixFileList 返回的数据填充。
type FileListResponse struct {
	BucketName string        `json:"BucketName,omitempty"`
	BucketID   string        `json:"BucketId,omitempty"`
	NextMarker string        `json:"NextMarker,omitempty"`
	DataSet    []FileDataSet `json:"DataSet,omitempty"`
}

//AdaptedAPI 用来指示 FileListResponse 可以用作哪些接口的返回 body 序列化。
//这里是用来序列化 PrefixFileList 接口返回 body。
func (f *FileListResponse) AdaptedAPI() string {
	return "PrefixFileList"
}

//String 输出 FileListResponse 里面的字段详情。
func (f *FileListResponse) String() string {
	bytes, err := json.MarshalIndent(f, "", " ")
	if err == nil {
		return fmt.Sprintf("%s\n", bytes)
	}
	return ""
}

//BucketResponse 管理 Bucket 返回的数据填充。
type BucketResponse struct {
	RetCode    int    `json:"RetCode,omitempty"`
	Action     string `json:"Action,omitempty"`
	BucketName string `json:"BucketName,omitempty"`
	BucketID   string `json:"BucketId,omitempty"`
	Message    string `json:"Message,omitempty"`
}

//AdaptedAPI 用来指示 BucketResponse 可以用作哪些接口的返回 body 序列化。
//这里是用来序列化 CreateBucket,DeleteBucket,UpdateBucket,DescribeBucket 接口返回 body。
func (b *BucketResponse) AdaptedAPI() string {
	return "CreateBucket,DeleteBucket,UpdateBucket,DescribeBucket"
}

//String 输出 BucketResponse 里面的字段详情。
func (b *BucketResponse) String() string {
	bytes, err := json.MarshalIndent(b, "", " ")
	if err == nil {
		return fmt.Sprintf("%s\n", bytes)
	}
	return ""
}

//Adapter 用来限制 MarshalResult 里面的第二个参数类型。在这里，我们只能传上面带 Response 结尾的结构体。
type Adapter interface {
	AdaptedAPI() string
	String() string
}

//MarshalResult 序列化一些接口的返回值。目前支持 CreateBucket,DeleteBucket,UpdateBucket,DescribeBucket,PrefixFileList,DescribeBucket 返回值的序列化。
func MarshalResult(u *UFileRequest, adapter Adapter) error {
	return json.Unmarshal(u.LastResponseBody, adapter)
}
