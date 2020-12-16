package main

import (
	"fmt"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"net/http"
)

const (
	ConfigFile = "./config.json"
	FilePath = "test.txt"
	KeyName = "test.txt"
	newKeyName = "new_test2.txt"
	force = "false"
	MimeType = ""
)

func main() {

	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	//req, err := ufsdk.NewFileRequest(config, nil)
	//if err != nil {
	//	panic(err.Error())
	//}
	header := make(http.Header)

	//1、上传标准存储类型文件
	header.Add("X-Ufile-Storage-Class", "IA")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}

	// 简单上传本地文件
	err = req.PutFile(FilePath, KeyName, MimeType)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	log.Println("文件上传成功!!")

	err = req.Rename(KeyName, newKeyName, force)
	if err != nil {
		log.Fatalf("文件重命名失败，具体错误详情：%s", err.Error())
	}
	log.Println("文件重命名成功,返回信息为：", req.LastResponseHeader)

	err = req.HeadFile(newKeyName)
	if err != nil {
		log.Fatalf("获取文件基本信息失败，具体错误详情：%s", err.Error())
	}
	log.Println("获取文件基本信息成功,返回信息为：", req.LastResponseHeader)
	log.Println("获取文件基本信息成功,返回信息为：", req.LastResponseBody)

	fmt.Scanln()
	// 强制更名
	req.RequestHeader.Set("X-Ufile-Storage-Class", "STANDARD")
	err = req.PutFile("test.json", KeyName, MimeType)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	log.Println("文件上传成功!!")

	err = req.Rename(KeyName, newKeyName, "ttt")
	if err != nil {
		log.Printf("文件重命名失败，具体错误详情：%s", err.Error())
		log.Fatalf("文件重命名失败，具体错误详情：%s", string(req.DumpResponse(true)))
	}
	log.Println("文件重命名成功,返回信息为：", req.LastResponseHeader)

	err = req.HeadFile(newKeyName)
	if err != nil {
		log.Fatalf("获取文件基本信息失败，具体错误详情：%s", err.Error())
	}
	log.Println("获取文件基本信息成功,返回信息为：", req.LastResponseHeader)
}