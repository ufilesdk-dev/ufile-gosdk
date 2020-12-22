package main

import (
	"flag"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"github.com/ufilesdk-dev/ufile-gosdk/example/helper"
	"log"
	"os"
)

const (
	ConfigFile = "config.json"
	IOPutKeyName = "PutKeyName"
	IOMPutKeyName = "MPutKeyName"
	MimeType = ""
)

func main() {
	flag.Parse()

	log.SetFlags(log.Lshortfile)
	if _, err := os.Stat(helper.FakeSmallFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(helper.FakeSmallFilePath, helper.FakeSmallFileSize)
	}
	if _, err := os.Stat(helper.FakeBigFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(helper.FakeBigFilePath, helper.FakeBigFileSize)
	}

	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	// 流式上传本地小文件
	f, err := os.Open(helper.FakeSmallFilePath)
	if err != nil {
		panic(err.Error())
	}
	err = req.IOPut(f, IOPutKeyName, MimeType)
	if err != nil {
		log.Fatalf("%s\n", req.DumpResponse(true))
	}
	log.Println("文件上传成功")

	checkEtag := req.CompareFileEtag(IOPutKeyName, helper.FakeSmallFilePath)
	if !checkEtag {
		log.Fatalln("CompareFileEtag 失败。")
	}
	log.Println("CompareFileEtag 成功。")

	//  流式分片上传本地文件
	f1, err := os.Open(helper.FakeBigFilePath)
	if err != nil {
		panic(err.Error())
	}
	err = req.IOMutipartAsyncUpload(f1, IOMPutKeyName, MimeType)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	log.Println("文件上传成功!!")

	checkEtag = req.CompareFileEtag(IOMPutKeyName, helper.FakeBigFilePath)
	if !checkEtag {
		log.Fatalln("CompareFileEtag 失败。")
	}
	log.Println("CompareFileEtag 成功。")
}
