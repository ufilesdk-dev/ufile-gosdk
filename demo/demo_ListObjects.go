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
	var Delimiter string
	var MaxKeys int
	var fileList ufsdk.ListObjectsResponse

	// 获取指定前缀的所有文件(非递归)
	Prefix = "prefix/2"
	Marker = ""
	for ;; {
		list, err := req.ListObjects(Prefix, Marker , "", 0)
		if err != nil {
			log.Fatalf("获取文件列表失败，错误信息为：%s", err.Error())
		}
		log.Printf("获取文件列表成功, 文件列表长度为 %d", len(list.Contents))
		log.Printf("获取文件列表成功, list.NextMarker %s", list.NextMarker)
		fileList.Contents = append(fileList.Contents, list.Contents...)
		if len(list.Contents) == 0 || len(list.NextMarker) <= 0{
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("获取文件列表成功, 文件列表长度为 %d", len(fileList.Contents))
	_, _ = fmt.Scanln()

	// 列举指定前缀的所有文件(递归)
	Prefix = "prefix/"
	Marker = ""
	fileList.Contents = fileList.Contents[:0]
	for ;; {
		list, err := req.ListObjects(Prefix, Marker, "", 0)
		if err != nil {
			log.Fatalf("获取文件列表失败，错误信息为：%s", err.Error())
		}
		log.Printf("获取文件列表成功, 文件列表长度为 %d", len(list.Contents))
		log.Printf("获取文件列表成功, list.NextMarker %s", list.NextMarker)
		fileList.Contents = append(fileList.Contents, list.Contents...)
		if len(list.Contents) == 0 || len(list.NextMarker) <= 0{
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("获取文件列表成功, 文件列表长度为 %d", len(fileList.Contents))
	_, _ = fmt.Scanln()

	// 列举指定数量文件
	fileList.Contents = fileList.Contents[:0]
	MaxKeys = 100
	Marker = ""
	for {
		list, err := req.ListObjects("", Marker, "", 0)
		if err != nil {
			log.Fatalf("获取文件列表失败，错误信息为：%s", err.Error())
		}
		log.Printf("获取文件列表成功, 文件列表长度为 %d", len(list.Contents))
		log.Printf("获取文件列表成功, list.NextMarker %s", list.NextMarker)
		log.Printf("获取文件列表成功, MaxKeys %d", MaxKeys)
		if len(list.Contents) >= MaxKeys {
			fileList.Contents = append(fileList.Contents, list.Contents[:MaxKeys]...)
			break
		}
		fileList.Contents = append(fileList.Contents, list.Contents...)
		MaxKeys = MaxKeys - len(list.Contents)
		if len(list.NextMarker) <= 0 {
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("获取文件列表成功, 文件列表长度为 %d", len(fileList.Contents))
	_, _ = fmt.Scanln()

	// 列举指定 Marker 之后的所有文件(递归)
	Marker = "prefix/dir"
	fileList.Contents = fileList.Contents[:0]
	for ;; {
		list, err := req.ListObjects("", Marker, "", 0)
		if err != nil {
			log.Fatalf("获取文件列表失败，错误信息为：%s", err.Error())
		}
		log.Printf("获取文件列表成功, 文件列表长度为 %d", len(list.Contents))
		log.Printf("获取文件列表成功, list.NextMarker %s", list.NextMarker)
		fileList.Contents = append(fileList.Contents, list.Contents...)
		if len(list.NextMarker) <= 0{
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("获取文件列表成功, 文件列表长度为 %d", len(fileList.Contents))
	_, _ = fmt.Scanln()

	// 列举指定前缀、指定Marker之后、指定递归方式、指定最大数量的文件列表
	Prefix = "prefix"
	Marker = "prefix"
	Delimiter = "/"
	MaxKeys = 200
	fileList.Contents = fileList.Contents[:0]
	for ;; {
		list, err := req.ListObjects(Prefix, Marker, Delimiter, 0)
		if err != nil {
			log.Fatalf("获取文件列表失败，错误信息为：%s", err.Error())
		}
		log.Printf("获取文件列表成功, 文件列表长度为 %d", len(list.Contents))
		log.Printf("获取文件列表成功, list.NextMarker %s", list.NextMarker)
		if len(list.Contents) >= MaxKeys {
			fileList.Contents = append(fileList.Contents, list.Contents[:MaxKeys]...)
			break
		}
		fileList.Contents = append(fileList.Contents, list.Contents...)
		MaxKeys = MaxKeys - len(list.Contents)
		if len(list.NextMarker) <= 0 {
			break
		}
		Marker = list.NextMarker
	}
	log.Printf("获取文件列表成功, 文件列表长度为 %d", len(fileList.Contents))
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
//		list, err := req.ListObjects(prefix, "", Delimiter, 0)
//		if err != nil {
//			return nil, err
//		}
//		if len(list.Contents) == 0 {
//			return objectInfos, nil
//		}
//		objectInfos = append(objectInfos, list.Contents...)
//	}
//}


