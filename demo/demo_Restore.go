package main

import (
	"fmt"
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
	restoreSmallFile()
	//restoreBigFile()
}

// 小文件解冻
func restoreSmallFile() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	header := make(http.Header)
	header.Add("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}

	// 上传一个归档存储文件
	err = req.PutFile(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}
	log.Printf("归档存储文件上传成功")

	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("获取文件基本信息失败，错误信息为：%s", err.Error())
	}
	log.Printf("获取文件基本信息成功，内容为：%s", req.LastResponseHeader)

	_, _ = fmt.Scanln()

	// 解冻一个归档存储文件
	err = req.Restore(KeyName)
	if err != nil {
		log.Println("解冻文件失败:", string(req.DumpResponse(true)))
		log.Fatalf("解冻文件 %s 失败，错误信息为：%s", KeyName, err.Error())
	}
	time.Sleep(time.Second * 10)

	// 查看解冻是否成功
	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("获取文件基本信息失败，错误信息为：%s", err.Error())
	}
	if _, ok := req.LastResponseHeader["X-Ufile-Restore"]; ok {
		log.Printf("文件解冻成功，内容为：%s", req.LastResponseHeader.Get("X-Ufile-Restore"))
		log.Printf("文件解冻成功，内容为：%s", req.LastResponseHeader)
	}

}

// 大文件解冻
func restoreBigFile() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	header := make(http.Header)
	header.Set("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}

	// 上传一个归档存储文件
	err = req.MPut("installer.exe", "installer", "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 解冻一个归档存储文件
	err = req.Restore("installer")
	if err != nil {
		log.Fatalf("解冻文件 %s 失败，错误信息为：%s", KeyName, err.Error())
	}
	time.Sleep(time.Second * 10)

	// 查看解冻是否成功
	err = req.HeadFile("installer")
	if err != nil {
		log.Fatalf("获取文件基本信息失败，错误信息为：%s", err.Error())
	}
	if _, ok := req.LastResponseHeader["X-Ufile-Restore"]; ok {
		log.Printf("文件解冻成功，内容为：%s", req.LastResponseHeader.Get("X-Ufile-Restore"))
	}

}

