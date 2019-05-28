package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	uploadFile    = "./FakeBigFile.txt"
	configFile    = "config.json"
	remoteFileKey = "/test.txt"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在上传文件。。。。")

	err = req.AsyncMPut(uploadFile, remoteFileKey, "", "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	checkEtag := req.CompareFileEtag(remoteFileKey, uploadFile)
	if checkEtag {
		log.Println("文件上传成功。")
	} else {
		log.Println("文件上传失败。")
	}
}
