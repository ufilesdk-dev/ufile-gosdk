package filemgr

import (
	"ufile-gosdk/common"
	"io/ioutil"
	"encoding/json"
	"io"
	"os"
	"net/url"
	"strconv"
	"net/http"
	"ufile-gosdk/uflog"
	"fmt"
	"strings"
)

type GetFileResponse struct {
	RetCode       int
	ErrMsg        string
	Etag          string
	ContentType   string
	ContentLength string
	ContentRange  string
}

func GetFile(bucketname string, path string, filename string, byterange string, isprivate bool, expires int64, publickey string, privatekey string, proxyhost string, connectionTimeout int) (GetFileResponse, error) {
	getfileresponse := GetFileResponse{}

	data := url.Values{}
	data.Add("UCloudPublicKey", publickey)

	data.Add("Expires", strconv.FormatInt(expires, 10))

	singaturedata := common.Signature(privatekey, data)

	authorization := common.Authorization(bucketname, filename, publickey, privatekey, "GET", "", "", "", "")

	if len(proxyhost) <= 4 {
		return getfileresponse, fmt.Errorf("ProxyHost Illegal")
	}
	suffix := string([]byte(proxyhost)[3:])

	pathrune := []rune(path)

	if string(pathrune[len(pathrune) - 1:]) != "/" {
		path = path + "/"
	}

	var res *http.Response
	var err error

	if isprivate {
		res, err = common.FileApiGetPrivate(bucketname, suffix, filename, byterange, authorization, singaturedata, connectionTimeout)

	} else {
		res, err = common.FileApiGetPublic(bucketname, suffix, filename, byterange, authorization, connectionTimeout)
	}
	if err != nil {
		return getfileresponse, err
	}

	if res.StatusCode != 200 {
		result, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return getfileresponse, err
		}

		err = json.Unmarshal(result, &getfileresponse)
		if err != nil {
			return getfileresponse, err
		}

		uflog.ERROR(getfileresponse)

		return getfileresponse, fmt.Errorf(getfileresponse.ErrMsg)
	}

	filepath := strings.Split(filename, "/")
	pathX := path

	for k, v := range filepath {
		if k == len(filepath) - 1 {
			break
		}

		if v != "" {
			pathX = pathX + v + "/"
		}
	}

	_, err = os.Stat(pathX)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pathX, os.ModePerm)
		if err != nil {
			return getfileresponse, err
		}
	}

	f, err := os.Create(path + filename)
	if err != nil {
		return getfileresponse, err
	}
	io.Copy(f, res.Body)
	defer f.Close()

	getfileresponse.RetCode = 0
	getfileresponse.Etag = res.Header.Get("Etag")
	getfileresponse.ContentLength = res.Header.Get("Content-Length")
	getfileresponse.ContentRange = res.Header.Get("Content-Range")
	getfileresponse.ContentType = res.Header.Get("Content-Type")

	uflog.INFO(getfileresponse)

	return getfileresponse, nil

}
