package main

import (
    ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
    "log"
	"bytes"
)

const (
    ConfigFile = "config.json"
    keyName = "testappend.txt"
    MimeType = ""
)

func main() {
    log.SetFlags(log.Lshortfile)
    config, err := ufsdk.LoadConfig(ConfigFile)
    if err != nil {
        panic(err.Error())
    }
    req, err := ufsdk.NewFileRequest(config, nil)
    if err != nil {
        panic(err.Error())
    }

	offset := int32(0)
	for i := 0; i < 5; i++ {
		content := []byte("hello ")
		err = req.AppendObject(keyName, offset, bytes.NewBuffer(content))
		if err != nil {
			log.Fatalf("%s\n", err.Error())
			return
		}
		offset += int32(len(content))
	}
    log.Println("文件上传成功!!")
}
