package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// 给定文件的绝对路径创建文件
func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 给定文件的绝对路径创建文件
func Create(path string) error {
	dir := Dir(path)
	if !Exists(dir) {
		Mkdir(dir)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// 打开文件
func Open(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 打开文件
func OpenWithFlag(path string, flag int) (*os.File, error) {
	f, err := os.OpenFile(path, flag, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 获取文件或目录信息
func Info(path string) *os.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return &info
}

// 修改时间(秒)
func MTime(path string) int64 {
	f, e := os.Stat(path)
	if e != nil {
		return 0
	}
	return f.ModTime().Unix()
}

// 修改时间(毫秒)
func MTimeMillisecond(path string) int64 {
	f, e := os.Stat(path)
	if e != nil {
		return 0
	}
	return int64(f.ModTime().Nanosecond() / 1000000)
}

// 文件大小(bytes)
func Size(path string) int64 {
	f, e := os.Stat(path)
	if e != nil {
		return 0
	}
	return f.Size()
}


// 文件移动/重命名
func Move(src string, dst string) error {
	return os.Rename(src, dst)
}

// 文件移动/重命名
func Rename(src string, dst string) error {
	return Move(src, dst)
}

// 文件复制
func Copy(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	err = dstFile.Sync()
	if err != nil {
		return err
	}
	srcFile.Close()
	dstFile.Close()
	return nil
}

// 文件/目录删除
func Remove(path string) error {
	return os.RemoveAll(path)
}

// 文件是否可
func IsReadable(path string) bool {
	result := true
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		result = false
	}
	file.Close()
	return result
}


// 修改文件/目录权限
func Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// 打开目录，并返回其下一级子目录名称列表，按照文件名称大小写进行排序
func ScanDir(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}

	list, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil
	}
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	return list
}

// 将所给定的路径转换为绝对路径
// 并判断文件路径是否存在，如果文件不存在，那么返回空字符串
func RealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !Exists(p) {
		return ""
	}
	return p
}

// (文本)读取文件内容
func GetContents(path string) string {
	return string(GetBinContents(path))
}

// (二进制)读取文件内容
func GetBinContents(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

// 写入文件内容
func putContents(path string, data []byte, flag int, perm os.FileMode) error {
	// 支持目录递归创建
	dir := Dir(path)
	if !Exists(dir) {
		Mkdir(dir)
	}
	// 创建/打开文件
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(data)
	if err != nil {
		return err
	} else if n < len(data) {
		return io.ErrShortWrite
	}
	return nil
}

// Truncate
func Truncate(path string, size int) error {
	return os.Truncate(path, int64(size))
}

// (文本)写入文件内容
func PutContents(path string, content string) error {
	return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

// (文本)追加内容到文件末尾
func PutContentsAppend(path string, content string) error {
	return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

// (二进制)写入文件内容
func PutBinContents(path string, content []byte) error {
	return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

// (二进制)追加内容到文件末尾
func PutBinContentsAppend(path string, content []byte) error {
	return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

// 获取当前执行文件的绝对路径
func SelfPath() string {
	p, _ := filepath.Abs(os.Args[0])
	return p
}

// 获取当前执行文件的目录绝对路径
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// 获取指定文件路径的文件名称
func Basename(path string) string {
	return filepath.Base(path)
}

// 获取指定文件路径的目录地址绝对路径
func Dir(path string) string {
	return filepath.Dir(path)
}

// 获取指定文件路径的文件扩展名
func Ext(path string) string {
	return filepath.Ext(path)
}


