package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "config.json"
	FilePath = "test-files/test.txt"
	//KeyName = "test.txt"
	MimeType = ""
)

func main() {
	//testKeyName()

	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	// 上传一个文件
	KeyName := "/prefixB/test.txt"
	err = req.PutFile(FilePath, KeyName, MimeType)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("%s\n", err.Error())
	}
	log.Println("文件上传成功!!")

	// 获取文件的基本信息
	err = req.HeadFile(KeyName)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("查询文件信息失败，具体错误详情：%s", err.Error())
	}
	log.Println("查询文件信息成功，文件基本信息：", req.LastResponseHeader)
}


func testKeyName() {
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	// 上传一个文件
	KeyName := "test_" + "#" + "_file.txt"
	err = req.PutFile(FilePath, KeyName, MimeType)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
	} else {
		log.Println("文件上传成功!!")
	}

	KeyName = "test_" + "?" + "_file.txt"
	err = req.PutFile(FilePath, KeyName, MimeType)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
	} else {
		log.Println("文件上传成功!!")
	}

	KeyName = "test_" + "%" + "_file.txt"
	err = req.PutFile(FilePath, KeyName, MimeType)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
	} else {
		log.Println("文件上传成功!!")
	}


}