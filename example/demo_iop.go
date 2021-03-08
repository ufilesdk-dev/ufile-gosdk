package main

import (
	"log"
	"os"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	localUpFile      = "./stream-up.png"
	localDlFile      = "./stream-dl.png"
	iopConfigFile    = "./config.json"
	iopRemoteFileKey = "picture.png"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(iopConfigFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	if _, err := os.Stat(localUpFile); os.IsNotExist(err) {
		panic(err.Error())
	}
	log.Println("正在上传文件...")

	//构建iop命令，缩放为原图50%
	iopcmdString := "iopcmd=thumbnail&type=1&scale=50"
	// 通过直接指定iop字符串执行上传iop, iopcmdString为自己构建的iop命令
	err = req.PutFileWithIopString(localUpFile, iopRemoteFileKey, "", iopcmdString)
	if err != nil {
		log.Fatalf("iop上传文件失败，错误信息为：%s\n", req.DumpResponse(true))
	} else {
		log.Println("iop上传文件成功")
	}

	log.Println("正在下载文件...")
	file, err := os.OpenFile(localDlFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Println("创建文件失败，错误信息为：", err.Error())
		return
	}

	//构建iop命令，缩放为原图50%
	iopcmdString = "iopcmd=thumbnail&type=1&scale=50"
	// 通过直接指定iop字符串执行下载iop, iopcmdString为自己构建的iop命令
	err = req.DownloadFileWithIopString(file, iopRemoteFileKey, iopcmdString)
	if err != nil {
		log.Println("下载文件出错，出错信息为：", err.Error())
	}
	file.Close()
}
