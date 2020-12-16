package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	KeyName = "testtest.txt"
)

func main() {

	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	err = req.HeadFile(KeyName)
	if err != nil {
		log.Println("查询文件信息失败,返回信息为：", string(req.DumpResponse(true)))
		log.Fatalf("查询文件信息失败，具体错误详情：%s", string(req.LastResponseBody))
	}
	for key, val := range req.LastResponseHeader {
		log.Println(key, "：", val)
	}
	log.Println("查询文件信息成功,返回信息为：", req.LastResponseHeader)
}
