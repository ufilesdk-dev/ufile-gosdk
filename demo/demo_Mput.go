package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "test-files/mongo.tgz"
	KeyName = "mongo"
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

	// 同步分片上传本地文件
	err = req.MPut(FilePath,  KeyName, MimeType)
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
