package tool

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Salt(password string, salt string, iteration int) string {
	b := []byte(password)
	s := []byte(salt)
	h := md5.New()
	h.Write(s)
	h.Write(b)
	var res []byte
	res = h.Sum(nil)

	for i := 0; i < iteration-1; i++ {
		h.Reset()
		h.Write(res)
		res = h.Sum(nil)
	}
	return hex.EncodeToString(res)
}
