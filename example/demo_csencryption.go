package main

import (
	"bytes"
	"crypto/md5"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"io/ioutil"
	"log"
	"os"
)

const (
	uploadFile    = "./plaintext.txt"
	configFile    = "config.json"
	remoteFileKey = "ciphertext.txt"
	saveAsName    = "./download.txt"
)

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
