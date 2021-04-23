
# UCloud 对象存储 SDK <a href="https://godoc.org/github.com/ufilesdk-dev/ufile-gosdk"><img src="https://godoc.org/github.com/ufilesdk-dev/ufile-gosdk?status.svg" alt="GoDoc"></a>
> Modules are interface and implementation.  
> The best modules are where interface is much simpler than implementation.  
> **By: John Ousterhout**

## UFile 对象存储基本概念
在对象存储系统中，存储空间（Bucket）是文件（File）的组织管理单位，文件（File）是存储空间的逻辑存储单元。对于每个账号，该账号里存放的每个文件都有唯一的一对存储空间（Bucket）与键（Key）作为标识。我们可以把 Bucket 理解成一类文件的集合，Key 理解成文件名。由于每个 Bucket 需要配置和权限不同，每个账户里面会有多个 Bucket。在 UFile 里面，Bucket 主要分为公有和私有两种，公有 Bucket 里面的文件可以对任何人开放，私有 Bucket 需要配置对应访问签名才能访问。

### 签名
本 SDK 接口是基于 HTTP 的，为了连接的安全性，UFile 使用 HMAC SHA1 对每个连接进行签名校验。使用本 SDK 可以忽略签名相关的算法过程，只要把公私钥写入到配置文件里面（注意不要传到版本控制里面），读取并传给 UFileRequest 里面的 New 方法即可。  
签名相关的算法与详细实现请见 [Auth 模块](auth.go)

## 安装
`go get github.com/ufilesdk-dev/ufile-gosdk`

### 执行测试
`cd example; go run demo_file.go`

## 功能列表
### 文件操作相关功能
[Put 上传](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.PutFile)
[Post 上传](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.PostFile)  
分片上传 [同步分片上传](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.MPut)，[异步分片上传](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.AsyncMPut)  
手动分片上传，[步骤一](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.InitiateMultipartUpload)，[步骤二](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.UploadPart)，[步骤三](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.FinishMultipartUpload)。[取消分片上传](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.AbortMultipartUpload)  
[文件秒传](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.UploadHit)  
[获取文件列表](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.PrefixFileList)  
[获取目录文件列表](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.ListObjects) 
[获取私有空间下载地址](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.GetPrivateURL)，[获取公有空间下载地址](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.GetPublicURL)。  
[删除文件](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.DeleteFile)  
[查看文件信息](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.HeadFile)  
[下载文件](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.DownloadFile)  
[比对本地与远程文件的 Etag](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.CompareFileEtag)  
[拷贝文件](https://godoc.org/github.com/ufilesdk-dev/ufile-gosdk#UFileRequest.Copy)  
[重命名文件](https://godoc.org/github.com/ufilesdk-dev/ufile-gosdk#UFileRequest.Rename)  
[Put 带回调上传](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.PutFileWithPolicy)
[同步分片上传-带回调](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.MPutWithPolicy)，
[异步分片上传-带回调](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.AsyncMPutWithPolicy)  
[文件客户端加密上传](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.PutWithEncryptFile), [文件客户端解密下载](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.DownloadWithDecryptFile)

### Bucket 操作相关功能
[创建 bucket](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.CreateBucket)  
[删除 bucket](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.DeleteBucket)  
[获取 bucket 列表](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.DescribeBucket)  
[修改 bucket](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#UFileRequest.UpdateBucket)  

### 签名构造
[构造文件管理签名](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#Auth.Authorization)  
[构造私有空间下载签名](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#Auth.AuthorizationPrivateURL)  
[构造 bucket 管理签名](https://godoc.org/pkg/github.com/ufilesdk-dev/ufile-gosdk/#Auth.AuthorizationBucketMgr)  

## 示例代码
SDK 主要分为两个模块，一个是 bucket 管理，一个是 file 管理。使用对象存储你需要频繁的调用 file 管理相关的接口，bucket 管理用到的地方不会太频繁。以下是用 SDK 上传一个文件的例子：
```go
import ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
config, err := ufsdk.LoadConfig(configFile)
if err != nil {
    panic(err.Error())
}
err = req.IOPut(f, "KeyName", "")
f.Close()
if err != nil {
    log.Fatalf("%s\n", req.DumpResponse(true))
}

// 流式上传大文件
f1, err := os.Open("FilePath1")
if err != nil {
    panic(err.Error())
}
err = req.IOMutipartAsyncUpload(f1, "KeyName", "")
f1.Close()
if err != nil {
    log.Fatalf("%s\n", req.DumpResponse(true))
}
```

[回到目录](#table-of-contents)


<a name="分片上传"></a>
### 分片上传

- demo程序

```go
// 加载配置，创建请求
config, err := ufsdk.LoadConfig("config.json")
if err != nil {
	panic(err.Error())
}
req, err := ufsdk.NewFileRequest(config, nil)
if err != nil {
	panic(err.Error())
}

err = req.MPut("FilePath", "KeyName", "MimeType")
if err != nil {
	log.Println("DumpResponse：", string(req.DumpResponse(true)))
}
```

[回到目录](#table-of-contents)

<a name="上传回调"></a>
### 上传回调

- demo程序

```go
// 加载配置，创建请求
config, err := ufsdk.LoadConfig("config.json")
if err != nil {
	panic(err.Error())
}
req, err := ufsdk.NewFileRequest(config, nil)
if err != nil {
	panic(err.Error())
}

// 同步分片上传回调
err = req.MPutWithPolicy("FilePath", "KeyName", "MimeType", "Policy")
if err != nil {
	log.Println("DumpResponse：", string(req.DumpResponse(true)))
}

// 异步分片上传回调
err = req.AsyncMPutWithPolicy("FilePath", "KeyName", "MimeType", "Policy")
if err != nil {
	log.Println("DumpResponse：", string(req.DumpResponse(true)))
}

// 异步分片并发上传回调
jobs := 20 // 并发数为 20
err = req.AsyncUploadWithPolicy("FilePath", "KeyName", "MimeType", jobs, "Policy")
if err != nil {
	log.Println("DumpResponse：", string(req.DumpResponse(true)))
}
```

[回到目录](#table-of-contents)

<a name="文件下载"></a>
### 文件下载

- demo程序

```go
// 加载配置，创建请求
config, err := ufsdk.LoadConfig("config.json")
if err != nil {
	panic(err.Error())
}
req, err := ufsdk.NewFileRequest(config, nil)
if err != nil {
	panic(err.Error())
}
// 普通下载
err = req.Download("DownLoadURL")
if err != nil {
	log.Println("DumpResponse：", string(req.DumpResponse(true)))
}
// 流式下载
err = req.Download("buffer", "KeyName")
if err != nil {
	log.Println("DumpResponse：", string(req.DumpResponse(true)))
}
```

## 客户端加密
用户需先在配置文件中提供密钥，密钥位数限定16,24或32，分别对应AES-128，AES-192或AES-256算法，配置文件参考[example/config.json.example](/example/config.json.example)

### 示例代码
```go
import ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
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
```
更详细的代码请参考 [example/demo_csencryption.go](/example/demo_csencryption.go) 

## 文档说明
本 SDK 使用 [godoc](https://blog.golang.org/godoc-documenting-go-code) 约定的方法对每个 export 出来的接口进行注释。
你可以直接访问生成好的[在线文档](https://godoc.org/github.com/ufilesdk-dev/ufile-gosdk)。  

# 联系我们

> - UCloud US3 [官方网站](https://www.ucloud.cn/site/product/ufile.html)
> - UCloud US3 [官方文档中心](https://docs.ucloud.cn/ufile/README)
> - UCloud US3 [开发者文档](https://ucloud-us3.github.io/go-sdk/概述.html)
> - UCloud 官方技术支持：[提交工单](https://accountv2.ucloud.cn/work_ticket/create)
