package main

import (
	"bytes"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"io"
	"os"
	"time"
)

const (
	ConfigFile = "./config.json"
	FilePath = "mongo.tgz"
	KeyName = "mongo.tgz"
	MimeType = ""
)

func checkErr (err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	test4(test3())
}

// 完整上传流程
func test1() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Open(FilePath)
	checkErr(err)
	defer file.Close()

	// 初始化分片
	state, err := req.InitiateMultipartUpload(KeyName, MimeType)
	checkErr(err)
	// 逐个上传分片， 若出错则终止上传
	chunk := make([]byte, state.BlkSize)
	var pos int
	for {
		bytesRead, fileErr := file.Read(chunk)
		if fileErr == io.EOF || bytesRead == 0 { //后面直接读到了结尾
			break
		}
		buf := bytes.NewBuffer(chunk[:bytesRead])
		err := req.UploadPart(buf, state, pos)
		if err != nil {
			checkErr(err)
			req.AbortMultipartUpload(state)
		}
		pos++
	}
	// 完成分片
	err = req.FinishMultipartUpload(state)
	checkErr(err)
}

//  测试断点续传
func test2() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Open(FilePath)
	checkErr(err)
	defer file.Close()

	// 初始化分片
	state, err := req.InitiateMultipartUpload(KeyName, MimeType)
	checkErr(err)
	// 逐个上传分片， 若出错则终止上传
	chunk := make([]byte, state.BlkSize)
	var pos int
	for {
		bytesRead, fileErr := file.Read(chunk)
		if fileErr == io.EOF || bytesRead == 0 { //后面直接读到了结尾
			break
		}
		buf := bytes.NewBuffer(chunk[:bytesRead])
		err := req.UploadPart(buf, state, pos)
		if err != nil {
			checkErr(err)
		}
		pos++
		if pos == 10 {
			break
		}
	}

	time.Sleep(time.Second * 60)
	time.Sleep(time.Second * 60 * 10)

	pos = 10
	for {
		bytesRead, fileErr := file.Read(chunk)
		if fileErr == io.EOF || bytesRead == 0 { //后面直接读到了结尾
			break
		}
		buf := bytes.NewBuffer(chunk[:bytesRead])
		err := req.UploadPart(buf, state, pos)
		if err != nil {
			checkErr(err)
		}
		pos++
	}

	time.Sleep(time.Second * 60)
	// 完成分片
	err = req.FinishMultipartUpload(state)
	checkErr(err)
}


//  测试断点续传
func test3() *ufsdk.MultipartState {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Open(FilePath)
	checkErr(err)
	defer file.Close()

	// 初始化分片
	state, err := req.InitiateMultipartUpload(KeyName, MimeType)
	checkErr(err)
	// 逐个上传分片， 若出错则终止上传
	chunk := make([]byte, state.BlkSize)
	var pos int
	for {
		bytesRead, fileErr := file.Read(chunk)
		if fileErr == io.EOF || bytesRead == 0 { //后面直接读到了结尾
			break
		}
		buf := bytes.NewBuffer(chunk[:bytesRead])
		err := req.UploadPart(buf, state, pos)
		if err != nil {
			checkErr(err)
		}
		pos++
		if pos == 20 {
			break
		}
	}
	return state
}

//  测试断点续传
func test4(state *ufsdk.MultipartState) {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Open(FilePath)
	checkErr(err)
	defer file.Close()

	// 逐个上传分片， 若出错则终止上传
	chunk := make([]byte, state.BlkSize)
	var pos int
	for {
		bytesRead, fileErr := file.Read(chunk)
		if fileErr == io.EOF || bytesRead == 0 { //后面直接读到了结尾
			break
		}
		buf := bytes.NewBuffer(chunk[:bytesRead])
		err := req.UploadPart(buf, state, pos)
		if err != nil {
			checkErr(err)
		}
		pos++
	}
	// 完成分片
	err = req.FinishMultipartUpload(state)
	checkErr(err)
}