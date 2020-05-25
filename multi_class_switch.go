package ufsdk

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type SwitchFile struct {
	Key          string `xml:"Key"`
	StorageClass string `xml:"StorageClass"`
}

type MultiClassSwitchFile struct {
	XMLName xml.Name     `xml:"ClassSwitch"`
	IsQuit  bool         `xml:"Quiet"`
	Objects []SwitchFile `xml:"Object"`
}

type ObjectMultiSwitchResult struct {
	XMLName xml.Name     `xml:"ClassSwitchResult"`
	Objects []SwitchFile `xml:"ClassSwitch"`
}

//MultiDeleteFile 批量类型转换接口
//input 批量类型转换接口，IsQuit 返回的格式，默认false
//简单模式（quiet）：ufile返回的消息体中只包含转换过程中出错的Object结果。如果所有转换都成功，则没有消息体。
//详细模式（verbose）：ufile返回的消息体中会包含所有转换Object的结果。默认采用详细模式。
//有效值：true（开启简单模式）、false（开启详细模式）
func (u *UFileRequest) MultiClassSwitchFile(input []SwitchFile, isQuit bool) (res *ObjectMultiSwitchResult, err error) {
	//鉴权
	mimeType := "multipart/form-data"
	h := make(http.Header)
	h.Add("Content-Type", mimeType)

	authorization := u.Auth.Authorization("POST", u.BucketName, "xml", h)

	buffer, err := getXmlDataSwitch(input, isQuit) //上传的文件名
	if err != nil {
		return nil, err
	}
	//内容体构造
	boundry := makeBoundry()
	body := makeFormBodyV2(authorization, boundry, "xml", mimeType, false, buffer)
	//lastline 一定要写，否则后端解析不到。
	lastline := fmt.Sprintf("\r\n--%s--\r\n", boundry)
	body.Write([]byte(lastline))

	url := u.MultiClassSwitch()
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

	return u.FinishMultiClassSwitch(req)
}

func getXmlDataSwitch(input []SwitchFile, isQuit bool) (io.Reader, error) {
	xmlFile := new(MultiClassSwitchFile)
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

func (u *UFileRequest) FinishMultiClassSwitch(req *http.Request) (res *ObjectMultiSwitchResult, err error) {
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
	retXML := new(ObjectMultiSwitchResult)
	err = xml.Unmarshal(u.LastResponseBody, retXML)
	if err != nil {
		return nil, err
	}

	return retXML, nil
}
