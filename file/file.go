package file

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/datainq/rwmc"
)

func NewFileReader(file string) (f io.ReadCloser, err error) {
	f, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(file, ".7z") {
		f, err = zlib.NewReader(f)
	} else if strings.HasSuffix(file, ".gz") {
		f, err = gzip.NewReader(f)
	}
	return f, err
}

func MaybeAddCompression(file string, w io.WriteCloser) (io.WriteCloser, error) {
	if strings.HasSuffix(file, ".7z") {
		w1, err := zlib.NewWriterLevel(w, zlib.BestCompression)
		if err != nil {
			return w1, err
		}
		return rwmc.NewWriteMultiCloser(w1, w), nil
	} else if strings.HasSuffix(file, ".gz") {
		w1, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			return w1, err
		}
		return rwmc.NewWriteMultiCloser(w1, w), nil
	}
	return w, nil
}

func NewFileWriter(file string) (w io.WriteCloser, err error) {
	w, err = os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	if err != nil {
		return nil, err
	}
	return MaybeAddCompression(file, w)
}

// CopyFile copies the contents from src to dst atomically.
// If dst does not exist, CopyFile creates it with permissions perm.
// If the copy fails, CopyFile aborts and dst is preserved.
//
// Modified version of: https://go-review.googlesource.com/c/go/+/1591
func CopyFile(dst, src string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err = os.Chmod(tmp.Name(), perm); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err = os.Rename(tmp.Name(), dst); err != nil {
		os.Remove(tmp.Name())
	}
	return err
}
