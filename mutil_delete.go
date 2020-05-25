package ufsdk

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

//UFileFile 单个文件文件名
type UFileFile struct {
	Key string `xml:"Key"`
}

//MultiDeleteFileStruct 批量删除xml结构体
type MultiDeleteFileStruct struct {
	XMLName xml.Name    `xml:"Delete"`
	IsQuit  bool        `xml:"Quiet"`
	Objects []UFileFile `xml:"Object"`
}

//ObjectMultiDeleteResult  返回结构体
type ObjectMultiDeleteResult struct {
	XMLName xml.Name    `xml:"DeleteResult"`
	Objects []UFileFile `xml:"Deleted"`
}

//MultiDeleteFile 批量删除文件
//input 批量删除的文件，IsQuit 返回的格式，默认false
//简单模式（quiet）：ufile返回的消息体中只包含删除过程中出错的Object结果。如果所有删除都成功，则没有消息体。
//详细模式（verbose）：ufile返回的消息体中会包含所有删除Object的结果。默认采用详细模式。
//有效值：true（开启简单模式）、false（开启详细模式）
func (u *UFileRequest) MultiDeleteFile(input []UFileFile, isQuit bool) (res *ObjectMultiDeleteResult, err error) {
	//鉴权
	mimeType := "multipart/form-data"
	h := make(http.Header)
	h.Add("Content-Type", mimeType)

	authorization := u.Auth.Authorization("POST", u.BucketName, "xml", h)

	buffer, err := getXmlData(input, isQuit) //上传的文件名
	if err != nil {
		return nil, err
	}
	//内容体构造
	boundry := makeBoundry()
	body := makeFormBodyV2(authorization, boundry, "xml", mimeType, false, buffer)
	//lastline 一定要写，否则后端解析不到。
	lastline := fmt.Sprintf("\r\n--%s--\r\n", boundry)
	body.Write([]byte(lastline))

	url := u.MultiDelete()
	//发送http请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundry)
	contentLength := body.Len()
	req.Header.Add("Content-Length", strconv.Itoa(contentLength))
	defer req.Body.Close()

	for k, v := range u.RequestHeader {
		for i := 0; i < len(v); i++ {
			req.Header.Add(k, v[i])
		}
	}

	return u.FinishMultiDelete(req)
}

func getXmlData(input []UFileFile, isQuit bool) (io.Reader, error) {
	xmlFile := new(MultiDeleteFileStruct)
	xmlFile.IsQuit = isQuit
	xmlFile.Objects = input
	xmlData, err := xml.Marshal(xmlFile)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(xmlData)

	return buffer, nil
}

func makeFormBodyV2(authorization, boundry, keyName, mimeType string, verifyMD5 bool, file io.Reader) *bytes.Buffer {
	boundry = "--" + boundry + "\r\n"
	boundryBytes := []byte(boundry)
	body := new(bytes.Buffer)

	body.Write(boundryBytes)
	body.Write(makeFormField("Authorization", authorization))
	body.Write(boundryBytes)
	body.Write(makeFormField("Content-Type", mimeType))
	body.Write(boundryBytes)
	body.Write(makeFormField("FileName", keyName))
	body.Write(boundryBytes)

	if verifyMD5 {
		h := md5.New()
		io.Copy(h, file)
		md5Str := fmt.Sprintf("%x", h.Sum(nil))
		body.Write(makeFormField("Content-MD5", md5Str))
		body.Write(boundryBytes)
	}

	addtionalStr := fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n", keyName)
	addtionalStr += fmt.Sprintf("Content-Type: %s\r\n\r\n", mimeType)
	body.Write([]byte(addtionalStr))
	body.ReadFrom(file)
	body.Write([]byte("\r\n"))
	body.Write(boundryBytes)

	return body
}

func (u *UFileRequest) FinishMultiDelete(req *http.Request) (res *ObjectMultiDeleteResult, err error) {
	resp, err := u.requestWithResp(req)
	if err != nil {
		return nil, err
	}

	err = u.responseParse(resp)
	if err != nil {
		return nil, err
	}

	if !VerifyHTTPCode(resp.StatusCode) {
		return nil, fmt.Errorf("Remote response code is %d - %s not 2xx call DumpResponse(true) show details",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	retXML := new(ObjectMultiDeleteResult)
	err = xml.Unmarshal(u.LastResponseBody, retXML)
	if err != nil {
		return nil, err
	}

	return retXML, nil
}
