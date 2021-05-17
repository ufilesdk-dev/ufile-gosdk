package ufsdk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

// SetMetaRequest  用于 设置元数据中的requestBody
type SetMetaRequest struct {
	Op    string `json:"op"`
	MetaK string `json:"metak"`
	MetaV string `json:"metav"`
}

// SetMimeType 设置key的mime类型
//key 对象的key
//mimeType 要设置的mime类型
func (u *UFileRequest) SetMimeType(key, mimeType string) (err error) {
	//构建SetMeta请求体
	setMimeTypeRequest := SetMetaRequest{
		Op:    "set",
		MetaK: "mimetype",
		MetaV: mimeType,
	}
	requestBody, err := json.Marshal(setMimeTypeRequest)
	if err != nil {
		return err
	}
	h := make(http.Header)
	h.Add("Content-Type", "application/json;charset=UTF-8")
	//生成签名
	authorization := u.Auth.Authorization("POST",u.BucketName,key,h)
	reqUrl := u.genFileURL(key)+"?opmeta"
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(requestBody)))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	return u.request(req)
}
