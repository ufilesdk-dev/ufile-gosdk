package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "mongo.zip"
	KeyName = "mongo.zip"
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
	log.Println("文件上传成功。")

	ok := req.CompareFileEtag(KeyName, FilePath)
	if !ok {
		log.Fatalln("CompareFileEtag 失败。")
	}
	log.Println("CompareFileEtag 成功。")

	err = req.DeleteFile(KeyName)
	if err != nil {
		log.Fatalln("文件删除失败。")
	}
	log.Println("文件删除成功。")
}