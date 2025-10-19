package utils

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"iter"
)

// TarFile tar 包中的文件
type TarFile struct {
	*tar.Reader
	*tar.Header
}

// Open 打开文件
func (t *TarFile) Open() (io.ReadCloser, error) {
	return t, nil
}

// Close 关闭文件
func (t *TarFile) Close() error {
	return nil
}

// ExtractTar 读取 .tar.gz .tgz 压缩包中的文件，逐个返回文件的：文件名、文件接口
func (p *Path) IterTarfile() iter.Seq2[string, File] {
	return func(yield func(name string, file File) bool) {
		var (
			f   io.ReadCloser
			r   io.Reader
			h   *tar.Header
			err error
		)
		if f, err = p.Open(); err != nil {
			return
		}
		defer f.Close()

		switch {
		case p.HasExt(".tgz", ".tar.gz"):
			if r, err = gzip.NewReader(f); err != nil {
				return
			}
		case p.HasExt(".tar.bz2"):
			r = bzip2.NewReader(f)
		default:
			r = f
		}
		t := tar.NewReader(r)
		for h, err = t.Next(); err == nil; h, err = t.Next() {
			if !yield(h.Name, &TarFile{t, h}) {
				break
			}
		}
	}
}
