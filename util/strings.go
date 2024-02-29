package util

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"time"
	"unsafe"
)

const (
	LETTLE_IDX_BITS       = 6                      // 6 bits to represent a letter index
	LETTLE_IDX_MASK       = 1<<LETTLE_IDX_BITS - 1 // All 1-bits, as many as letterIdxBits
	LETTLE_IDX_MAX        = 63 / LETTLE_IDX_BITS   // # of letter indices fitting in 63 bits
	LETTER_BYTES          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NUM_BYTES             = "0123456789"
	NUM_PLUS_LETTER_BYTES = NUM_BYTES + LETTER_BYTES
)

var randStringBytesMaskImprSrcUnsafeSrc = rand.NewSource(time.Now().UnixNano())

//该函数不保证去除重复，不可直接当作Id生成器来使用
func RandomStrByCharacterSet(n int, characterSet string) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, randStringBytesMaskImprSrcUnsafeSrc.Int63(), LETTLE_IDX_MAX; i >= 0; {
		if remain == 0 {
			cache, remain = randStringBytesMaskImprSrcUnsafeSrc.Int63(), LETTLE_IDX_MAX
		}
		if idx := int(cache & LETTLE_IDX_MASK); idx < len(characterSet) {
			b[i] = characterSet[idx]
			i--
		}
		cache >>= LETTLE_IDX_BITS
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func RandomNumStr(n int) string {
	return RandomStrByCharacterSet(n, NUM_BYTES)
}

func RandomNumLetterStr(n int) string {
	return RandomStrByCharacterSet(n, NUM_PLUS_LETTER_BYTES)
}

func EncodeBase64Str(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64Str(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func ToJsonStr(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}
