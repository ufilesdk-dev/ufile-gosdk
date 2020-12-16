package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "installer.exe"
	KeyName = "Docker3.exe"
	MimeType = ""
	Jobs = 20
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
	err = req.AsyncUpload(FilePath,  KeyName, MimeType, Jobs)
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

