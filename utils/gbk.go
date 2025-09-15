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

// // 计算字符串的长度，汉字算两个字节
func Wlen(s string) int {
	bytes, err := Encode(s)
	CheckFatal(err)
	return len(bytes)
}
