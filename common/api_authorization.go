package common

import (
	"crypto/sha1"
	"crypto/hmac"
	"encoding/base64"
)

func Authorization(bucketname string, filename string, publickey string, privatekey string, httpverb string, contentmd5 string, contenttype string, mydata string, canonicalizeducloudheaders string) string {

	var CanonicalizedResource = "/" + bucketname + "/" + filename

	var StringToSign = httpverb + "\n" + contentmd5 + "\n" + contenttype + "\n" + mydata + "\n" + canonicalizeducloudheaders + CanonicalizedResource

    //uflog.INFO("StringToSign", StringToSign)

	var Signature = Base64(HmacSha1(privatekey, StringToSign))
	var Authorization = "UCloud" + " " + publickey + ":" + Signature

	//uflog.INFO("Authorization", Authorization)

	return Authorization

}

func Base64(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}

func HmacSha1(UCloudPrivateKey string, StringToSign string) []byte {
	mac := hmac.New(sha1.New, []byte(UCloudPrivateKey))
	mac.Write([]byte(StringToSign))
	return mac.Sum(nil)
}
