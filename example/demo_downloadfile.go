package main

import (
	"github.com/kuixiao/ufile-gosdk/example/helper"
	ufsdk "github.com/kuixiao/ufile-gosdk"
	"log"
	"os"
	"time"
)

const (
	uploadFile    = "./FakeSmallFile.txt"
	configFile    = "config.json"
	remoteFileKey = "remoteFileKey.txt"
	downloadPath    = "download.txt"
	downloadFilePath    = "downloadFile.txt"
)

func main() {
	// 准备下载请求与要下载的文件
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
	if _, err := os.Stat(uploadFile); os.IsNotExist(err) {
		helper.GenerateFakefile(uploadFile, helper.FakeSmallFileSize)
	}
	err = req.PutFile(uploadFile, remoteFileKey, "")
	if err != nil {
		log.Println("上传文件失败，错误信息为：", req.DumpResponse(true))
	}

	// 下载文件
	log.Println("下载文件: ")
	reqUrl := req.GetPrivateURL(remoteFileKey, 10*time.Second)
	err = req.Download(reqUrl)
	if err != nil {
		log.Fatalln(string(req.DumpResponse(true)))
	}
	log.Printf("下载文件成功！数据为：%s", string(req.LastResponseBody))
	// 保存到本地
	f, err := os.OpenFile(downloadPath, os.O_CREATE|os.O_WRONLY, 0755)
	defer f.Close()
	if err != nil {
		log.Fatalf("创建文件失败，错误信息为：%s", err.Error())
	}
	_, err = f.WriteString(string(req.LastResponseBody))
	if err != nil {
		log.Println("下载数据保存到本地文件失败：", err.Error())
	} else {
		log.Printf("下载数据保存到本地文件%s成功：", downloadPath)
	}
	// 检查所保存文件与远程文件etag是否一致
	etagCheck := req.CompareFileEtag(remoteFileKey, downloadPath)
	if !etagCheck {
		log.Println("文件下载出错，etag 比对不一致。")
	} else {
		log.Println("文件下载成功, etag 一致")
	}

	// 流式下载文件
	log.Println("流式下载文件: ")
	file, err := os.OpenFile(downloadFilePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln("创建文件失败，错误信息为：", err.Error())
	}
	err = req.DownloadFile(file, remoteFileKey)
	if err != nil {
		log.Println("下载文件出错，出错信息为：", err.Error())
	} else {
		log.Printf("下载文件保存到本地文件%s成功：", downloadFilePath)
	}
	defer file.Close() //提前关闭文件，防止 etag 计算不准。
	// 检查所保存文件与远程文件etag是否一致
	etagCheck = req.CompareFileEtag(remoteFileKey, downloadFilePath)
	if !etagCheck {
		log.Println("文件下载出错，etag 比对不一致。")
	} else {
		log.Println("文件下载成功, etag 一致")
	}

	// 删除文件
	err = req.DeleteFile(remoteFileKey)
	if err != nil {
		log.Println("远程文件删除失败：", err.Error())
	} else {
		log.Println("远程文件删除成功")
	}


}
