package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	ConfigFile = "./config.json"
	//FilePath = "test.txt"
	//KeyName = "test.txt"
	//MimeType = ""
)

func main() {
	nilHttpClient()
	timeoutRequest()
	headerRequest()
	httpsRequest()
	bucketRequest()
}

// nil client
func nilHttpClient() {
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	err = req.PutFile("test-files/test.txt", "test.txt", "")
	if err != nil {
		log.Printf("PutFile失败，错误信息为：%s", string(req.DumpResponse(true)))
		log.Fatalf("PutFile失败，错误信息为：%s", err.Error())
	}
	log.Printf("PutFile 成功")

	_, _ = fmt.Scanln()
	_ = req.DeleteFile("test.txt")
}

// timeout request
func timeoutRequest() {

	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	// 设置超时
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	req.Context, _ = context.WithTimeout(context.Background(), time.Second * 10)
	err = req.PutFile("test-files/test.txt", "test.txt", "")
	if err != nil {
		log.Printf("PutFile失败，错误信息为：%s", string(req.DumpResponse(true)))
		log.Fatalf("PutFile失败，错误信息为：%s", err.Error())
	}
	log.Printf("PutFile 成功")

	_, _ = fmt.Scanln()
	_ = req.DeleteFile("test.txt")
}

// 自定义请求 HTTP Header
func headerRequest() {
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	// 自定义请求 HTTP Header
	header := make(http.Header)
	header.Add("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.PutFile("test-files/test.txt", "test.txt", "")
	if err != nil {
		log.Printf("PutFile失败，错误信息为：%s", string(req.DumpResponse(true)))
		log.Fatalf("PutFile失败，错误信息为：%s", err.Error())
	}
	log.Printf("PutFile 成功")

	_, _ = fmt.Scanln()
	_ = req.DeleteFile("test.txt")

}

// 自定义client
func clientRequest() {
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	// 自定义请求 HTTP Header
	client := &http.Client{}
	client.Timeout = time.Second * 5


	header := make(http.Header)
	header.Add("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.PutFile("test-files/test.txt", "test.txt", "")
	if err != nil {
		log.Printf("PutFile失败，错误信息为：%s", string(req.DumpResponse(true)))
		log.Fatalf("PutFile失败，错误信息为：%s", err.Error())
	}
	log.Printf("PutFile 成功")

	_, _ = fmt.Scanln()
	_ = req.DeleteFile("test.txt")

}

// https
func httpsRequest() {
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	// https
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := ufsdk.NewFileRequest(config, client)
	if err != nil {
		panic(err.Error())
	}


	err = req.PutFileByHttps("test-files/test.txt", "test.txt", "")
	if err != nil {
		log.Printf("PutFile失败，错误信息为：%s", string(req.DumpResponse(true)))
		log.Fatalf("PutFile失败，错误信息为：%s", err.Error())
	}
	log.Printf("PutFile 成功，成功信息为：%s", string(req.DumpResponse(true)))
	log.Printf("PutFile 成功")

	_, _ = fmt.Scanln()
	_ = req.DeleteFile("test.txt")
}

// bucket request
func bucketRequest() {
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewBucketRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	bucketRet, err := req.CreateBucket("test-demo-newbucket", "cn-bj", "private", "")
	if err != nil {
		log.Fatalf("创建 bucket 出错，错误信息为：%s\n", err.Error())
	}
	log.Println("创建 Bucket 成功，bucket 为", bucketRet)

	_, _ = fmt.Scanln()
	_, _ = req.DeleteBucket("test-demo-newbucket", "")
}