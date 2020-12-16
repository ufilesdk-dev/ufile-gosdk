package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	ConfigFile = "./config.json"
	FilePath = "test.txt"
	KeyName = "test.txt"
	DeletedDir = "test/"
	Prefix = ""
	Delimiter = ""
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

	// 删除一个文件
	err = req.PutFile(FilePath,  KeyName, "")
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	log.Println("文件上传成功!!")

	err = req.DeleteFile(KeyName)
	if err != nil {
		log.Println(string(req.DumpResponse(true)))
		log.Fatalf("删除文件 %s 失败，错误信息为：%s", KeyName, err.Error())
	}
	log.Println("删除文件成功")

	// 删除指定目录下的文件(非递归)
	for ;; {
		list, err := req.ListObjects(DeletedDir, "", Delimiter, 0)
		if err != nil {
			log.Fatalf("获取文件列表失败，错误信息为：%s", err.Error())
		}
		if len(list.Contents) == 0 {
			break
		}
		log.Printf("获取文件列表成功, 文件列表长度为 %d", len(list.Contents))
		for _, object := range list.Contents {
			log.Printf("%s 删除成功", object.Key)
			err = req.DeleteFile(object.Key)
			if err != nil {
				log.Fatalf("删除文件 %s 失败，错误信息为：%s", object.Key, err.Error())
			}
		}
	}
	log.Printf("成功删除目录 %s 下的文件", DeletedDir)

	// 删除带有 Prefix 的文件(递归, DeletedDir 为空将删除 Bucket 所有文件)
	for ;; {
		list, err := req.ListObjects(Prefix, "", Delimiter, 0)
		if err != nil {
			log.Fatalf("获取文件列表失败，错误信息为：%s", err.Error())
		}
		if len(list.Contents) == 0 {
			break
		}
		log.Printf("获取文件列表成功, 文件列表长度为 %d", len(list.Contents))
		for _, object := range list.Contents {
			log.Printf("%s 删除成功", object.Key)
			err = req.DeleteFile(object.Key)
			if err != nil {
				log.Fatalf("删除文件 %s 失败，错误信息为：%s", object.Key, err.Error())
			}
		}
	}
	log.Printf("成功删除前缀为 %s 的文件", Prefix)

}
