package main

import (
	"context"
	"log"
	"time"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	uploadFile    = "./FakeBigFile.txt"
	configFile    = "config.json"
	remoteFileKey = "test.txt"
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

	exampeAsyncMput(req, 500*time.Millisecond)
	exampeMput(req, 500*time.Millisecond)
}

func exampeMput(req *ufsdk.UFileRequest, timeout time.Duration) {
	var cancelFunc context.CancelFunc
	req.Context, cancelFunc = context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	err := req.MPut(uploadFile, remoteFileKey, "")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("错误：context 异常。")
}

func exampeAsyncMput(req *ufsdk.UFileRequest, timeout time.Duration) {
	var cancelFunc context.CancelFunc
	req.Context, cancelFunc = context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	err := req.AsyncMPut(uploadFile, remoteFileKey, "")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("错误：context 异常。")
}
