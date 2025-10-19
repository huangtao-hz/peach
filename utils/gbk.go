// Package utils  GBK转码包
// Writer : Huang Tao 2020/03/01
// 支持将 GBK 编码的二进制解码为字符串，或将字符串使用 GBK 编码成二进制
// 修订   ： 2025-08-05
package utils

import (
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	// Encoder GBK 编码器
	Encoder *encoding.Encoder = simplifiedchinese.GBK.NewEncoder()
	// Decoder GBK 解码器
	Decoder *encoding.Decoder = simplifiedchinese.GBK.NewDecoder()
)

// NewReader 新建一个 GBK 编码的 Reader
func NewReader(r io.Reader) *transform.Reader {
	return transform.NewReader(r, Decoder)
}

// NewWriter 新建一个 GBK 编码的 Writer
func NewWriter(w io.Writer) *transform.Writer {
	return transform.NewWriter(w, Encoder)
}

// Encode 将 UTF8 字符串转换成 GBK 编码的 bytes
func Encode(s string) ([]byte, error) {
	return Encoder.Bytes([]byte(s))
}

// Decode 将 GBK 编码的 bytes 转换成 UTF8 字符串
func Decode(b []byte) (string, error) {
	return Decoder.String(string(b))
}

// Wlen 计算字符串的长度，汉字算两个字节
func Wlen(s string) int {
	bytes, err := Encode(s)
	CheckFatal(err)
	return len(bytes)
}

// IsGBK 判断是否为GBK编码
func IsGBK(bytes []byte) bool {
	for i := 0; i < len(bytes); i++ {
		if bytes[i] >= 0x81 && bytes[i] <= 0xFE {
			if i+1 >= len(bytes) || bytes[i+1] < 0x40 || (bytes[i+1] > 0x7E && bytes[i+1] < 0x80) || bytes[i+1] > 0xFE {
				return false
			}
			i++
		} else if bytes[i] > 0x7F {
			return false
		}
	}
	return true
}

// IsUTF8 判断是否为 UTF8 编码
func IsUTF8(bytes []byte) bool {
	for i := 0; i < len(bytes); {
		if bytes[i]&0x80 == 0x00 {
			i++
		} else if bytes[i]&0xE0 == 0xC0 {
			if i+1 >= len(bytes) || bytes[i+1]&0xC0 != 0x80 {
				return false
			}
			i += 2
		} else if bytes[i]&0xF0 == 0xE0 {
			if i+2 >= len(bytes) || bytes[i+1]&0xC0 != 0x80 || bytes[i+2]&0xC0 != 0x80 {
				return false
			}
			i += 3
		} else if bytes[i]&0xF8 == 0xF0 {
			if i+3 >= len(bytes) || bytes[i+1]&0xC0 != 0x80 || bytes[i+2]&0xC0 != 0x80 || bytes[i+3]&0xC0 != 0x80 {
				return false
			}
			i += 4
		} else {
			return false
		}
	}
	return true
}
