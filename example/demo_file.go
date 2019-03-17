package main

import (
	"log"
	"os"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"github.com/ufilesdk-dev/ufile-gosdk/example/helper"
)

const (
	putUpload = iota
	postUpload
	mput
	asyncmput
)

func main() {
	log.SetFlags(log.Lshortfile)
	if _, err := os.Stat(helper.FakeSmallFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(helper.FakeSmallFilePath, helper.FakeSmallFileSize)
	}
	if _, err := os.Stat(helper.FakeBigFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(helper.FakeBigFilePath, helper.FakeBigFileSize)
	}
	config, err := ufsdk.LoadConfig(helper.ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	var fileKey string
	fileKey = helper.GenerateUniqKey()
	scheduleUploadhelper(helper.FakeSmallFilePath, fileKey, putUpload, req)
	fileKey = helper.GenerateUniqKey()
	scheduleUploadhelper(helper.FakeSmallFilePath, fileKey, postUpload, req)

	fileKey = helper.GenerateUniqKey()
	scheduleUploadhelper(helper.FakeBigFilePath, fileKey, mput, req)
	fileKey = helper.GenerateUniqKey()
	scheduleUploadhelper(helper.FakeBigFilePath, fileKey, asyncmput, req)
}

func scheduleUploadhelper(filePath, keyName string, uploadType int, req *ufsdk.UFileRequest) {
	log.Println("上传的文件 Key 为：", keyName)
	var err error
	switch uploadType {
	case putUpload:
		log.Println("正在使用PUT接口上传文件...")
		err = req.PutFile(filePath, keyName, "")
		break
	case postUpload:
		log.Println("正在使用 POST 接口上传文件...")
		err = req.PostFile(filePath, keyName, "")
	case mput:
		log.Println("正在使用同步分片上传接口上传文件...")
		err = req.MPut(filePath, keyName, "")
	case asyncmput:
		log.Println("正在使用异步分片上传接口上传文件...")
		err = req.AsyncMPut(filePath, keyName, "")
	}
	if err != nil {
		log.Println("文件上传失败!!，错误信息为：", err.Error())
		//如果 err 给出的提示信息不够，你可 dump 整个 response 出来查看 http 的返回。
		log.Printf("%s\n", req.DumpResponse(true))
		return
	}
	log.Println("文件上传成功!!")
	log.Println("公有空间文件下载 URL 是：", req.GetPublicURL(keyName))
	log.Println("私有空间文件下载 URL 是：", req.GetPrivateURL(keyName, 24*60*60)) //过期时间为一天

	log.Println("正在获取文件的基本信息。")
	err = req.HeadFile(keyName)
	if err != nil {
		log.Println("查询文件信息失败，具体错误详情：", err.Error())
		return
	}

	log.Println("正在秒传文件...")
	err = req.UploadHit(filePath, keyName)
	if err != nil {
		log.Println("文件秒传失败，错误信息为：", err.Error())
	} else {
		log.Printf("秒传文件返回的信息是：%s\n", req.LastResponseBody)
	}

	log.Println("正在获取文件列表...")
	list, err := req.PrefixFileList(keyName, "", 10)
	if err != nil {
		log.Println("获取文件列表失败，错误信息为：", err.Error())
		return
	}
	log.Printf("获取文件列表返回的信息是：\n%s\n", list)

	log.Println("正在删除刚刚上传的文件")
	err = req.DeleteFile(keyName)
	if err != nil {
		log.Println("删除文件失败，错误信息为：", err.Error())
		return
	}
	log.Println("删除文件成功")
}
