package filemgr

import (
	"io/ioutil"
	"ufile-gosdk/common"
	"encoding/json"
	"ufile-gosdk/uflog"
	"fmt"
)

type DeleteFileResponse struct {
	RetCode       int
	ErrMsg        string
	ContentLength string
	XSessionId    string
}

func DeleteFile(bucketname string, filename string, publickey string, privatekey string, proxyhost string) (DeleteFileResponse, error) {
	deletefileresponse := DeleteFileResponse{}

	authorization := common.Authorization(bucketname, filename, publickey, privatekey, "DELETE", "", "", "", "")

	if len(proxyhost) <= 4 {
		return deletefileresponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	res, err := common.FileApiDelete(bucketname, suffix, filename, authorization)
	if err != nil {
		return deletefileresponse, err
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return deletefileresponse, err
	}

	if res.StatusCode != 200 && res.StatusCode != 204 {
		deletefileresponse.XSessionId = res.Header.Get("X-Sessionid")

		uflog.ERROR(string(result))

		err = json.Unmarshal(result, &deletefileresponse)
		if err != nil {
			return deletefileresponse, err
		}

		return deletefileresponse, fmt.Errorf(deletefileresponse.ErrMsg)
	}

	deletefileresponse.RetCode = 0
	deletefileresponse.ContentLength = res.Header.Get("Content-Length")

	uflog.INFO(deletefileresponse)

	return deletefileresponse, nil
}