package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"io/ioutil"
	"log"
	"time"
)

const (
	filePath   = "./FakeSmallFile.txt"
	configFile = "config.json"
	filekey    = "test.txt"
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
	err = req.PutFile(filePath, filekey, "")
	if err != nil {
		log.Printf("上传文件失败，错误信息为：%s\n", req.DumpResponse(true))
		return
	}

	log.Println("正在下载文件。。。。")
	err = req.Download(req.GetPrivateURL(filekey, 24*time.Hour))
	if err != nil {
		log.Println("下载文件出错，出错信息为：", err.Error())
	}

	err = ioutil.WriteFile(filekey, req.LastResponseBody, 0755)
	if err != nil {
		log.Println("写本地文件失败，失败信息为：", err.Error())
	} else {
		log.Println("文件下载成功，文件名为：", filekey)
	}
}
