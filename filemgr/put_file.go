package filemgr

import (
	"ufile-gosdk/common"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"ufile-gosdk/uflog"
	"fmt"
)

type PutFileResponse struct {
	RetCode int
	ErrMsg  string
	Etag    string
}

func PutFile(bucketname string, path string, filename string, contentType string, publickey string, privatekey string, proxyhost string) (PutFileResponse, error) {
	putfileresponse := PutFileResponse{}

	data := url.Values{}

	authorization := common.Authorization(bucketname, filename, publickey, privatekey, "PUT", "", contentType, "", "")

	if len(proxyhost) <= 4 {
		return putfileresponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	res, err := common.FileApiPut(bucketname, suffix, path, filename, data, authorization, contentType)
	if err != nil {
		return putfileresponse, err
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return putfileresponse, err
	}

	if res.StatusCode != 200 {
		err = json.Unmarshal(result, &putfileresponse)
		if err != nil {
			return putfileresponse, err
		}

		uflog.ERROR(putfileresponse)

		return putfileresponse, fmt.Errorf(putfileresponse.ErrMsg)
	}

	putfileresponse.RetCode = 0
	putfileresponse.Etag = res.Header.Get("Etag")

	uflog.INFO(putfileresponse)

	return putfileresponse, nil
}