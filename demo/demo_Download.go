package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
	"os"
	"time"
)

const (
	ConfigFile = "./config.json"
	KeyName = "test.txt"
	FilePath = "test1.txt"
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

	// 可以通过 GetPrivateURL 获取下载URL，第二个参数为URL有效时间
	reqUrl := req.GetPrivateURL(KeyName, 10*time.Second)
	log.Println("reqUrl:", reqUrl)
	// Download 要求传入 object 下载URL，object下载数据存储在 req.LastResponseBody 中
	err = req.Download(reqUrl)
	if err != nil {
		log.Println(string(req.DumpResponse(true)))
		log.Fatalf("下载文件失败！返回错误信息为：%s", err.Error())
	}
	log.Printf("下载文件成功！数据为：%s", string(req.LastResponseBody))

	// 保存到本地文件
	file, err := os.OpenFile(FilePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalf("创建文件失败，错误信息为：%s", err.Error())
	}
	_, _ = file.WriteString(string(req.LastResponseBody))
	_ = file.Close()
}
