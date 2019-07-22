package main

import (
	"net/http"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	uploadFile    = "./FakeBigFile.txt"
	configFile    = "config.json"
	remoteStFileKey = "test_standard1.txt"
	remoteIaFileKey = "test_ia2.txt"
	remoteArFileKey = "test_archive3.txt"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}

	//存储类型，目前支持的类型分别是标准:"STANDARD"、低频:"IA"、冷存:"ARCHIVE"
	header := make(http.Header)

	//1、上传标准存储类型文件
	header.Add("X-Ufile-Storage-Class", "STANDARD")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在上传标准存储类型文件。。。。")
	err = req.MPut(uploadFile, remoteStFileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	//2、上传低频存储类型文件
	header.Set("X-Ufile-Storage-Class", "IA")
	req, err = ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在上传低频存储类型文件。。。。")
	err = req.MPut(uploadFile, remoteIaFileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	//3、上传归档存储类型文件
	header.Set("X-Ufile-Storage-Class", "ARCHIVE")
	req, err = ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在上传归档存储类型文件。。。。")
	err = req.MPut(uploadFile, remoteArFileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	//4、解冻归档存储类型文件
	req, err = ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在解冻归档存储类型文件。。。。")
	err = req.Restore(remoteArFileKey)
	if err != nil {
		log.Println("文件解冻失败，失败原因：", err.Error())
		return
	}
	log.Println("文件解冻成功。")

	//5、转换文件存储类型
	req, err = ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在转换归档存储类型文件为低频类型。。。。")
	err = req.ClassSwitch(remoteArFileKey, "IA")
	if err != nil {
		log.Println("文件转换存储类型失败，失败原因：", err.Error())
		return
	}
	log.Println("文件转换存储类型成功。")

	//6、获取文件列表
	req, err = ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在获取文件列表。。。。")
	list, err := req.PrefixFileList("test_", "", 10)
	if err != nil {
		log.Println("获取文件列表失败，错误信息为：", err.Error())
		return
	}
	log.Printf("获取文件列表返回的信息是：\n%s\n", list)
}
