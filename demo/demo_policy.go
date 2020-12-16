package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "config.json"
	KeyName = "mongo"
	MimeType = ""
	//putPolicy = `{"callbackUrl":"https://175.24.32.143:4321/post","callbackBody":"{\"Name\":\"Alice\", \"Age\":20}","callbackBodyType": "application/json"}`
	//mputPolicy = `{"callbackUrl":"http://175.24.32.143:8091/hello","callbackBody":"Name=$(x:name)&Age=$(x:age)","callbackVar":{"x:name": "Alice","x:age": "20"}}`
	//asyncMputPolicy = `{"callbackUrl":"http://175.24.32.143:8092/hello","callbackBody":"{\"Name\":\"Alice\", \"Age\":20}","callbackBodyType": "application/json"}`
	jsonPolicy = `{"callbackUrl":"https://175.24.32.143:4321/post","callbackBody":"{\"Name\":\"Alice\", \"Age\":20}","callbackBodyType": "application/json"}`
	// 注意 "x:age": "20"  必须为字符串
	queryPolicy = `{"callbackUrl":"http://175.24.32.143:8091/hello","callbackBody":"Name=$(x:name)&Age=$(x:age)","callbackVar":{"x:name": "Alice","x:age": "20"}}`
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
		log.Fatalf("PutWithPolicy 上传回调失败，错误信息为：%s\n", req.DumpResponse(true))
	}
	log.Printf("PutWithPolicy 上传回调成功")
	_ = req.DeleteFile(KeyName)

	// 同步分片上传回调
	err = req.MPutWithPolicy(FilePath, KeyName, MimeType, jsonPolicy)
	if err != nil {
		log.Fatalf("MPutWithPolicy 上传文件失败，错误信息为：%s\n", req.DumpResponse(true))
	}
	log.Printf("MPutWithPolicy 上传回调成功")
	_ = req.DeleteFile(KeyName)

	// 异步分片上传回调
	err = req.AsyncMPutWithPolicy(FilePath, KeyName, MimeType, queryPolicy)
	if err != nil {
		log.Fatalf("AsyncMPutWithPolicy 上传文件失败，错误信息为：%s\n", req.DumpResponse(true))
	}
	log.Printf("AsyncMPutWithPolicy 上传回调成功")
	_ = req.DeleteFile(KeyName)

	// 异步分片并发上传回调
	jobs := 20 // 并发数为 20
	err = req.AsyncUploadWithPolicy(FilePath, KeyName, MimeType, jobs, queryPolicy)
	if err != nil {
		log.Fatalf("AsyncUploadWithPolicy 上传回调失败，错误信息为：%s\n", req.DumpResponse(true))
	}
	log.Printf("AsyncUploadWithPolicy 上传回调成功")
	_ = req.DeleteFile(KeyName)
}
