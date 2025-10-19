package utils

import (
	"archive/zip"
	"iter"
)

// ConvName 解决中文文件名乱码的问题
func ConvName(name *string) {
	if b := []byte(*name); IsGBK(b) {
		if d, err := Decode(b); err == nil {
			*name = d
		}
	}
}

// IterZip 遍历 .zip 文件，逐个返回文件的：文件名、文件接口
func (p *Path) IterZip() iter.Seq2[string, File] {
	return func(yield func(name string, file File) bool) {
		if r, err := zip.OpenReader(p.path); err == nil {
			defer r.Close()
			for _, zfile := range r.File {
				if zfile.NonUTF8 {
					ConvName(&zfile.Name)
				}
				if !yield(zfile.Name, zfile) {
					break
				}
			}
		}
	}
}
