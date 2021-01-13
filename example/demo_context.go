package main

import (
	"context"
	"log"
	"time"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	ConfigFile = "./config.json"
	FilePath = "mongo.zip"
	KeyName = "mongo.zip"
	MimeType = ""
)

func main() {
	log.SetFlags(log.Lshortfile)

	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	exampleMput(req, 500*time.Millisecond)
	exampleAsyncMput(req, 500*time.Millisecond)
}

func exampleMput(req *ufsdk.UFileRequest, timeout time.Duration) {
	var cancelFunc context.CancelFunc
	req.Context, cancelFunc = context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	err := req.MPut(FilePath, KeyName, MimeType)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("错误：context 异常。")
}

func exampleAsyncMput(req *ufsdk.UFileRequest, timeout time.Duration) {
	var cancelFunc context.CancelFunc
	req.Context, cancelFunc = context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	err := req.AsyncMPut(FilePath, KeyName, MimeType)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("错误：context 异常。")
}
