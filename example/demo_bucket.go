package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	configFile = "config.json"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}
	bucketName := generateRandomBucket()
	req, err := ufsdk.NewBucketRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	ExampleCreateBucket(bucketName, req)
	ExampleDescribeBuckt(bucketName, req)
	ExampleUpdateBucket(bucketName, req)
	ExampleDeleteBucket(bucketName, req)
}

func ExampleCreateBucket(bucketName string, req *ufsdk.UFileRequest) {
	log.Println("正在创建 bucket 名字为：", bucketName, "...")
	err := req.CreateBucket(bucketName, "cn-bj", "private", "")
	if err != nil {
		log.Println("创建 bucket 出错，错误信息为：", err.Error())
	}
	var response ufsdk.BucketListResponse
	ufsdk.MarshalResult(req, &response)
	if response.RetCode != 0 {
		log.Println("创建 bucket 出错，服务器返回错误信息为：", response.Message)
	} else {
		log.Println("创建 bucket 成功")
		log.Println(response.String())
	}
}

func ExampleDescribeBuckt(bucketName string, req *ufsdk.UFileRequest) {
	log.Println("正在获取 bucket 详细信息...")
	err := req.DescribeBucket(bucketName, 0, 10, "")
	if err != nil {
		log.Println("获取 bucket 信息出错，错误信息为：", err.Error())
	}
	var response ufsdk.BucketListResponse
	ufsdk.MarshalResult(req, &response)
	if response.RetCode != 0 {
		log.Println("获取 bucket 信息出错，服务器返回的错误信息为：", response.Message)
	} else {
		log.Println("获取 bucket 详细信息成功。")
		log.Println(response.String())
	}
}

func ExampleUpdateBucket(bucketName string, req *ufsdk.UFileRequest) {
	log.Println("正在更新 bucket 信息")
	err := req.UpdateBucket(bucketName, "public", "")
	if err != nil {
		log.Println("更新 bucket 信息失败，错误信息为：", err.Error())
	}
	var response ufsdk.BucketResponse
	ufsdk.MarshalResult(req, &response)
	if response.RetCode != 0 {
		log.Println("更新 bucket 信息失败，服务器返回的错误信息为：", response.Message)
	} else {
		log.Println("更新 bucket 信息成功。")
		log.Println(response.String())
	}
}

func ExampleDeleteBucket(bucketName string, req *ufsdk.UFileRequest) {
	log.Println("正在删除 bucket ")
	err := req.DeleteBucket(bucketName, "")
	if err != nil {
		log.Println("删除 bucket 失败，错误信息为：", err.Error())
	}
	var response ufsdk.BucketResponse
	ufsdk.MarshalResult(req, &response)
	if response.RetCode != 0 {
		log.Println("删除 bucket 失败，服务器返回的错误信息为：", response.Message)
	} else {
		log.Println("删除 bucket 成功。")
		log.Println(response.String())
	}
}

func generateRandomBucket() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := 8
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + seededRand.Intn(25))
	}
	return strings.ToLower(string(bytes))
}
