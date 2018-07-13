package main

import (
	"ufile-gosdk"
	"fmt"
)

func bucketmgr() {
	ufile := ufile_gosdk.GoUfile{
		Publickey:"**************************************",
		Privatekey:"**************************************",
		ProxyHost:ufile_gosdk.CN_SHANGHAI,
	}

	res1, err := ufile.CreateBucket("ufile-go-test", "private", "")
	fmt.Println("res:", res1, ", err = ", err)
	if err != nil || res1.RetCode != 0 {
		fmt.Println("CreateBucket error")
	}

	res2, err := ufile.DescribeBucket("ufile-go-test", "")
	fmt.Println("res:", res2, ", err = ", err)
	if err != nil || res2.RetCode != 0 {
		fmt.Println("DescribeBucket error")
	}

	res3, err := ufile.UpdateBucket("ufile-go-test", "public", "")
	fmt.Println("res:", res3, ", err = ", err)
	if err != nil || res3.RetCode != 0 {
		fmt.Println("UpdateBucket error")
	}

	res4, err := ufile.PrefixFileList("ufile-go-test", "test", "", 10)
	fmt.Println("res:", res4, ", err = ", err)
	if err != nil || res4.RetCode != 0 {
		fmt.Println("PrefixFileList error")
	}

	res5, err := ufile.DeleteBucket("zqy-test2", "org-1iwg5b")
	fmt.Println("res:", res5, ", err = ", err)
	if err != nil || res5.RetCode != 0 {
		fmt.Println("DeleteBucket error")
	}
}

func main() {
	bucketmgr()
}
