package main

import (
	ufsdk "github.com/kuixiao/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "./config.json"
	KeyName = "config.json"
	MimeType = ".xml"
)

func main() {

	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	//config, err := ufsdk.LoadConfig("ConfigFile")
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	// 简单上传本地文件
	err = req.PostFile(FilePath, KeyName, MimeType)
	//err = req.PutFile("FilePath", "FileKey", "MimeType")
	if err != nil {
		log.Println(string(req.DumpResponse(true)))
		log.Fatalf("%s\n", err.Error())
	}
	log.Println("文件上传成功!!")

	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("查询文件信息失败，具体错误详情：%s", err.Error())
	}
	log.Println("查询文件信息成功,返回信息为：", req.LastResponseHeader)
}
