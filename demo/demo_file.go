package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"os"
	"time"
)

const (
	ConfigFile = "./config.json"
	FilePath = "test.txt"
	KeyName = "test.txt"
	MimeType = ""
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

	// 创建测试文件
	f, err := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic("创建测试文件失败，失败信息为：" + err.Error())
	}
	defer f.Close()
	_, err = f.WriteString("I am a test file!!!")
	if err != nil {
		panic("创建测试文件失败，失败信息为：" + err.Error())
	}

	// 上传一个文件
	err = req.PutFile(FilePath,  KeyName, MimeType)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("%s\n", err.Error())
	}
	log.Println("文件上传成功!!")
	log.Println("公有空间文件下载 URL 是：", req.GetPublicURL(KeyName))
	log.Println("私有空间文件下载 URL 是：", req.GetPrivateURL(KeyName, 24*60*60 * time.Second)) //过期时间为一天

	// 获取文件的基本信息
	err = req.HeadFile(KeyName)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("查询文件信息失败，具体错误详情：%s", err.Error())
	}
	log.Println("查询文件信息成功，文件基本信息：", req.LastResponseHeader)

	// 下载文件
	err = req.Download(req.GetPrivateURL(KeyName, 24*60*60 * time.Second))
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("下载文件失败！返回错误信息为：%s", err.Error())
	}
	log.Println("下载文件成功！数据为：", string(req.LastResponseBody))

	// 正在秒传文件
	err = req.UploadHit(FilePath, "UploadHit" + KeyName)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Println("文件秒传失败，错误信息为：", err.Error())
	} else {
		log.Println("秒传文件成功")
	}

	// 拷贝文件
	err = req.Copy("Copy" + KeyName, config.BucketName, KeyName)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Println("文件拷贝失败，错误信息为：", err.Error())
	} else {
		log.Println("文件拷贝成功，拷贝文件返回的信息是：", req.LastResponseBody)
	}

	// 获取文件列表
	list, err := req.PrefixFileList(KeyName, "", 10)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("获取文件列表失败，错误信息为：%s", err.Error())
	}
	log.Println("获取文件列表返回的信息是：", list)

	// 获取目录文件列表
	listV2, err := req.ListObjects(KeyName, "", "/", 10)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("获取目录文件列表失败，错误信息为：%s", err.Error())
	}
	log.Println("获取目录文件列表返回的信息是：", listV2)

	// 重命名文件
	err = req.Rename(KeyName, "Rename" + KeyName, "true")
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Println("文件重命名失败，错误信息为：", err.Error())
	} else {
		log.Println("文件重命名成功，重命名文件返回的信息是：\n", req.LastResponseBody)
	}

	// 删除文件
	err = req.DeleteFile("Rename" + KeyName)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("删除文件失败，错误信息为：%s", err.Error())
	}
	log.Println("删除文件成功")

	_ = req.DeleteFile("Rename" + KeyName)
	_ = req.DeleteFile("Copy" + KeyName)
	_ = req.DeleteFile("UploadHit" + KeyName)
}
