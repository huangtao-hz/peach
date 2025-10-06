package archive

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"strings"
)

// File 定义压缩包文件，正常的 Path 也符合这个规范
type File interface {
	FileInfo() fs.FileInfo
	Open() (io.ReadCloser, error)
}

// ReadFileFunc 定义读取文件的函数，该函数应先判断文件名，然后再处理
type ReadFileFunc func(File)

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

// ExtractTar 读取 .tar.gz .tgz 压缩包中的文件
func ExtractTar(file string, fns ...ReadFileFunc) (err error) {
	var (
		f io.ReadCloser
		r io.Reader
		h *tar.Header
	)

	if f, err = os.Open(file); err != nil {
		return
	}
	defer f.Close()

	switch {
	case strings.HasSuffix(file, ".tgz") || strings.HasSuffix(file, ".tar.gz"):
		if r, err = gzip.NewReader(f); err != nil {
			return
		}
	case strings.HasSuffix(file, ".tar.bz2"):
		r = bzip2.NewReader(f)
	default:
		r = f
	}

	t := tar.NewReader(r)
	for h, err = t.Next(); err == nil; h, err = t.Next() {
		for _, fn := range fns {
			fn(&TarFile{t, h})
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}
