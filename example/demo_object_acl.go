package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	FilePath      = "./FakeSmallFile.txt"
	ConfigFile    = "config.json"
	FileKey       = "test_acl.txt"
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

	//1、上传文件
	log.Println("正在上传文件。。。。")
	err = req.MPut(FilePath, FileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	//2、设置acl, default：继承Bucket ACL，public-read：公共读
	log.Println("正在设置文件acl。。。。")
	acl := "public-read"
	err = req.PutObjectAcl(FileKey, acl)
	if err != nil {
		log.Fatalln("设置文件acl失败，失败原因：", err.Error())
		return
	}
	log.Println("设置文件acl成功。")

	//3、获取acl
	log.Println("正在获取文件acl。。。。")
	acl, err = req.GetObjectAcl(FileKey)
	if err != nil {
		log.Fatalln("获取文件acl失败，失败原因：", err.Error())
		return
	}
	log.Println("获取文件acl:%s", acl)
}
