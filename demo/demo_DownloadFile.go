package main

import (
	"bytes"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"os"
)

const (
	ConfigFile = "./config.json"
	KeyName = "test.txt"
	FilePath = "localFile.txt"
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

	// 下载到缓存
	var buf bytes.Buffer
	err = req.DownloadFile(&buf, KeyName)
	if err != nil {
		log.Fatalf("流式下载到缓存出错，出错信息为：%s", err.Error())
	}
	log.Printf("流式下载到缓存成功！数据为：%s", buf.String())

	// 下载到本地文件
	file, err := os.OpenFile(FilePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalf("创建文件失败，错误信息为：%s", err.Error())
	}
	err = req.DownloadFile(file, KeyName)
	if err != nil {
		log.Fatalf("流式下载到本地文件出错，错误信息为：%s", err.Error())
	}
	defer file.Close()
	log.Printf("流式下载到本地文件成功!")
}
