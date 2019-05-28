package main

import (
	"flag"
	"log"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

var (
	key    = flag.String("k", "", "upload file key.")
	path   = flag.String("p", "", "upload file path.")
	thread = flag.Int("t", 0, "upload concurrent thread count.")
	config = flag.String("c", "", "config file")
)

func main() {
	flag.Parse()
	if *config == "" || *path == "" || *key == "" {
		flag.PrintDefaults()
		return
	}
	ufConfig, err := ufsdk.LoadConfig(*config)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(ufConfig, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在上传文件。。。。")

	err = req.AsyncUpload(*path, *key, "", "", *thread)
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")
}
