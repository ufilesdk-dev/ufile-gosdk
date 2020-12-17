package main

import (
	ufsdk "github.com/kuixiao/ufile-gosdk"
	"github.com/kuixiao/ufile-gosdk/example/helper"
	"log"
	"os"
)

const (
	ConfigFile = "./config.json"
	FilePath = "./FakeSmallFile.txt"
	KeyName = "PutKeyName"
)

func main() {
	if _, err := os.Stat(FilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(FilePath, helper.FakeSmallFileSize)
	}

	config, err := ufsdk.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	// 流式上传本地小文件
	f, err := os.Open(FilePath)
	if err != nil {
		panic(err.Error())
	}
	err = req.IOPut(f, KeyName, "")
	if err != nil {
		log.Fatalf("%s\n", req.DumpResponse(true))
	}
	log.Println("文件上传成功")

	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	log.Printf(" %s", req.LastResponseHeader)

}
