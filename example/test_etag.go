package main

import (
	"log"
	"strings"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	filePath   = "./FakeBigFile.txt"
	configFile = "config.json"
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

	etag := ufsdk.GetFileEtag(filePath)

	err = req.AsyncMPut(filePath, "test.txt", "", "")
	if err != nil {
		log.Println("错误返回，失败的HTTP Response 为：")
		log.Printf("%s", req.DumpResponse(true))
		panic(err.Error())
	}
	log.Println("成功返回，返回的 HTTP Response 为: ")
	log.Printf("%s\n", req.DumpResponse(false))
	err = req.HeadFile("test.txt")
	if err != nil {
		panic(err.Error())
	}
	remoteEtag := strings.Trim(req.LastResponseHeader.Get("Etag"), "\"")
	if remoteEtag != etag {
		log.Println("etag 算法不一致", remoteEtag, etag)
	} else {
		log.Println("etag 算法一致")
	}
}
