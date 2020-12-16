package main

import (
	"fmt"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

func main() {

	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig("./config.json")
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}


	var Prefix string
	var Marker string
	var Limit int
	var fileList ufsdk.FileListResponse

	// 获取指定前缀的所有文件(非递归)
	Prefix = "prefix/"
	Marker = "sa"
	Limit = 0
	for ;; {
		list, err := req.PrefixFileList(Prefix, Marker , Limit)
		if err != nil {
			log.Println(string(req.DumpResponse(true)))
			log.Fatalf("前缀列表查询失败，错误信息为：%s", err.Error())
		}
		log.Printf("前缀列表查询成功, 文件列表长度为 %d", len(list.DataSet))
		log.Printf("前缀列表查询成功, list.NextMarker %s", list.NextMarker)
		fileList.DataSet = append(fileList.DataSet, list.DataSet...)
		if len(list.DataSet) == 0 || len(list.NextMarker) <= 0{
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("前缀列表查询成功, 文件列表长度为 %d", len(fileList.DataSet))
	_, _ = fmt.Scanln()
	
	
	
	// 列举指定数量文件
	fileList.DataSet = fileList.DataSet[:0]
	Limit = 100
	Marker = ""
	for {
		list, err := req.PrefixFileList("", Marker, 0)
		if err != nil {
			log.Printf("前缀列表查询失败，错误信息为：%s", string(req.DumpResponse(true)))
			log.Fatalf("前缀列表查询失败，错误信息为：%s", err.Error())
		}
		log.Printf("前缀列表查询成功, 文件列表长度为 %d", len(list.DataSet))
		log.Printf("前缀列表查询成功, list.NextMarker %s", list.NextMarker)
		log.Printf("前缀列表查询成功, Limit %d", Limit)
		if len(list.DataSet) >= Limit {
			fileList.DataSet = append(fileList.DataSet, list.DataSet[:Limit]...)
			break
		}
		fileList.DataSet = append(fileList.DataSet, list.DataSet...)
		Limit = Limit - len(list.DataSet)
		if len(list.NextMarker) <= 0 {
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("前缀列表查询成功, 文件列表长度为 %d", len(fileList.DataSet))
	_, _ = fmt.Scanln()

	// 列举指定 Marker 之后的所有文件(递归)
	Marker = "prefix/3"
	fileList.DataSet = fileList.DataSet[:0]
	for ;; {
		list, err := req.PrefixFileList("", Marker, 0)
		if err != nil {
			log.Fatalf("前缀列表查询失败，错误信息为：%s", err.Error())
		}
		log.Printf("前缀列表查询成功, 文件列表长度为 %d", len(list.DataSet))
		log.Printf("前缀列表查询成功, list.NextMarker %s", list.NextMarker)
		fileList.DataSet = append(fileList.DataSet, list.DataSet...)
		if len(list.NextMarker) <= 0{
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("前缀列表查询成功, 文件列表长度为 %d", len(fileList.DataSet))
	_, _ = fmt.Scanln()

	// 列举指定前缀、指定Marker之后、指定最大数量的文件列表
	Prefix = "prefix"
	Marker = "/prefix"
	Limit = 200
	fileList.DataSet = fileList.DataSet[:0]
	for ;; {
		list, err := req.PrefixFileList(Prefix, Marker, 0)
		if err != nil {
			log.Fatalf("前缀列表查询失败，错误信息为：%s", err.Error())
		}
		log.Printf("前缀列表查询成功, 文件列表长度为 %d", len(list.DataSet))
		log.Printf("前缀列表查询成功, list.NextMarker %s", list.NextMarker)
		if len(list.DataSet) >= Limit {
			fileList.DataSet = append(fileList.DataSet, list.DataSet[:Limit]...)
			break
		}
		fileList.DataSet = append(fileList.DataSet, list.DataSet...)
		Limit = Limit - len(list.DataSet)
		if len(list.NextMarker) <= 0 {
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("前缀列表查询成功, 文件列表长度为 %d", len(fileList.DataSet))
}

// 获取指定前缀文件列表(非递归)
//func getFileListByPrefix(configPath, prefix string) (objectInfos []ufsdk.ObjectInfo, err error){
//	// 加载配置，创建请求
//	config, err := ufsdk.LoadConfig(configPath)
//	if err != nil {
//		panic(err.Error())
//	}
//	req, err := ufsdk.NewFileRequest(config, nil)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	for ;; {
//		list, err := req.PrefixFileList(prefix, "", Delimiter, 0)
//		if err != nil {
//			return nil, err
//		}
//		if len(list.DataSet) == 0 {
//			return objectInfos, nil
//		}
//		objectInfos = append(objectInfos, list.DataSet...)
//	}
//}


