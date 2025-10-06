package archive

import (
	"archive/zip"
	"peach/utils"
)

// ConvName 解决中文文件名乱码的问题
func ConvName(name string) string {
	b := []byte(name)
	if d, err := utils.Decode(b); err == nil {
		return d
	}
	return name
}

// ExtractZip 读取 .zip 压缩包中的文件
func ExtractZip(file string, fns ...ReadFileFunc) (err error) {
	if r, err := zip.OpenReader(file); err == nil {
		defer r.Close()
		for _, zfile := range r.File {
			if zfile.NonUTF8 {
				zfile.Name = ConvName(zfile.Name)
			}
			for _, fn := range fns {
				fn(zfile)
			}
		}
	}
	return
}
