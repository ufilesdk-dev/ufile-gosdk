package main

import (
	"ufile-gosdk"
	"fmt"
	"os"
)

func filemgr() {
	ufile := ufile_gosdk.GoUfile{
		Publickey:"***************************************",
		Privatekey:"***************************************",
		ProxyHost:ufile_gosdk.CN_SHANGHAI,
	}

	// 初始化配置
	ufileconfig := ufile_gosdk.Config{
		RetryCount        :10,
		RetryInterval     :6,
		BlokSize          :4*1024*1024,
		Expires           :0,
		ConnectionTimeout :500,
	}
	ufile.SetConfig(ufileconfig)
	fmt.Println(ufile.GetConfig())

	// 上传文件到公有空间
	res1, err := ufile.PutFile("ufile-go-public", "/Users/Desktop/UCloud.pdf", "test1/test2/ucloud_public.pdf", "")
	fmt.Println("res1:", res1, ", err = ", err)
	if err != nil || res1.RetCode != 0 {
		fmt.Println("PutFile error")
	}

	// 上传文件到私有空间
	res2, err := ufile.PutFile("ufile-go-private", "/Users/Desktop/UCloud.pdf", "test3/test2/ucloud_private.pdf", "")
	fmt.Println("res2:", res2, ", err = ", err)
	if err != nil || res2.RetCode != 0 {
		fmt.Println("PutFile error")
	}

	// 从公有空间下载文件
	res3, err := ufile.GetFile("ufile-go-public", "/Users/Desktop/", "test1/test2/ucloud_public.pdf", "", false)
	fmt.Println("res3:", res3, ", err = ", err)
	if err != nil || res3.RetCode != 0 {
		fmt.Println("GetFile error")
	}

	// 从私有空间下载文件
	res4, err := ufile.GetFile("ufile-go-private", "/Users/Desktop/", "test3/test2/ucloud_private.pdf", "", true)
	fmt.Println("res4:", res4, ", err = ", err)
	if err != nil || res4.RetCode != 0 {
		fmt.Println("GetFile error")
	}

	// 秒传到公有空间，hit
	res5, err := ufile.UploadHit("ufile-go-public", "/Users/Desktop/UCloud.pdf", "test2/ucloud_public.pdf")
	fmt.Println("res5:", res5, ", err = ", err)
	if err != nil || res5.RetCode != 0 {
		fmt.Println("UploadHit error")
	}

	// 秒传到公有空间，not hit
	res6, err := ufile.UploadHit("ufile-go-public", "/Users/Desktop/ddd.png", "test2/ddd.png")
	fmt.Println("res6:", res6, ", err = ", err)
	if err != nil || res6.RetCode != 0 {
		fmt.Println("UploadHit error")
	}

	// 秒传到私有空间，hit
	res7, err := ufile.UploadHit("ufile-go-private", "/Users/Desktop/UCloud.pdf", "test2/ucloud_private.pdf")
	fmt.Println("res7:", res7, ", err = ", err)
	if err != nil || res7.RetCode != 0 {
		fmt.Println("UploadHit error")
	}

	// 秒传到私有空间，not hit
	res8, err := ufile.UploadHit("ufile-go-private", "/Users/Desktop/ddd.png", "test2/ddd.png")
	fmt.Println("res8:", res8, ", err = ", err)
	if err != nil || res8.RetCode != 0 {
		fmt.Println("UploadHit error")
	}

	// 从公有空间删除文件
	res9, err := ufile.DeleteFile("ufile-go-public", "test2/ucloud_public.pdf")
	fmt.Println("res9:", res9, ", err = ", err)
	if err != nil || res9.RetCode != 0 {
		fmt.Println("DeleteFile error")
	}

	// 从私有空间删除文件
	res10, err := ufile.DeleteFile("ufile-go-private", "test2/ucloud_private.pdf")
	fmt.Println("res10:", res10, ", err = ", err)
	if err != nil || res10.RetCode != 0 {
		fmt.Println("DeleteFile error")
	}

	/*
		Stream分片上传
	 */
	// 初始化stream上传分片
	res11, err := ufile.InitiateMultipartUpload("ufile-go-private", "manager.tar.gz")
	fmt.Println("res11:", res11, ", err = ", err)
	if err != nil || res11.RetCode != 0 {
		fmt.Println("InitiateMultipartUpload error")
	}

	f, err := os.Open("/Users/Downloads/manager.tar.gz")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	defer f.Close()

	// 开始分片上传
	fmt.Println("开始分片上传stream:")

	resp, err := ufile.UploadPartStream("ufile-go-private", f, "manager.tar.gz", res11.UploadId, 0)
	fmt.Println("resp:", resp, ", err = ", err)
	if err != nil || resp.RetCode != 0 {
		fmt.Println("UploadPartStream error")
	}
	if err != nil && resp.RetCode != 0 {
		// 失败重传
		fmt.Println("失败重传stream：")
		resp, err := ufile.ResumeUploadPartStream("ufile-go-private", f, "manager.tar.gz", res11.UploadId, resp.PartNumber, resp.ETag)
		fmt.Println("resp:", resp, ", err = ", err)
		if resp.RetCode != 0 {
			fmt.Println("ResumeUploadPartStream error")
		}
	}

	// 结束分片上传
	fmt.Println("结束分片上传stream:")
	res12, err := ufile.FinishMultipartUpload("ufile-go-private", "manager.tar.gz", resp.ETag, res11.UploadId)
	fmt.Println("res12:", res12, ", err = ", err)
	if err != nil || res12.RetCode != 0 {
		fmt.Println("FinishMultipartUpload error")
	}

	/*
		文件分片上传
	 */
	// 初始化文件上传分片
	fmt.Println("初始化文件上传分片:")
	res13, err := ufile.InitiateMultipartUpload("ufile-go-public", "manager1.tar.gz")
	fmt.Println("res13:", res13, ", err = ", err)
	if err != nil || res13.RetCode != 0 {
		fmt.Println("InitiateMultipartUpload error")
	}

	// 开始分片上传
	fmt.Println("开始分片上传文件:")

	resp1, err := ufile.UploadPartFile("ufile-go-public", "/Users/Downloads/manager.tar.gz", "manager1.tar.gz", res13.UploadId, 0)
	fmt.Println("resp:", resp1, ", err = ", err)
	if err != nil || resp1.RetCode != 0 {
		fmt.Println("UploadPartFile error")
	}
	if err != nil && resp1.RetCode != 0 {
		// 失败重传
		fmt.Println("失败重传文件：")
		resp1, err := ufile.ResumeUploadPartFile("ufile-go-public", "/Users/Downloads/manager.tar.gz", "manager1.tar.gz", res13.UploadId, resp1.PartNumber, resp1.ETag)
		fmt.Println("resp:", resp1, ", err = ", err)
		if resp1.RetCode != 0 {
			fmt.Println("ResumeUploadPartFile error")
		}
	}

	// 结束分片上传
	fmt.Println("结束分片上传文件:")
	res14, err := ufile.FinishMultipartUpload("ufile-go-public", "manager1.tar.gz", resp1.ETag, res13.UploadId)
	fmt.Println("res14:", res14, ", err = ", err)
	if err != nil || res14.RetCode != 0 {
		fmt.Println("FinishMultipartUpload error")
	}
}

func main() {
	filemgr()
}
