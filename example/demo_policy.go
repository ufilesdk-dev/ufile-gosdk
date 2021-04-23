package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "./demo_policy.go"
	KeyName = "demo_policy.go"
	MimeType = ""
	jsonPolicy = `{"callbackUrl":"https://xxxxx","callbackBody":"{\"Name\":\"Alice\", \"Age\":20}","callbackBodyType": "application/json"}`
	// 注意 "x:age": "20"  必须为字符串
	queryPolicy = `{"callbackUrl":"http://xxxxx","callbackBody":"Name=$(x:name)&Age=$(x:age)","callbackVar":{"x:name": "Alice","x:age": "20"}}`
)

func main() {
	// 加载配置，构造请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	// 简单上传回调
	err = req.PutFileWithPolicy(FilePath, KeyName, MimeType, jsonPolicy)
	if err != nil {
		log.Printf("PutWithPolicy 上传回调失败，错误信息为：%s\n", req.DumpResponse(true))
	} else {
		log.Printf("PutWithPolicy 上传回调成功")
	}

	// 同步分片上传回调
	err = req.MPutWithPolicy(FilePath, KeyName, MimeType, jsonPolicy)
	if err != nil {
		log.Printf("MPutWithPolicy 上传文件失败，错误信息为：%s\n", req.DumpResponse(true))
	} else {
		log.Printf("MPutWithPolicy 上传回调成功")
	}

	// 异步分片上传回调
	err = req.AsyncMPutWithPolicy(FilePath, KeyName, MimeType, queryPolicy)
	if err != nil {
		log.Printf("AsyncMPutWithPolicy 上传文件失败，错误信息为：%s\n", req.DumpResponse(true))
	} else {
		log.Printf("AsyncMPutWithPolicy 上传回调成功")
	}

	// 异步分片并发上传回调
	jobs := 20 // 并发数为 20
	err = req.AsyncUploadWithPolicy(FilePath, KeyName, MimeType, jobs, queryPolicy)
	if err != nil {
		log.Printf("AsyncUploadWithPolicy 上传回调失败，错误信息为：%s\n", req.DumpResponse(true))
	} else {
		log.Printf("AsyncUploadWithPolicy 上传回调成功")
	}

	err = req.DeleteFile(KeyName)
	if err != nil {
		log.Printf("DeleteFile失败，错误信息为：%s\n", req.DumpResponse(true))
	}
	log.Printf("DeleteFile成功")
}
