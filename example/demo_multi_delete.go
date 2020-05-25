package main

import (
	"fmt"
	"log"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	configTestFile       = "config.json"
	remoteDeleteFile1Key = "test1-file"
	remoteDeleteFile2Key = "test2-file"
	remoteDeleteFile3Key = "test3-file"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(configTestFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	//写入需要进行删除的文件
	deleteArray := []ufsdk.UFileFile{
		ufsdk.UFileFile{Key: remoteDeleteFile1Key},
		ufsdk.UFileFile{Key: remoteDeleteFile2Key},
		ufsdk.UFileFile{Key: remoteDeleteFile3Key},
	}

	//调用批量删除接口 deleteArray:需要进行删除的文件列表，最大长度为1000
	//选择批量删除的结果返回类型
	//false:ufile返回的消息体中会包含所有删除Object的结果。默认采用详细模式。
	//true: ufile返回的消息体中只包含删除过程中出错的Object结果。如果所有删除都成功，则没有消息体。
	//如果您填写的文件列表，删除全部失败（包括全部404），则会报错500
	res, err := req.MultiDeleteFile(deleteArray, false)
	if err != nil {
		fmt.Println("批量删除失败,失败原因:", err.Error())
		return
	}

	// 打印object的删除标记。
	for _, file := range res.Objects {
		fmt.Println("key:", file.Key)
	}
	//获取文件列表
	req, err = ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在获取文件列表。。。。")
	list, err := req.PrefixFileList("", "", 10)
	if err != nil {
		log.Println("获取文件列表失败，错误信息为：", err.Error())
		return
	}
	log.Printf("获取文件列表返回的信息是：\n%s\n", list)
}
