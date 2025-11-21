package main

import (
	"log"
	"strconv"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	FilePath      = "./random_file.bin"
	ConfigFile    = "config.json"
	srcFileKey    = "srcFile.txt"
	dstFileKey    = "dstFile.txt"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	srcFileBucket := req.BucketName
	//1、上传文件
	log.Println("正在上传文件。。。。")
	err = req.MPut(FilePath, srcFileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	//2、获取文件信息
	log.Println("正在获取文件的基本信息。。。。")
	err = req.HeadFile(srcFileKey)
	if err != nil {
		log.Println("查询文件信息失败，具体错误详情：", err.Error())
		return
	}
	log.Println("查询文件信息成功。")
	fileSize, err :=  strconv.ParseInt(req.LastResponseHeader.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Fatalln("获取文件大小失败，失败原因：", err.Error())
		return
	}
	srcEtag := req.LastResponseHeader.Get("Etag")
	log.Println("文件Etag为：", srcEtag)

	//3、执行拷贝到想要rename的新文件
	if fileSize > 100 * 1024 * 1024 {
		log.Println("正在执行大文件分片拷贝重命名。。。。")
		err = req.MCopy(dstFileKey, srcFileBucket, srcFileKey, fileSize)
	} else {
		log.Println("正在执行普通文件拷贝重命名。。。。")
		err = req.Copy(dstFileKey, srcFileBucket, srcFileKey)
	}
	if err != nil {
		log.Fatalln("文件拷贝重命名失败，失败原因：", err.Error())
		return
	}
	log.Println("文件拷贝重命名成功。")

	//4、校验新文件是否存在且和源文件一致
	log.Println("正在校验新文件。。。。")
	err = req.HeadFile(dstFileKey)
	if err != nil {
		log.Fatalln("新文件校验失败，失败原因：", err.Error())
		return
	}
	dstEtag := req.LastResponseHeader.Get("Etag")
	if srcEtag != dstEtag {
		log.Fatalln("新文件校验失败，源文件Etag与新文件Etag不一致，源文件Etag：", srcEtag, "新文件Etag：", dstEtag)
		return
	}
	log.Println("新文件校验成功，Etag为：", dstEtag)

	//5、根据业务需要删除源文件
	log.Println("正在删除源文件。。。。")
	err = req.DeleteFile(srcFileKey)
	if err != nil {
		log.Fatalln("源文件删除失败，失败原因：", err.Error())
		return
	}
	log.Println("源文件删除成功。")
}
