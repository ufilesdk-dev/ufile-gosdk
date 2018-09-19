# Ucloud 对象存储 SDK

> Modules are interface and implementation.    
> The best modules are where interface is much simpler than implementation.  
> **By: John Ousterhout**

## UFile 对象存储基本概念
在对象存储系统中，我们把文件存放在 bucket 里面，key 是用来索引文件用的。你可以把 bucket 理解成文件系统里么的 folder(文件夹)，key 理解成文件名。由于每个 bucket 需要配置和权限不同，所以每个账户里面会有多个 bucket。在 ufile 里面，bucket 主要分为公有和私有两种，公有 bucket 里面的文件可以对任何人开放，而私有 bucket 需要用户配置访问的签名才能访问。

### 签名
本 SDK 的用户接口是基于 HTTP 的，为了连接的安全性，我们使用 HMAC SHA1 对每个连接进行签名校验。使用 SDK 你可以忽略签名相关的操作，你只要把公私钥写入到配置文件里面（注意不要传到版本控制里面），读取并传给 UFileRequest 里面的创建 instance 的函数即可。
签名相关的算法与详细实现请见 [Auth 模块](auth.go)

## 示例代码
SDK 主要分为两个模块，一个是 bucket 管理，一个是 file 管理。使用对象存储你需要频繁的调用 file 管理相关的接口，bucket 管理用到的地方不会太频繁。一下是用 SDK 上传一个文件的例子：
```go
    import ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}
    req := ufsdk.NewUFileRequest(config, nil)
    err = req.PutFile(filePath, keyName, "")
	if err != nil {
        log.Println("文件上传失败!!，错误信息为：", err.Error())
        //把 HTTP 详细的 response dump 出来
        req.DumpResponse(true)
    }
```
更详细的代码请参考 [example/test_ufile.go](/example/test_ufile.go) 和 [example/test_ubucket.go](example/test_ubucket.go)

## 文档说明
本 SDK 使用 [godoc](https://blog.golang.org/godoc-documenting-go-code) 约定的方法对每个 export 出来的接口进行注释。
你可以直接访问生成好的[在线文档](https://godoc.org/github.com/ufilesdk-dev/ufile-gosdk)。

## 如何排错？
使用 UFileRequest 里面的方法对返回的 error 进行检查。如果不为 nil，调用 Error() 查看错误信息。调用 DumpResponse(true) 查看详细的 HTTP 具体返回值。 