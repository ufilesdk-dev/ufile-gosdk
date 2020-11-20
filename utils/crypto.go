package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
)

//Crypto模块包含了客户端加解密相关内容
//本工具目前仅支持AES-CTR加密方式
type Crypto struct {
	key   []byte
	ctr   cipher.Stream
	nonce []byte
}

const (
	defaultNonce = "rcqg4nt58wu401m5"
)

//NewCrypto构造一个新的Crypto
//key 加密密钥，位数限定16,24或32，分别对应AES-128，AES-192或AES-256算法
//startPos Crypto初始化时对应的文件位置，一般为0，仅在异步分片上传加密文件时会传入对应的值
//本工具中，CTR的Nonce为默认值，不提供修改接口
//如若修改Nonce，请自行实现并务必牢记加解密保持一致
func NewCrypto(key []byte, startPos uint64) (cry *Crypto, err error) {
	//创建cipher.Block接口
	c, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	//创建分组模式
	nonce := []byte(defaultNonce)

	if startPos != 0 {
		ivLen := len(nonce)
		iv := binary.BigEndian.Uint64(nonce[ivLen-8:]) + startPos/uint64(ivLen) //字节切片，最后8位上限是256^8 = 2^64 ,+1代表着一个4M的切片，所以这里最大加密文件上限为2^64 *4M = 2^46 TB，远超上传上限
		binary.BigEndian.PutUint64(nonce[ivLen-8:], iv)
	}

	stream := cipher.NewCTR(c, nonce)

	return &Crypto{key, stream, nonce}, nil
}

//实现数据加解密
func (c *Crypto) XOR(text1 []byte) []byte {
	text2 := make([]byte, len(text1))
	c.ctr.XORKeyStream(text2, text1)

	return text2
}
