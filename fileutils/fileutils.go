package fileutils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

const BufferSize = 8 * 1024 * 1024

func WriteFile(fpath string, data []byte) error {
	dir, err := filepath.Abs(filepath.Dir(fpath))
	if err != nil {
		return err
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fpath, data, 0666)
}

func WriteFileWithForce(fname string, data []byte) error {
	if strings.Index(fname, "~") == 0 {
		u, err := user.Current()
		if err != nil {
			return err
		}

		fname = u.HomeDir + fname[1:]
	}

	dir, err := filepath.Abs(filepath.Dir(fname))
	if err != nil {
		return err
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fname, data, 0666)
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func ListFileNames(path string) ([]string, error) {
	var (
		list = []string{}
	)

	fds, err := ioutil.ReadDir(path)
	if err != nil {
		return list, err
	}

	for _, f := range fds {
		list = append(list, f.Name())
	}
	return list, nil
}

func CreateDirectory(dirPath string) error {
	f, e := os.Stat(dirPath)
	if e != nil && os.IsNotExist(e) {
		return os.MkdirAll(dirPath, 0755)
	}
	if e == nil && !f.IsDir() {
		return fmt.Errorf("create dir:%s error, not a directory", dirPath)
	}
	return e
}

func DeleteFile(filePath string) error {
	if !PathExist(filePath) {
		return fmt.Errorf("delete file:%s error, file not exist", filePath)
	}
	if IsDir(filePath) {
		return fmt.Errorf("delete file:%s error, is a directory instead of a file", filePath)
	}
	return os.Remove(filePath)
}

func DeleteFiles(filePaths ...string) {
	if len(filePaths) > 0 {
		for _, f := range filePaths {
			DeleteFile(f)
		}
	}
}

func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	if PathExist(path) {
		return os.OpenFile(path, flag, perm)
	}
	if err := CreateDirectory(filepath.Dir(path)); err != nil {
		return nil, err
	}

	return os.OpenFile(path, flag, perm)
}

func Link(src string, linkName string) error {
	if PathExist(linkName) {
		if IsDir(linkName) {
			return fmt.Errorf("link %s to %s: error, link name already exists and is a directory", linkName, src)
		}
		if err := DeleteFile(linkName); err != nil {
			return err
		}

	}
	return os.Link(src, linkName)
}

func SymbolicLink(src string, target string) error {
	return os.Symlink(src, target)
}

func CopyFile(src string, dst string) (err error) {
	var (
		s *os.File
		d *os.File
	)
	if !IsRegularFile(src) {
		return fmt.Errorf("copy file:%s error, is not a regular file", src)
	}
	if s, err = os.Open(src); err != nil {
		return err
	}
	defer s.Close()

	if PathExist(dst) {
		return fmt.Errorf("copy file:%s error, dst file already exists", dst)
	}

	if d, err = OpenFile(dst, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755); err != nil {
		return err
	}
	defer d.Close()

	buf := make([]byte, BufferSize)
	for {
		n, err := s.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 || err == io.EOF {
			break
		}
		if _, err := d.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func MoveFile(src string, dst string) error {
	if !IsRegularFile(src) {
		return fmt.Errorf("move file:%s error, is not a regular file", src)
	}
	if PathExist(dst) && !IsDir(dst) {
		if err := DeleteFile(dst); err != nil {
			return err
		}
	}
	return os.Rename(src, dst)
}

func MoveFileAfterCheckMd5(src string, dst string, md5 string) error {
	if !IsRegularFile(src) {
		return fmt.Errorf("move file with md5 check:%s error, is not a "+
			"regular file", src)
	}

	m := Md5Sum(src)
	if m != md5 {
		return fmt.Errorf("move file with md5 check:%s error, md5 of source "+
			"file doesn't match against the given md5 value", src)
	}
	return MoveFile(src, dst)
}

func PathExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func IsDir(name string) bool {
	f, e := os.Stat(name)
	if e != nil {
		return false
	}
	return f.IsDir()
}

func IsRegularFile(name string) bool {
	f, e := os.Stat(name)
	if e != nil {
		return false
	}
	return f.Mode().IsRegular()
}

func Md5Sum(name string) string {
	if !IsRegularFile(name) {
		return ""
	}

	f, err := os.Open(name)
	if err != nil {
		return ""
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, BufferSize)
	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetSys(info os.FileInfo) (*syscall.Stat_t, bool) {
	sys, ok := info.Sys().(*syscall.Stat_t)
	return sys, ok
}
