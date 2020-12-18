package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	ufsdk "github.com/kuixiao/ufile-gosdk"
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
	bucketName := generateRandomBucketName()
	req, err := ufsdk.NewBucketRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	log.Println("正在创建 bucket 名字为：", bucketName, "...")
	bucketRet, err := req.CreateBucket(bucketName, "cn-bj", "private", "")
	if err != nil {
		log.Printf("创建 bucket 出错，错误信息为：%s\n", err.Error())
	} else {
		log.Println("创建 Bucket 成功，bucket 为", bucketRet)
	}

	log.Println("正在获取 bucket 详细信息...")
	bucketList, err := req.DescribeBucket(bucketName, 0, 10, "")
	if err != nil {
		log.Println("获取 bucket 信息出错，错误信息为：", err.Error())
	} else {
		log.Println("获取 bucket list 成功，list 为", bucketList)
	}

	log.Println("正在更新 bucket 信息")
	bucketRet, err = req.UpdateBucket(bucketName, "public", "")
	if err != nil {
		log.Println("更新 bucket 信息失败，错误信息为：", err.Error())
	} else {
		log.Println("Bucket 更新成功。")
	}

	log.Println("正在删除 bucket 信息")
	bucketRet, err = req.DeleteBucket(bucketName, "")
	if err != nil {
		log.Println("删除 bucket 失败，错误信息为：", err.Error())
	} else {
		log.Println("删除 bucket 成功")
	}
}

func generateRandomBucketName() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := 8
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + seededRand.Intn(25))
	}
	return strings.ToLower(string(bytes))
}
