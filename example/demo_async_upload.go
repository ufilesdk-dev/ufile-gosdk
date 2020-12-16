package main

import (
	"flag"
	"log"
	"os"

	ufsdk "github.com/kuixiao/ufile-gosdk"
	"github.com/kuixiao/ufile-gosdk/example/helper"
)

var (
	key    = flag.String("k", "", "upload file key.")
	path   = flag.String("p", "", "upload file path.")
	thread = flag.Int("t", 0, "upload concurrent thread count.")
	config = flag.String("c", "", "config file")
)
const (
	uploadFile    = "./FakeBigFile.txt"
	configFile    = "config.json"
	remoteFileKey = "AsyncUpload.txt"
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
	log.Println("正在上传文件。。。。")

	if _, err := os.Stat(uploadFile); os.IsNotExist(err) {
		helper.GenerateFakefile(uploadFile, helper.FakeBigFileSize)
	}

	err = req.AsyncUpload(uploadFile, remoteFileKey, "", *thread)
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	checkEtag := req.CompareFileEtag(remoteFileKey, uploadFile)
	if !checkEtag {
		log.Fatalln("CompareFileEtag 失败。")
	}
	log.Println("CompareFileEtag 成功。")

	err = req.DeleteFile(remoteFileKey)
	if err != nil {
		log.Fatalln("文件删除失败。")
	}
	log.Println("文件删除成功。")
}
