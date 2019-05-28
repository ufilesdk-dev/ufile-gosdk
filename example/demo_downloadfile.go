package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"os"
)

const (
	uploadFile    = "./FakeBigFile.txt"
	configFile    = "config.json"
	remoteFileKey = "test.txt"
	saveAsName    = "download.txt"
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

	err = req.PutFile(uploadFile, remoteFileKey, "", "")
	if err != nil {
		log.Printf("上传文件失败，错误信息为：%s\n", req.DumpResponse(true))
		return
	}

	log.Println("正在下载文件。。。。")
	file, err := os.OpenFile(saveAsName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Println("创建文件失败，错误信息为：", err.Error())
		return
	}

	err = req.DownloadFile(file, remoteFileKey)
	if err != nil {
		log.Println("下载文件出错，出错信息为：", err.Error())
	}
	file.Close() //提前关闭文件，防止 etag 计算不准。

	etagCheck := req.CompareFileEtag(remoteFileKey, saveAsName)
	if !etagCheck {
		log.Println("文件下载出错，etag 比对不一致。")
	} else {
		log.Println("文件下载成功")
	}
}
