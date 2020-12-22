package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"github.com/ufilesdk-dev/ufile-gosdk/example/helper"
	"log"
	"os"
)

const (
	ConfigFile = "./config.json"
	FilePath = "mongo.zip"
	KeyName = "mongo.zip"
	MimeType = ""
)

func main() {
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在上传文件。。。。")

	if _, err := os.Stat(FilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(FilePath, helper.FakeBigFileSize)
	}

	err = req.AsyncMPut(FilePath, KeyName, MimeType)
	if err != nil {
		log.Fatalln("文件上传失败，失败原因：", err.Error())
	}
	log.Println("文件上传成功。")

	checkEtag := req.CompareFileEtag(KeyName, FilePath)
	if !checkEtag {
		log.Fatalln("CompareFileEtag 失败。")
	}
	log.Println("CompareFileEtag 成功。")

	err = req.DeleteFile(KeyName)
	if err != nil {
		log.Fatalln("文件删除失败。")
	}
	log.Println("文件删除成功。")
}
