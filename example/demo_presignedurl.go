package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"github.com/ufilesdk-dev/ufile-gosdk/example/helper"
)

func main() {
	log.SetFlags(log.Lshortfile)
	if _, err := os.Stat(helper.FakeBigFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(helper.FakeBigFilePath, helper.FakeBigFileSize)
	}
	config, err := ufsdk.LoadConfig(helper.ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	u, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	fileKey := helper.GenerateUniqKey()
	err = putfile(helper.FakeBigFilePath, fileKey, u)
	if err != nil {
		log.Println("上传文件失败，具体错误详情：", err.Error())
		return
	}
	log.Println("上传文件成功")

	fileKey = helper.GenerateUniqKey()
	err = mputfile(helper.FakeBigFilePath, fileKey, u)
	if err != nil {
		log.Println("分片上传文件失败，具体错误详情：", err.Error())
		return
	}
	log.Println("分片上传文件成功")
}

func putfile(filePath string, keyName string, u *ufsdk.UFileRequest) error {
	// 请确保在服务端生成该签名URL时设置的请求头与在使用URL时设置的请求头一致
	// u.RequestHeader = http.Header{}
	// u.RequestHeader.Set("key", "value")
	// 根据请求添加query
	// u.RequestHeader = url.Values{}
	// u.query.Set("key", "value")

	signedUrl := u.GenPresignedURL(keyName, 0, "PUT")
	log.Println("上传文件的url 为：", signedUrl)
	//使用url上传文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", signedUrl, file)
	if err != nil {
		return err
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("返回上传状态码: ", resp.StatusCode)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Remote response code is %d - %s not 2xx call DumpResponse(true) show details",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}

// 分片上传
func mputfile(filePath string, keyName string, u *ufsdk.UFileRequest) error {
	// 创建一个新的HTTP客户端
	client := &http.Client{}

	//初始化分片上传
	state, err := u.InitiateMultipartUpload(keyName, "")
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	chunk := make([]byte, state.BlkSize)
	var partNumber int
	for {
		bytesRead, fileErr := file.Read(chunk)
		if fileErr == io.EOF || bytesRead == 0 { //后面直接读到了结尾
			break
		}

		// 获取分片上传的url
		query := url.Values{}
		query.Set("uploadId", state.UploadID)
		query.Set("partNumber", strconv.Itoa(partNumber))
		u.RequestQuery = query
		signedUrl := u.GenPresignedURL(keyName, 0, "PUT")
		log.Println("上传分片的url 为：", signedUrl)

		// 使用url上传分片
		buf := bytes.NewBuffer(chunk[:bytesRead])
		req, err := http.NewRequest("PUT", signedUrl, buf)
		if err != nil {
			return err
		}

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// 读取响应
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Println("返回上传状态码: ", resp.StatusCode)
		if resp.StatusCode != 200 {
			return fmt.Errorf("Remote response code is %d - %s not 2xx call DumpResponse(true) show details",
				resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		if err != nil {
			u.AbortMultipartUpload(state)
			return err
		}
		// 解析出响应中的etag
		etag := strings.Trim(resp.Header.Get("Etag"), "\"") //为保证线程安全，这里就不保留 lastResponse
		if etag == "" {
			etag = strings.Trim(resp.Header.Get("ETag"), "\"") //为保证线程安全，这里就不保留 lastResponse
		}
		state.Etags[partNumber] = etag
		partNumber++
	}

	return u.FinishMultipartUpload(state)
}
