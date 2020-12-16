package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"os"
)

const (
	FilePath  		= "1.jpg"
	ConfigFile		= "config.json"
	KeyName   		= "1.jpg"
	SaveAsKeyName 	= "pictureSaveAs.jpg"
)

func main() {
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	// 使用 iop 上传 （缩放为原图50%，并持久化为 SaveAsKeyName）
	iopCmdString := "iopcmd=thumbnail&type=1&scale=50|saveAs=" + SaveAsKeyName
	err = req.PutFileWithIopString(FilePath, KeyName, "", iopCmdString)
	if err != nil {
		log.Fatalf("上传文件失败，错误信息为：%s\n", req.DumpResponse(true))
	}
	log.Printf("文件上传成功！")


	err = req.HeadFile(KeyName)
	if err != nil {
		log.Fatalf("error: %s ; \n response: %s", err.Error(), string(req.DumpResponse(true)))
	}
	log.Printf("FileInfo: %s", req.LastResponseHeader)


	file, err := os.OpenFile("111.jpg", os.O_CREATE|os.O_WRONLY, 0755)
	defer file.Close()
	if err != nil {
		log.Fatalf("创建文件失败，错误信息为：%s", err.Error())
	}

	//指定iop字符串执行下载iop （缩放为原图50%）
	iopCmdString = "iopcmd=thumbnail&type=1&scale=50"
	err = req.DownloadFileWithIopString(file, KeyName, iopCmdString)
	if err != nil {
		log.Fatalf("下载文件出错，出错信息为：%s", err.Error())
	}
	log.Printf("文件下载成功！")

}
