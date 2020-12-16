package main

import (
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	"log"
)

const (
	configFile = "config.json"
	bucketName = "bucket-create-test"
	regionId = "cn-bj"
	bucketType = "private"
	updateType = "public"
	projectId = "org-4bf30w"
	offset = 0
	limit = 10
)

func main() {
	// 加载配置，创建请求
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}
	req, err := ufsdk.NewBucketRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}
	// 创建存储空间
	bucketRet, err := req.CreateBucket(bucketName, regionId, bucketType, projectId)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Fatalf("创建 bucket 出错，错误信息为：%s\n", err.Error())
	}
	log.Println("创建 Bucket 成功，bucket 为", bucketRet)

	// 获取存储空间信息
	//  {
	//   "BucketName": "applcations",
	//   "BucketId": "ufile-vp0csiqf",
	//   "Domain": {
	//    "Src": [
	//     "applcations.cn-sh2.ufileos.com"
	//    ],
	//    "Cdn": [
	//     "applcations.ufile.ucloud.com.cn"
	//    ]
	//   },
	//   "Type": "public",
	//   "CreateTime": 1601117058,
	//   "Biz": "general",
	//   "Region": "cn-sh2"
	//  },
	// 获取指定 BucketName 详细信息
	bucketInfo, err := req.DescribeBucket(bucketName, offset, limit, projectId)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Println("获取 bucket 信息出错，错误信息为：", err.Error())
	} else {
		log.Println("获取 bucket 信息成功，bucketInfo 为：", bucketInfo)
	}

	// 获取指定 projectId 下所有Bucket信息
	bucketList, err := req.DescribeBucket("", offset, limit, projectId)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Println("获取 bucket 信息出错，错误信息为：", err.Error())
	} else {
		log.Println("获取 bucket list 成功，list 为", bucketList)
	}

	//更新存储空间
	bucketRet, err = req.UpdateBucket(bucketName , updateType , projectId)
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Println("更新 bucket 信息失败，错误信息为：", err.Error())
	} else {
		log.Println("Bucket 更新成功。返回信息为：", bucketRet)
	}
	// 删除存储空间
	bucketRet, err = req.DeleteBucket(bucketName, projectId )
	if err != nil {
		log.Println("DumpResponse：", string(req.DumpResponse(true)))
		log.Println("删除 bucket 失败，错误信息为：", err.Error())
	} else {
		log.Println("删除 bucket 成功。返回信息为：", bucketRet)
	}
}
