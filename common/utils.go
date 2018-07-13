package common

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

)

var BLKSIZE int64 = 2 << 21
var consoleW = 0

func TransWinSep(key string) string {

	return strings.Replace(key, "\\", "/", -1)
}

func MimeType(filename string) (mimeType string, err error) {

	if ext := filepath.Ext(filename); ext != "" {
		if m, ok := CompleteMimeTypes[ext]; ok {
			return m, nil
		}
	}
	if runtime.GOOS == "windows" {
		mimeType = mime.TypeByExtension(path.Ext(filename))
		return
	}

	slice, err := exec.Command("/usr/bin/file", "--mime-type", filename).Output()
	if err != nil {
		// 来自于不同平台的文件，file命令会报错
		mimeType = mime.TypeByExtension(path.Ext(filename))
		return mimeType, nil
	}
	s := strings.Split(string(slice), " ")
	if len(s) > 1 {
		mimeType = strings.TrimRight(s[len(s)-1], "\n")
		return
	}
	err = errors.New("invalid command format")
	return
}

func CanonicalizedUcloudHeaders(ucloudHeader map[string]string) string {

	keys := make([]string, 0)
	for k := range ucloudHeader {
		k = strings.ToLower(k)
		ucloudHeader[k] = strings.TrimSpace(ucloudHeader[k])
		keys = append(keys, k)
	}
	sort.Strings(keys)

	s := ""
	for _, k := range keys {
		v := ucloudHeader[k]
		s += k + ":" + v + "\n"
	}
	return s
}

func Resource(bucket string, key string) string {

	return "/" + bucket + "/" + key
}

type CountedReader struct {
	r     io.Reader
	readn int64
	total int64
	show  bool
	err   error
}


func NewCountedReader(r io.Reader, size int64, show bool) *CountedReader {
	return &CountedReader{r: r, readn: 0, total: size, show: show}
}

func ETag(file string) (etag string, err error) {

	f, err := os.Open(file)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	var blkcnt uint32 = uint32(fi.Size() / BLKSIZE)
	if fi.Size()%BLKSIZE != 0 {
		blkcnt += 1
	}

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, blkcnt)

	h := sha1.New()
	buf := make([]byte, 0, 24)
	buf = append(buf, bs...)
	if fi.Size() <= BLKSIZE {
		io.Copy(h, f)
	} else {
		var i uint32
		for i = 0; i < blkcnt; i++ {
			shaBlk := sha1.New()
			io.Copy(shaBlk, io.LimitReader(f, BLKSIZE))
			io.Copy(h, bytes.NewReader(shaBlk.Sum(nil)))
		}
	}
	buf = h.Sum(buf)
	etag = base64.URLEncoding.EncodeToString(buf)
	return etag, nil
}

func ETagByBuffer(buffer *bytes.Buffer) (etag string, err error) {

	var size int64 = int64(buffer.Len())
	var blkcnt uint32 = uint32(size / BLKSIZE)
	if size%BLKSIZE != 0 {
		blkcnt += 1
	}

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, blkcnt)

	h := sha1.New()
	buf := make([]byte, 0, 24)
	buf = append(buf, bs...)
	if size <= BLKSIZE {
		io.Copy(h, buffer)
	} else {
		var i uint32
		for i = 0; i < blkcnt; i++ {
			shaBlk := sha1.New()
			io.Copy(shaBlk, io.LimitReader(buffer, BLKSIZE))
			io.Copy(h, bytes.NewReader(shaBlk.Sum(nil)))
		}
	}
	buf = h.Sum(buf)
	etag = base64.URLEncoding.EncodeToString(buf)
	return etag, nil
}

type ProxyRet struct {
	Retcode int    `json:"RetCode"`
	Errmsg  string `json:"ErrMsg"`
}

//====================================================================================================

type TimeoutConn struct {
	Conn         net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewTimeoutConn(conn net.Conn, readTimeout, writeTimeout time.Duration) *TimeoutConn {

	return &TimeoutConn{
		Conn:         conn,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func (c *TimeoutConn) Read(b []byte) (n int, err error) {
	if c.ReadTimeout > 0 {
		c.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	}
	return c.Conn.Read(b)
}

func (c *TimeoutConn) Write(b []byte) (n int, err error) {
	if c.WriteTimeout > 0 {
		c.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	}
	return c.Conn.Write(b)
}

func (c *TimeoutConn) Close() error {
	return c.Conn.Close()
}

func (c *TimeoutConn) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *TimeoutConn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *TimeoutConn) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

func (c *TimeoutConn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *TimeoutConn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

type TimeoutClient struct {
	Conn net.Conn
	*http.Client
}

func NewTimeoutClient(dialTimeout, readTimeout, writeTimeout time.Duration, ssl bool) (*TimeoutClient, error) {

	tc := &TimeoutClient{}
	transport := &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, dialTimeout)
			if err != nil {
				return nil, err
			}
			tc.Conn = conn
			return NewTimeoutConn(conn, readTimeout, writeTimeout), nil
		},
		ResponseHeaderTimeout: readTimeout,
		DisableKeepAlives:     true,
	}

	if ssl {
		sslcfg := &tls.Config{
			InsecureSkipVerify: true,
		}
		transport.TLSClientConfig = sslcfg
	}

	client := &http.Client{
		Transport: transport,
	}
	tc.Client = client
	return tc, nil
}

func SchemePrefix(domain string) (string, bool) {

	if strings.Index(domain, "https://") != -1 {
		return "https://", true
	}
	return "http://", false
}