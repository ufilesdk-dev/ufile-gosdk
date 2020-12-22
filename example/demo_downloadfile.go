package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"github.com/ufilesdk-dev/ufile-gosdk/example/helper"
	"log"
	"os"
	"time"
)

const (
	ConfigFile = "./config.json"
	SmallFilePath = "config.json"
	SmallFileKeyName = "config.json"
	SmallFileMimeType = ""
	DownloadPath = "dl_config.json"

	BigFilePath = "mongo.zip"
	BigFileKeyName = "mongo.zip"
	BigFileMimeType = ""
	DownloadFilePath = "dl_mongo.zip"
)

func main() {
	// 准备下载请求与要下载的文件
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在上传文件。。。。")
	if _, err := os.Stat(SmallFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(SmallFilePath, 1024)
	}
	err = req.PutFile(SmallFilePath, SmallFileKeyName, SmallFileMimeType)
	if err != nil {
		log.Println("上传文件失败，错误信息为：", string(req.DumpResponse(true)))
	}
	log.Printf("文件%s上传成功", SmallFilePath)
	if _, err := os.Stat(BigFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(BigFilePath, helper.FakeBigFileSize)
	}
	err = req.MPut(BigFilePath, BigFileKeyName, BigFileMimeType)
	if err != nil {
		log.Println("上传文件失败，错误信息为：", string(req.DumpResponse(true)))
	}
	log.Printf("文件%s上传成功", BigFilePath)

	// 下载文件
	log.Println("下载文件: ")
	reqUrl := req.GetPrivateURL(SmallFileKeyName, 10*time.Second)
	err = req.Download(reqUrl)
	if err != nil {
		log.Fatalln(string(req.DumpResponse(true)))
	}
	log.Printf("下载文件成功！")
	// 保存到本地
	f, err := os.OpenFile(DownloadPath, os.O_CREATE|os.O_WRONLY, 0755)
	defer f.Close()
	if err != nil {
		log.Fatalf("创建文件失败，错误信息为：%s", err.Error())
	}
	_, err = f.WriteString(string(req.LastResponseBody))
	if err != nil {
		log.Println("下载数据保存到本地文件失败：", err.Error())
	} else {
		log.Printf("下载数据保存到本地文件%s成功：", DownloadPath)
	}
	// 检查所保存文件与远程文件etag是否一致
	etagCheck := req.CompareFileEtag(SmallFileKeyName, DownloadPath)
	if !etagCheck {
		log.Println("文件下载出错，etag 比对不一致。")
	} else {
		log.Println("文件下载成功, etag 一致")
	}

	// 流式下载文件
	log.Println("流式下载文件: ")
	file, err := os.OpenFile(DownloadFilePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln("创建文件失败，错误信息为：", err.Error())
	}
	err = req.DownloadFile(file, BigFileKeyName)
	if err != nil {
		log.Println("下载文件出错，出错信息为：", err.Error())
	} else {
		log.Printf("下载文件保存到本地文件%s成功：", DownloadFilePath)
	}
	file.Close() //提前关闭文件，防止 etag 计算不准。
	// 检查所保存文件与远程文件etag是否一致
	etagCheck = req.CompareFileEtag(BigFileKeyName, DownloadFilePath)
	if !etagCheck {
		log.Println("文件下载出错，etag 比对不一致。")
	} else {
		log.Println("文件下载成功, etag 一致")
	}

	// 删除远程小文件
	err = req.DeleteFile(SmallFileKeyName)
	if err != nil {
		log.Println("远程文件删除失败：", err.Error())
	} else {
		log.Printf("远程文件 %s 删除成功", SmallFileKeyName)
	}
	// 删除远程大文件
	err = req.DeleteFile(BigFileKeyName)
	if err != nil {
		log.Println("远程文件删除失败：", err.Error())
	} else {
		log.Printf("远程文件%s删除成功", BigFileKeyName)
	}


}
