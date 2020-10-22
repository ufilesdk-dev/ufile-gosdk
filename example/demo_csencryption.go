package main

import (
	"bytes"
	"crypto/md5"
	"io/ioutil"
	"log"
	"os"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	uploadFile    = "./plaintext.txt"
	configFile    = "config.json"
	remoteFileKey = "ciphertext.txt"
	saveAsName    = "./download.txt"
)

//PutWithCryptoFile 文件客户端加密上传
//func (u *UFileRequest) PutWithEncryptFile(filePath, keyName, mimeType string) error
//进行客户端加密上传时，需要用户提供加解密密钥，详情见配置文件相关文档
//本SDK支持加密算法AES-GCM-NoPadding，如有其它加密算法需求，需自行实现加解密方法
//注意在客户端加密的条件下，ufile暂不支持文件分片上传下载操作。
//mimeType 如果为空的，会调用 net/http 里面的 DetectContentType 进行检测。
//keyName 表示传到 ufile 的文件名。

//DownloadWithDecryptFile 文件客户端加密下载
//func (u *UFileRequest) DownloadWithDecryptFile(writer io.Writer, keyName string) error
//注意在客户端加密的条件下，ufile暂不支持文件分片上传下载操作,因此客户端加密后文件下载请使用此接口
//进行客户端加密下载时，需要用户提供加解密密钥，详情见配置文件相关文档

func main() {
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

	err = req.PutWithEncryptFile(uploadFile, remoteFileKey, "")
	if err != nil {
		log.Println("文件上传失败，失败原因：", err.Error())
		return
	}
	log.Println("文件上传成功。")

	log.Println("正在加密下载文件。。。。")
	file, err := os.OpenFile(saveAsName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Println("创建文件失败，错误信息为：", err.Error())
		return
	}
	defer file.Close()

	err = req.DownloadWithDecryptFile(file, remoteFileKey)
	if err != nil {
		log.Println("下载文件出错，出错信息为：", err.Error())
	}

	ok, err := compareFileMd5()
	if err != nil {
		log.Println("客户端加密：比较文件Md5出错，出错信息为：", err.Error())
	}
	if !ok {
		log.Println("客户端加密：比较文件Md5失败")
	} else {
		log.Println("客户端加密：文件客户端加密上传下载成功")
	}

}

func compareFileMd5() (bool, error) {
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
