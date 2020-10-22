package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

//Crypto模块包含了客户端加解密相关内容
//本工具目前仅支持AES-GCM-NoPadding加密方式
type Crypto struct {
	key   []byte
	gcm   cipher.AEAD
	nonce []byte
}

const (
	defaultNonce     = "000000000000"
	defaultNonceSize = 12
)

//NewCrypto构造一个新的Crypto，传入你的加密密钥
//加密密钥位数限定16,24或32，分别对应AES-128，AES-192或AES-256算法
//本工具中，GCM下NonceSize为默认长度12，不提供修改接口
//如若修改NonceSize，请自行实现并务必牢记加解密保持一致
func NewCrypto(key []byte) (cry *Crypto, err error) {
	c, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	g, err := cipher.NewGCMWithNonceSize(c, defaultNonceSize)
	if err != nil {
		return nil, err
	}

	nonce := []byte(defaultNonce)
	return &Crypto{key, g, nonce}, nil
}

//Encrypt实现数据加密
func (c *Crypto) Encrypt(plaintext []byte) ([]byte, error) {
	return c.gcm.Seal(nil, c.nonce, plaintext, nil), nil
}

//Decrypt实现数据解密
func (c *Crypto) Decrypt(ciphertext []byte) ([]byte, error) {
	return c.gcm.Open(nil, c.nonce, ciphertext, nil)
}
