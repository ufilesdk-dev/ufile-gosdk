package ufile_gosdk

import (
	"ufile-gosdk/bucketmgr"
	"ufile-gosdk/filemgr"
	"os"

	"ufile-gosdk/uflog"
)

type GoUfile struct {
	Publickey  string
	Privatekey string
	ProxyHost  string
}

type Config struct {
	RetryCount        int
	RetryInterval     int
	BlokSize          int64
	Expires           int64
	ConnectionTimeout int
}

var (
	commonConfig = Config{}
)

func init() {
	commonConfig = Config{
		RetryCount:  3,
		RetryInterval:1,
		BlokSize:4194304,
		Expires:300,
		ConnectionTimeout:5,
	}
}

func (ufile *GoUfile) SetConfig(setConfig Config) {
	if setConfig.RetryCount != 0 {
		commonConfig.RetryCount = setConfig.RetryCount
	}

	if setConfig.RetryInterval != 0 {
		commonConfig.RetryInterval = setConfig.RetryInterval
	}

	if setConfig.BlokSize != 0 {
		commonConfig.BlokSize = setConfig.BlokSize
	}

	if setConfig.Expires != 0 {
		commonConfig.Expires = setConfig.Expires
	}

	if setConfig.ConnectionTimeout != 0 {
		commonConfig.ConnectionTimeout = setConfig.ConnectionTimeout
	}
}

func (ufile *GoUfile) GetConfig() Config {
	return commonConfig
}


/**
 *bucketName	空间名
 *bucketType	空间访问类型
 *projectId	项目id
 *
 */
func (ufile *GoUfile) CreateBucket(bucketName string, bucketType string, projectId string) (bucketmgr.CreateBucketResponse, error) {
	uflog.INFO("***************CreateBucket:创建空间***************")

	response, err := bucketmgr.CreateBucket(bucketName, bucketType, Region[ufile.ProxyHost], projectId, ufile.Publickey, ufile.Privatekey)
	if err != nil {
		uflog.ERROR(err)
	}

	uflog.INFO("***************Complete:创建空间结束***************")

	return response, err

}

/**
 *bucketName	空间名
 *projectId	项目id
 *
 */
func (ufile *GoUfile) DeleteBucket(bucketName string, projectId string) (bucketmgr.DeleteBucketResponse, error) {
	uflog.INFO("***************DeleteBucket:删除空间***************")

	response, err := bucketmgr.DeleteBucket(bucketName, projectId, ufile.Publickey, ufile.Privatekey)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:删除空间结束***************")

	return response, err

}

/**
 *bucketName	空间名
 *projectId	项目id
 *
 */
func (ufile *GoUfile) DescribeBucket(bucketName string, projectId string) (bucketmgr.DescribeBucketResponse, error) {
	uflog.INFO("***************DescribeBucket:获取空间信息***************")

	response, err := bucketmgr.DescribeBucket(bucketName, projectId, ufile.Publickey, ufile.Privatekey)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:获取空间信息结束***************")

	return response, err

}

/**
 *bucketName	空间名
 *bucketType	空间访问类型
 *projectId	项目id
 *
 */
func (ufile *GoUfile) UpdateBucket(bucketName string, bucketType string, projectId string) (bucketmgr.UpdateBucketResponse, error) {

	uflog.INFO("***************UpdateBucket:更改空间属性***************")

	response, err := bucketmgr.UpdateBucket(bucketName, bucketType, projectId, ufile.Publickey, ufile.Privatekey)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:更改空间属性结束***************")

	return response, err

}

/**
 *bucketName	空间名
 *prefix	前缀
 *marker	标志字符串
 *limit		文件列表数目
 */
func (ufile *GoUfile) PrefixFileList(bucketName string, prefix string, marker string, limit int) (bucketmgr.PrefixFileListResponse, error) {
	uflog.INFO("***************PrefixFileList:前缀列表查询***************")

	response, err := bucketmgr.PrefixFileList(bucketName, prefix, marker, limit, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:前缀列表查询结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *path		文件路径
 *filename	文件名
 *contentType	文件类型
 */
func (ufile *GoUfile) PutFile(bucketName string, path string, filename string, contentType string) (filemgr.PutFileResponse, error) {
	uflog.INFO("***************PutFile:上传文件***************")

	response, err := filemgr.PutFile(bucketName, path, filename, contentType, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:上传文件结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *path		文件路径
 *filename	文件名
 *byterange	分片下载的文件范围
 *isprivate	是否是私有空间
 */
func (ufile *GoUfile) GetFile(bucketName string, path string, filename string, byterange string, isprivate bool) (filemgr.GetFileResponse, error) {
	uflog.INFO("***************GetFile:下载文件***************")

	response, err := filemgr.GetFile(bucketName, path, filename, byterange, isprivate, commonConfig.Expires, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost, commonConfig.ConnectionTimeout)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:下载文件结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *path		文件路径
 *filename	文件名
 */
func (ufile *GoUfile) UploadHit(bucketName string, path string, filename string) (filemgr.UploadHitResponse, error) {
	uflog.INFO("***************UploadHit:秒传文件***************")

	response, err := filemgr.UploadHit(bucketName, path, filename, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost, commonConfig.ConnectionTimeout)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:秒传文件结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *filename	文件名
 */
func (ufile *GoUfile) DeleteFile(bucketName string, filename string) (filemgr.DeleteFileResponse, error) {
	uflog.INFO("***************DeleteFile:删除文件***************")

	response, err := filemgr.DeleteFile(bucketName, filename, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:删除文件结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *filename	文件名
 */
func (ufile *GoUfile)InitiateMultipartUpload(bucketName string, filename string) (filemgr.MultipartUploadResponse, error) {
	uflog.INFO("***************InitiateMultipartUpload:初始化分片***************")

	response, err := filemgr.InitiateMultipartUpload(bucketName, filename, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:初始化分片结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *f		文件流
 *filename	文件名
 *uploadid	本次分片上传的上传Id
 *partnumber	本次分片上传的分片号码
 */
func (ufile *GoUfile)UploadPartStream(bucketName string, f *os.File, filename string, uploadid string, partnumber int64) (filemgr.UploadPartResponse, error) {
	uflog.INFO("***************UploadPartStream:上传分片文件流***************")

	response, err := filemgr.UploadPartStream(bucketName, f, filename, uploadid, partnumber, commonConfig.RetryCount, commonConfig.RetryInterval, "", ufile.Publickey, ufile.Privatekey, ufile.ProxyHost, commonConfig.BlokSize, commonConfig.ConnectionTimeout)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:上传分片文件流结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *path		文件路径
 *filename	文件名
 *uploadid	本次分片上传的上传Id
 *partnumber	本次分片上传的分片号码
 */
func (ufile *GoUfile)UploadPartFile(bucketName string, path string, filename string, uploadid string, partnumber int64) (filemgr.UploadPartResponse, error) {
	uflog.INFO("***************UploadPartFile:上传分片文件***************")

	response, err := filemgr.UploadPartFile(bucketName, path, filename, uploadid, partnumber, commonConfig.RetryCount, commonConfig.RetryInterval, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost, commonConfig.BlokSize, commonConfig.ConnectionTimeout)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:上传分片文件结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *f		文件流
 *filename	文件名
 *uploadid	本次分片上传的上传Id
 *partnumber	本次分片上传的分片号码
 *etags		已上传分片etags
 */
func (ufile *GoUfile)ResumeUploadPartStream(bucketName string, f *os.File, filename string, uploadid string, partnumber int64, etags []string) (filemgr.UploadPartResponse, error) {
	uflog.INFO("***************ResumeUploadPartStream:续传分片文件流***************")

	response, err := filemgr.ResumeUploadPartStream(bucketName, f, filename, uploadid, partnumber, etags, commonConfig.RetryCount, commonConfig.RetryInterval, "", ufile.Publickey, ufile.Privatekey, ufile.ProxyHost, commonConfig.BlokSize, commonConfig.ConnectionTimeout)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:续传分片文件流结束***************")

	return response, err
}

/**
 *bucketName	空间名
 *path		文件路径
 *filename	文件名
 *uploadid	本次分片上传的上传Id
 *partnumber	本次分片上传的分片号码
 *etags		已上传分片etags
 */
func (ufile *GoUfile)ResumeUploadPartFile(bucketName string, path string, filename string, uploadid string, partnumber int64, etags []string, ) (filemgr.UploadPartResponse, error) {
	uflog.INFO("***************ResumeUploadPartFile:续传分片文件***************")

	response, err := filemgr.ResumeUploadPartFile(bucketName, path, filename, uploadid, partnumber, etags, commonConfig.RetryCount, commonConfig.RetryInterval, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost, commonConfig.BlokSize, commonConfig.ConnectionTimeout)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:续传分片文件结束***************")

	return response, err
}

//func (ufile *GoUfile)GetMultiUploadPart(bucketName string, uploadid string) (filemgr.GetMultiUploadPartResponse, error) {
//	response, err := filemgr.GetMultiUploadPart(bucketName, uploadid, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost)
//
//	return response, err
//}

/**
 *bucketName	空间名
 *filename	文件名
 *uploadid	本次分片上传的上传Id
 *etags		已上传分片etags
 */
func (ufile *GoUfile)FinishMultipartUpload(bucketName string, filename string, etags []string, uploadid string) (filemgr.FinishMultipartUploadResponse, error) {
	uflog.INFO("***************FinishMultipartUpload:完成分片***************")

	response, err := filemgr.FinishMultipartUpload(bucketName, filename, etags, uploadid, ufile.Publickey, ufile.Privatekey, ufile.ProxyHost, commonConfig.ConnectionTimeout)
	if err != nil {
		uflog.ERROR(err)
	}
	uflog.INFO("***************Complete:完成分片结束***************")

	return response, err
}
