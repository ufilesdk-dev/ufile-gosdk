package utils

import (
	"crypto/aes"
	"crypto/cipher"
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

	//异步分片上传，根据分片大小和序号修改IV
	if startPos != 0 {
		ivLen := len(nonce)
		nonce = editIV(nonce, startPos/uint64(ivLen))
	}

	stream := cipher.NewCTR(c, nonce)

	return &Crypto{key, stream, nonce}, nil
}

//修改IV
func editIV(nonce []byte, startPos uint64) []byte {
	for startPos > 0 {
		startPos--
		for i := len(nonce) - 1; i >= 0; i-- {
			nonce[i]++
			if nonce[i] != 0 {
				break
			}
		}
	}
	return nonce
}

//实现数据加解密
func (c *Crypto) XOR(text1 []byte) []byte {
	text2 := make([]byte, len(text1))
	c.ctr.XORKeyStream(text2, text1)

	return text2
}
