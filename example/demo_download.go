package main

import (
	"flag"
	"log"
	"os"
	"fmt"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

var (
	key = flag.String("k", "", "download file key.")
	//	path   = flag.String("p", "", "upload file path.")
	config = flag.String("c", "", "config file")
)

func main() {
	flag.Parse()
	if *config == "" || *key == "" {
		flag.PrintDefaults()
		return
	}
	ufConfig, err := ufsdk.LoadConfig(*config)
	if err != nil {
		panic(err.Error())
	}

	req, err := ufsdk.NewFileRequest(ufConfig, nil)
	if err != nil {
		panic(err.Error())
	}
	log.Println("正在下载文件。。。。")

	defer func(){
		if(false){
			fmt.Println("last status: ", req.LastResponseStatus);
			fmt.Println("last header: ", req.LastResponseHeader);
			fmt.Println("last body: ", req.LastResponseBody);
			fmt.Println(string(req.DumpResponse(true)));
		}
	}()

	file, err := os.OpenFile(*key, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err.Error())
	}

	err = req.DownloadFile(file, *key)
	if err != nil {
		log.Println("fail:", err.Error())
		return
	}


	log.Println("success")
}
