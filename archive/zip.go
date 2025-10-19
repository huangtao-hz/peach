package archive

import (
	"archive/zip"
	"peach/utils"
)

// ConvName 解决中文文件名乱码的问题
func convName(name string) string {
	b := []byte(name)
	if d, err := utils.Decode(b); err == nil {
		return d
	}
	return name
}

// isUTF8 判断是否为 UTF8 编码
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

// ExtractZip 读取 .zip 压缩包中的文件
func ExtractZip(file string, fns ...ReadFileFunc) (err error) {
	if r, err := zip.OpenReader(file); err == nil {
		defer r.Close()
		for _, zfile := range r.File {
			b := []byte(zfile.Name)
			if zfile.NonUTF8 && IsGBK(b) {
				zfile.Name = convName(zfile.Name)
			}
			for _, fn := range fns {
				fn(zfile)
			}
		}
	}
	return
}
