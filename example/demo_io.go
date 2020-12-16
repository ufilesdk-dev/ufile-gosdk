package main

import (
	"errors"
	"flag"
	ufsdk "github.com/kuixiao/ufile-gosdk"
	"github.com/kuixiao/ufile-gosdk/example/helper"
	"log"
	"os"
)

const (
	PutType  = 0
	MputType = 1
	PutKeyName = "PutKeyName"
	MPutKeyName = "MPutKeyName"
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

	config, err := ufsdk.LoadConfig(helper.ConfigFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	testUpload(req, helper.FakeSmallFilePath, PutKeyName, PutType)

	testUpload(req, helper.FakeBigFilePath, MPutKeyName, MputType)
}

func testUpload(req *ufsdk.UFileRequest, filePath, keyName string, uploadType int)  {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	switch uploadType {
	case PutType:
		err = req.IOPut(f, keyName, "")
	case MputType:
		err = req.IOMutipartAsyncUpload(f, keyName, "")
	default:
		return
	}
	if err != nil {
		log.Fatalf("%s\n", req.DumpResponse(true))
	}
	if req.CompareFileEtag(keyName, filePath) == false {
		log.Fatalf("接口测试失败，上传的文件etag无法与本地文件etag匹配上。")
	}
	log.Printf("文件 %s 上传成功", filePath)
}
