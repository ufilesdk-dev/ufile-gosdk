package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"os"
)

const (
	ConfigFile = "./config.json"
	FilePath = "installer.exe"
	KeyName = "Docker3.exe"
	MimeType = ""
)

func main() {

	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	// 异步分片上传本地文件
	f, err := os.Open(FilePath)
	if err != nil {
		panic(err.Error())
	}
	err = req.IOMultipartAsyncUpload(f, KeyName, MimeType)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	log.Println("文件上传成功!!")

	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	log.Printf(" %s", req.LastResponseHeader)
}
