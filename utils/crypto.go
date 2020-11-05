package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

//Crypto模块包含了客户端加解密相关内容
//本工具目前仅支持AES-GCM-NoPadding加密方式
type Crypto struct {
	Key   []byte
	ctr   cipher.Stream
	nonce []byte
}

const (
	defaultNonce = "12345678abcdefgh"
)

//NewCrypto构造一个新的Crypto，传入你的加密密钥
//加密密钥位数限定16,24或32，分别对应AES-128，AES-192或AES-256算法
//本工具中，GCM下NonceSize为默认长度12，不提供修改接口
//如若修改NonceSize，请自行实现并务必牢记加解密保持一致
func NewCrypto(key []byte) (cry *Crypto, err error) {
	//创建cipher.Block接口
	c, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	//创建分组模式
	nonce := []byte(defaultNonce)
	stream := cipher.NewCTR(c, nonce)

	return &Crypto{key, stream, nonce}, nil
}

//实现数据加解密
func (c *Crypto) XOR(text1 []byte) []byte {
	// text2 := make([]byte, len(text1))

	c.ctr.XORKeyStream(text1, text1)
	// fmt.Println(len(text1), len(text2))

	return text1
}
