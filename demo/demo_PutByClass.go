package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"net/http"
)

const (
	ConfigFile = "./config.json"
	FilePath = "test.txt"
	KeyName = "test.txt"
	mPutFilePath = "mongo.tgz"
	ArchiveEKeyName = "Ar_test.txt"
	IAKeyName = "IA_test.txt"
	StanEKeyName = "Stan_test.txt"
)

func main() {
	putByClass()
}

// 上传归档存储文件
func test1(){
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	header := make(http.Header)
	// 上传一个归档存储类型文件 ( 首次使用 X-Ufile-Storage-Class 需使用 http.Header.Add 方法添加)
	header.Add("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.PutFile(FilePath, ArchiveEKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个低频存储类型文件 (修改Header已有字段使用 Set )
	req.RequestHeader.Set("X-Ufile-Storage-Class", "IA")
	err = req.PutFile(FilePath, IAKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个标准存储类型文件
	req.RequestHeader.Set("X-Ufile-Storage-Class", "STANDARD")
	err = req.PutFile(FilePath, StanEKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}
}

func test2(){
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	// 使用 NewFileRequest 构造的 req 无 Header，需先 make
	req.RequestHeader = make(http.Header)


	// 上传一个归档存储类型文件 ( 首次使用 X-Ufile-Storage-Class 需使用 http.Header.Add 方法添加)
	req.RequestHeader.Add("X-Ufile-Storage-Class", "ARCHIVE")
	err = req.PostFile(FilePath, ArchiveEKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个低频存储类型文件 (修改Header已有字段使用 Set )
	req.RequestHeader.Set("X-Ufile-Storage-Class", "IA")
	err = req.PostFile(FilePath, IAKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个标准存储类型文件
	req.RequestHeader.Set("X-Ufile-Storage-Class", "STANDARD")
	err = req.PostFile(FilePath, StanEKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}
}


// PutFile 指定存储类型上传
func putByClass() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}


	header := make(http.Header)
	// 上传一个归档存储类型文件
	header.Add("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.PutFile(FilePath, ArchiveEKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个低频存储类型文件
	req.RequestHeader.Set("X-Ufile-Storage-Class", "IA")
	err = req.PutFile(FilePath, IAKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个标准存储类型文件
	req.RequestHeader.Set("X-Ufile-Storage-Class", "STANDARD")
	err = req.PutFile(FilePath, StanEKeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}
}

// PostFile 指定存储类型上传
func postByClass() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	header := make(http.Header)
	// 上传一个归档存储类型文件
	header.Set("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.PostFile(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个低频存储类型文件
	header.Set("X-Ufile-Storage-Class", "IA")
	req, err = ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.PostFile(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个标准存储类型文件
	header.Set("X-Ufile-Storage-Class", "IA")
	req, err = ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.PostFile(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}
}

// MP指定存储类型上传
func mPutByClass() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	header := make(http.Header)
	// 上传一个归档存储类型文件
	header.Set("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.MPut(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个低频存储类型文件
	header.Set("X-Ufile-Storage-Class", "IA")
	req, err = ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.MPut(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个标准存储类型文件
	header.Set("X-Ufile-Storage-Class", "IA")
	req, err = ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.MPut(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}
}

// FinishMultipartUpload 指定存储类型上传
func finishByClass() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	header := make(http.Header)
	// 上传一个归档存储类型文件
	header.Set("X-Ufile-Storage-Class", "ARCHIVE")
	req, err := ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}


	err = req.MPut(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	err = req.AsyncMPut(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	err = req.AsyncUpload(FilePath, KeyName, "", 10)
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个低频存储类型文件
	header.Set("X-Ufile-Storage-Class", "IA")
	req.RequestHeader = header
	err = req.MPut(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}

	// 上传一个标准存储类型文件
	header.Set("X-Ufile-Storage-Class", "IA")
	req, err = ufsdk.NewFileRequestWithHeader(config, header, nil)
	if err != nil {
		panic(err.Error())
	}
	err = req.MPut(FilePath, KeyName, "")
	if err != nil {
		log.Fatalf("文件上传失败，失败原因：%s", err.Error())
	}
}