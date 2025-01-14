package ufsdk

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"mime"
	"path/filepath"
)

const (
	blkSIZE = 2 << 21
)

//Config 配置文件序列化所需的全部字段
type Config struct {
	PublicKey       string `json:"public_key"`
	PrivateKey      string `json:"private_key"`
	BucketHost      string `json:"bucket_host"`
	BucketName      string `json:"bucket_name"`
	FileHost        string `json:"file_host"`
	VerifyUploadMD5 bool   `json:"verfiy_upload_md5"`
	Endpoint        string `json:"endpoint"`
}

var extra2Mimetype = map[string]string{
    ".js": "application/javascript",
    ".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    ".xltx": "application/vnd.openxmlformats-officedocument.spreadsheetml.template",
    ".potx": "application/vnd.openxmlformats-officedocument.presentationml.template",
    ".ppsx": "application/vnd.openxmlformats-officedocument.presentationml.slideshow",
    ".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    ".sldx": "application/vnd.openxmlformats-officedocument.presentationml.slide",
    ".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    ".dotx": "application/vnd.openxmlformats-officedocument.wordprocessingml.template",
    ".xlam": "application/vnd.ms-excel.addin.macroEnabled.12",
    ".xlsb": "application/vnd.ms-excel.sheet.binary.macroEnabled.12",
    ".apk": "application/vnd.android.package-archive",
    ".ipa": "application/vnd.ios.package-archive",
}

//LoadConfig 从配置文件加载一个配置。
func LoadConfig(jsonPath string) (*Config, error) {
	file, err := openFile(jsonPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	err = json.Unmarshal(configBytes, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//VerifyHTTPCode 检查 HTTP 的返回值是否为 2XX，如果不是就返回 false。
func VerifyHTTPCode(code int) bool {
	if code < http.StatusOK || code > http.StatusIMUsed {
		return false
	}
	return true
}

//GetFileMimeType 获取文件的 mime type 值，接收文件路径作为参数。如果检测不到，则返回空。
func GetFileMimeType(path string) string {
	f, err := openFile(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	return getMimeTypeFromData(f)
}

func getMimeTypeFromData(f *os.File) string {
	buffer := make([]byte, 512)
	_, err := f.Seek(0, 0)
	if err != nil {
		return "plain/text"
	}
	_, err = f.Read(buffer)
	defer func() { f.Seek(0, 0) }() //revert file's seek
	if err != nil {
		return "plain/text"
	}

	return http.DetectContentType(buffer)
}

func getMimeTypeFromFilename(path string) string {
	ext := filepath.Ext(path)
	if mimetype, ok := extra2Mimetype[ext]; ok {
		return mimetype;
	}
	return mime.TypeByExtension(ext)
}

func openFile(path string) (*os.File, error) {
	return os.Open(path)
}

//getFileSize get opened file size
func getFileSize(f *os.File) int64 {
	fi, err := f.Stat()
	if err != nil {
		panic(err.Error())
	}
	return fi.Size()
}

//GetFileEtag 获取文件的 etag 值
func GetFileEtag(path string) string {
	f, err := openFile(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	return calculateEtag(f)
}

//Calculatek 计算文件的 etag 值。
func calculateEtag(f *os.File) string {
	fsize := getFileSize(f)
	blkcnt := uint32(fsize / blkSIZE)
	if fsize%blkSIZE != 0 {
		blkcnt++
	}

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, blkcnt)

	h := sha1.New()
	buf := make([]byte, 0, 24)
	buf = append(buf, bs...)
	if fsize <= blkSIZE {
		io.Copy(h, f)
	} else {
		var i uint32
		for i = 0; i < blkcnt; i++ {
			shaBlk := sha1.New()
			io.Copy(shaBlk, io.LimitReader(f, blkSIZE))
			io.Copy(h, bytes.NewReader(shaBlk.Sum(nil)))
		}
	}
	buf = h.Sum(buf)
	etag := base64.URLEncoding.EncodeToString(buf)
	return etag
}

func makeBoundry() string {
	h := md5.New()
	t := time.Now()
	io.WriteString(h, t.String())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func makeFormBody(authorization, boundry, keyName, mimeType string, verifyMD5 bool, file *os.File) (*bytes.Buffer, error) {
	boundry = "--" + boundry + "\r\n"
	boundryBytes := []byte(boundry)
	body := new(bytes.Buffer)

	body.Write(boundryBytes)
	body.Write(makeFormField("Authorization", authorization))
	body.Write(boundryBytes)
	body.Write(makeFormField("Content-Type", mimeType))
	body.Write(boundryBytes)
	body.Write(makeFormField("FileName", keyName))
	body.Write(boundryBytes)

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return body, err
	}

	if verifyMD5 {
		md5Str := fmt.Sprintf("%x", md5.Sum(b))
		body.Write(makeFormField("Content-MD5", md5Str))
		body.Write(boundryBytes)
	}

	addtionalStr := fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n", keyName)
	addtionalStr += fmt.Sprintf("Content-Type: %s\r\n\r\n", mimeType)
	body.Write([]byte(addtionalStr))
	body.Write(b)
	body.Write([]byte("\r\n"))
	body.Write(boundryBytes)

	return body, nil
}

func makeFormField(key, value string) []byte {
	keyStr := fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"\r\n\r\n", key)
	valueStr := fmt.Sprintf("%s\r\n", value)
	return []byte(keyStr + valueStr)
}

func structPrettyStr(data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", " ")
	if err == nil {
		return fmt.Sprintf("%s\n", bytes)
	}
	return ""
}

// FilePart is the file part definition
type FilePart struct {
	Number int   // Part number
	Offset int64 // Part offset
	Size   int64 // Part size
}

// SplitFileByPartSize splits big file into parts by the size of parts.
// Splits the file by the part size. Returns the FilePart when error is nil.
func SplitFileByPartSize(fileSize, partSize int64) ([]FilePart, error) {
	if fileSize <= 0 || partSize <= 0 {
		return nil, errors.New("fileSize or partSize invalid")
	}

	var partN = fileSize / partSize
	var parts []FilePart
	var part = FilePart{}
	for i := int64(0); i < partN; i++ {
		part.Number = int(i)
		part.Offset = i * partSize
		part.Size = partSize
		parts = append(parts, part)
	}

	if fileSize%partSize > 0 {
		part.Number = len(parts)
		part.Offset = int64(len(parts)) * partSize
		part.Size = fileSize % partSize
		parts = append(parts, part)
	}

	return parts, nil
}
