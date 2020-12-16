package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "test-files/test.txt"
	KeyName = "test-hit.json"
)

func main() {

	// 加载配置，创建请求
	//config, err := ufsdk.LoadConfig("ConfigFile")
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	// 简单上传本地文件
	err = req.PutFile(FilePath,  KeyName, "")
	if err != nil {
		log.Println("文件上传失败 ", err.Error())
		log.Println(string(req.DumpResponse(true)))
	} else{
		log.Println("文件上传成功!!")
	}

	// 秒传文件
	err = req.UploadHit("demo_Mput.go", "UploadHit" + KeyName)
	if err != nil {
		log.Println(string(req.DumpResponse(true)))
		log.Fatalf("秒传失败，错误信息为：%s", err.Error())
	}
	log.Printf("秒传成功，状态为：%d", req.LastResponseStatus)


	err = req.HeadFile("UploadHit" + KeyName)
	if err != nil {
		log.Fatalf("获取文件基本信息失败，错误信息为：%s", err.Error())
	}
	log.Printf("获取文件基本信息成功，内容为：%s", req.LastResponseHeader)
}
