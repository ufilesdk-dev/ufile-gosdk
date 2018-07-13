package common

import (
	"net/url"
	"encoding/hex"
	"crypto/sha1"
	"strings"
	"ufile-gosdk/uflog"
)

func Signature(privatekey string, data url.Values) url.Values {

	datastring, _ := url.PathUnescape(data.Encode())

	var signaturestring string
	datamap := strings.Split(datastring, "&")
	for _, v := range datamap {
		datavalue := strings.SplitN(v, "=", 2)
		for _, dv := range datavalue {
			signaturestring = signaturestring + dv
		}
	}

	signaturestring = signaturestring + privatekey

	h := sha1.New()

	h.Write([]byte(signaturestring))

	bs := h.Sum(nil)

	sbs := hex.EncodeToString(bs)

	data.Add("Signature", sbs)

	uflog.INFO("Signature", sbs)

	return data
}