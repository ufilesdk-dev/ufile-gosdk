package main

import (
	"errors"
	"flag"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"github.com/ufilesdk-dev/ufile-gosdk/example/helper"
	"log"
	"os"
)

var (
	put  = flag.String("put", "", "test put io interface")
	mput = flag.String("mput", "", "test mput io interface")
)

const (
	PutType  = 0
	MputType = 1
)

func main() {
	flag.Parse()
	if *put == "" && *mput == "" {
		flag.PrintDefaults()
		return
	}
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

	if *put != "" {
		err = testUpload(req, helper.FakeSmallFilePath, *put, PutType)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	if *mput != "" {
		err = testUpload(req, helper.FakeBigFilePath, *mput, MputType)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Println("接口测试通过。")
}

func testUpload(req *ufsdk.UFileRequest, filePath, keyName string, uploadType int) error {
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
		return errors.New("wrong upload type")
	}
	if err != nil {
		log.Printf("%s\n", req.DumpResponse(true))
		return err
	}
	if req.CompareFileEtag(keyName, filePath) == false {
		err = errors.New("接口测试失败，上传的文件etag无法与本地文件etag匹配上。")
	}
	return err
}
