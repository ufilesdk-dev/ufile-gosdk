package main

import (
	"fmt"
	"log"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	configTestFile       = "config.json"
	remoteSwitchFile1Key = "1.test"
	remoteSwitchFile2Key = "2.test"
	remoteSwitchFile3Key = "3.test"
	remoteSwitchFile4Key = "4.test"
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
	//写入需要进行转换的文件
	deleteArray := []ufsdk.SwitchFile{
		ufsdk.SwitchFile{Key: remoteSwitchFile1Key, StorageClass: "IA"},      //标准转低频
		ufsdk.SwitchFile{Key: remoteSwitchFile2Key, StorageClass: "ARCHIVE"}, //标准转归档
		ufsdk.SwitchFile{Key: remoteSwitchFile3Key, StorageClass: "ARCHIVE"}, //低频转归档
		ufsdk.SwitchFile{Key: remoteSwitchFile4Key, StorageClass: "IA"},      //归档低频 ,不可转换，解冻以后转换仍然失败。
	}

	//调用批量转换接口 deleteArray:需要进行转换的文件列表，最大长度为1000
	//选择批量转换的结果返回类型
	//false:ufile返回的消息体中会包含所有转换Object的结果。默认采用详细模式。
	//true: ufile返回的消息体中只包含转换过程中出错的Object结果。如果所有转换都成功，则没有消息体。
	//如果您填写的文件列表，转换全部失败（包括全部404），则会报错500
	res, err := req.MultiClassSwitchFile(deleteArray, false)
	if err != nil {
		fmt.Println("批量转换失败,失败原因:", err.Error())
		return
	}

	// 打印object的转换标记。
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
