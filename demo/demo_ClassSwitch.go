package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"net/http"
	"time"
)

const (
	ConfigFile = "./config.json"
	FilePath = "test.txt"
	KeyName = "test.txt"
)

func main() {
	//STANDARD2IA()
	//IA2ARCHIVE()
	STANDARD2ARCHIVE()
}

// 标准转低频
func STANDARD2IA() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	header := make(http.Header)
	header.Add("X-Ufile-Storage-Class", "STANDARD")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}

	// 上传标准存储类型文件
	err = req.PutFile(FilePath, KeyName, "")
	if err != nil {
		panic(err.Error())
	}
	// 转换类型
	err = req.ClassSwitch(KeyName, "IA")
	if err != nil {
		log.Fatalf("文件存储类型转换失败，错误信息为：%s", err.Error())
	}
	log.Printf("文件存储类型转换成功")

	// 获取文件基本信息
	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("获取文件基本信息失败，错误信息为：%s", err.Error())
	}
	log.Printf("文件基本信息为：%s", req.LastResponseHeader)
}

// 低频 转 归档
func IA2ARCHIVE() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	header := make(http.Header)
	header.Add("X-Ufile-Storage-Class", "IA")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}

	// 上传标准存储类型文件
	err = req.PutFile(FilePath, KeyName, "")
	if err != nil {
		panic(err.Error())
	}

	// 转换类型
	err = req.ClassSwitch(KeyName, "ARCHIVE")
	if err != nil {
		log.Fatalf("文件存储类型转换失败，错误信息为：%s", err.Error())
	}
	log.Printf("文件存储类型转换成功")

	// 获取文件基本信息
	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("获取文件基本信息失败，错误信息为：%s", err.Error())
	}
	log.Printf("文件基本信息为：%s", req.LastResponseHeader)

	// 验证转为归档后， 文件是否还可以下载
	err = req.Download(req.GetPrivateURL(KeyName, time.Second*60))
	if err != nil {
		log.Fatalf("下载失败，错误信息为：%s", err.Error())
	}
	log.Printf("文件内容为：%s", string(req.LastResponseBody))
}

// 标准转低频
func STANDARD2ARCHIVE() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	header := make(http.Header)
	header.Add("X-Ufile-Storage-Class", "STANDARD")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}

	// 上传标准存储类型文件
	err = req.PutFile(FilePath, KeyName, "")
	if err != nil {
		panic(err.Error())
	}
	// 转换类型
	err = req.ClassSwitch(KeyName, "ARCHIVE")
	if err != nil {
		log.Fatalf("文件存储类型转换失败，错误信息为：%s", err.Error())
	}
	log.Printf("文件存储类型转换成功")

	// 获取文件基本信息
	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("获取文件基本信息失败，错误信息为：%s", err.Error())
	}
	log.Printf("文件基本信息为：%s", req.LastResponseHeader)

	// 验证转为归档后， 文件是否还可以下载
	for {
		err = req.Download(req.GetPrivateURL(KeyName, time.Second*60))
		if err != nil {
			log.Fatalf("下载失败，错误信息为：%s", err.Error())
		}
		log.Printf("文件内容为：%s", string(req.LastResponseBody))
		time.Sleep(2*time.Second)
	}
}


