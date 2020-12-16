package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "test.txt"
	KeyName = "test-demotest.txt"
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

	// 简单上传本地文件
	ok := req.CompareFileEtag(KeyName, FilePath)
	if !ok {
		log.Printf("本地文件 %s 和远程文件 %s 不一致！",FilePath,  KeyName)
	} else {
		log.Printf("本地文件 %s 和远程文件 %s 一致！",FilePath,  KeyName)
	}
}
