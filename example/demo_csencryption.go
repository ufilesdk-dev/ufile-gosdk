package main

import (
	"bytes"
	"crypto/md5"
	"io/ioutil"
	"log"
	"os"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"github.com/ufilesdk-dev/ufile-gosdk/example/helper"
)

const (
	SmallFilePath = "./FakeSmallFile.txt"
	LargeFilePath = "./FakeLargeFile.txt"
	configFile    = "config.json"

	remoteFileKey = "ciphertext.txt"
	saveAsName    = "./download.txt"
)

//PutEncryptedFile 文件客户端加密上传
//func (u *UFileRequest) PutEncryptedFile(filePath, keyName, mimeType string) error
//进行客户端加密上传时，需要用户提供加解密密钥，详情见配置文件相关文档
//本SDK支持加密算法AES-CTR，如有其它加密算法需求，需自行实现加解密方法
//mimeType 如果为空的，会调用 net/http 里面的 DetectContentType 进行检测。
//keyName 表示传到 ufile 的文件名。

//DownloadEncryptedFile 文件客户端加密下载
//func (u *UFileRequest) DownloadEncryptedFile(writer io.Writer, keyName string) error
//进行客户端加密下载时，需要用户提供加解密密钥，详情见配置文件相关文档

//MPutEncryptedFile 加密并同步分片上传一个文件
//func (u *UFileRequest) MPutEncryptedFile(filePath, keyName, mimeType string) error
//filePath 是本地文件所在的路径，内部会自动对文件进行加密和分片上传，上传的方式是同步一片一片的加密再上传。
//mimeType 如果为空的话，会调用 net/http 里面的 DetectContentType 进行检测。
//keyName 表示传到 ufile 的文件名。
//大于 100M 的加密文件推荐使用本接口上传。

//AsyncMPutEncryptedFile 加密并异步分片上传一个文件
//func (u *UFileRequest) AsyncMPutEncryptedFile(filePath, keyName, mimeType string) error
//filePath 是本地文件所在的路径，内部会自动对文件进行加密和分片上传，上传的方式是使用异步的方式同时加密并传多个分片的块。
//mimeType 如果为空的话，会调用 net/http 里面的 DetectContentType 进行检测。
//keyName 表示传到 ufile 的文件名。
//大于 100M 的文件推荐使用本接口上传。
//同时并发上传的分片数量为10

//DownloadLargeEncryptedFile 客户端加密下载接口
//func (u *UFileRequest) DownloadLargeEncryptedFile(writer io.Writer, keyName string) error
//对下载大文件比较友好；支持流式下载
//进行客户端加密下载时，需要用户提供加解密密钥，详情见配置文件相关文档

func PutEncryptedFileExample() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在加密上传文件。。。。")

	//文件客户端加密上传
	err = req.PutEncryptedFile(SmallFilePath, remoteFileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	log.Println("正在加密下载文件。。。。")
	file, err := os.Create(saveAsName)
	if err != nil {
		log.Println("创建文件失败，错误信息为：", err.Error())
		return
	}
	defer file.Close()

	err = req.DownloadEncryptedFile(file, remoteFileKey)
	if err != nil {
		log.Println("下载文件出错，出错信息为：", err.Error())
	}

	ok, err := compareFileMd5(SmallFilePath, saveAsName)
	if err != nil {
		log.Println("客户端加密：比较文件Md5出错，出错信息为：", err.Error())
	}
	if !ok {
		log.Println("客户端加密：比较文件Md5失败")
	} else {
		log.Println("客户端加密：文件客户端加密上传下载成功")
	}
	log.Println()
}

func MPutEncryptedFileExample() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在同步分片上传加密文件。。。。")

	//文件客户端加密同步分片上传
	err = req.MPutEncryptedFile(LargeFilePath, remoteFileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("同步分片上传加密文件成功。")

	log.Println("正在加密下载大文件。。。。")
	file, err := os.Create(saveAsName)
	if err != nil {
		log.Println("创建文件失败，错误信息为：", err.Error())
		return
	}
	defer file.Close()

	err = req.DownloadLargeEncryptedFile(file, remoteFileKey)
	if err != nil {
		log.Println("加密下载大文件出错，出错信息为：", err.Error())
	}

	ok, err := compareFileMd5(LargeFilePath, saveAsName)
	if err != nil {
		log.Println("客户端加密：比较文件Md5出错，出错信息为：", err.Error())
	}
	if !ok {
		log.Println("客户端加密：比较文件Md5失败")
	} else {
		log.Println("客户端加密：文件客户端加密同步分片上传下载成功")
	}
	log.Println()

}

func AsyncMPutEncryptedFile() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在异步分片上传加密文件。。。。")

	//异步分片上传加密文件
	err = req.AsyncMPutEncryptedFile(LargeFilePath, remoteFileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("异步分片上传加密文件成功。")

	log.Println("正在加密下载大文件。。。。")
	file, err := os.Create(saveAsName)
	if err != nil {
		log.Println("创建文件失败，错误信息为：", err.Error())
		return
	}
	defer file.Close()

	err = req.DownloadLargeEncryptedFile(file, remoteFileKey)
	if err != nil {
		log.Println("下载大文件出错，出错信息为：", err.Error())
	}

	ok, err := compareFileMd5(LargeFilePath, saveAsName)
	if err != nil {
		log.Println("客户端加密：比较文件Md5出错，出错信息为：", err.Error())
	}
	if !ok {
		log.Println("客户端加密：比较文件Md5失败")
	} else {
		log.Println("客户端加密：文件客户端加密异步分片上传下载成功")
	}
}
func main() {
	//如果上传文件不存在，就生成文件以测试
	if _, err := os.Stat(SmallFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(SmallFilePath, helper.FakeSmallFileSize)
		log.Println("测试小文件已生成")
	}
	if _, err := os.Stat(LargeFilePath); os.IsNotExist(err) {
		helper.GenerateFakefile(LargeFilePath, helper.FakeBigFileSize)
		log.Println("测试大文件已生成")

	}

	PutEncryptedFileExample()
	MPutEncryptedFileExample()
	AsyncMPutEncryptedFile()
}

func compareFileMd5(uploadFile, saveAsName string) (bool, error) {
	beforePutFile, err := ioutil.ReadFile(uploadFile)
	if err != nil {
		return false, err
	}

	bMd5 := md5.Sum(beforePutFile)

	afterPutFile, err := ioutil.ReadFile(saveAsName)
	if err != nil {
		return false, err
	}
	aMd5 := md5.Sum(afterPutFile)

	return bytes.Equal(aMd5[:], bMd5[:]), nil
}
