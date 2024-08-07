package file

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type ClientFile interface {
	Open(filePath string) error
	CheckFileSize() error
	GetFileSize() int64
	GetFileName()string
	Read(b []byte) (n int, err error)
	Close() error
}

// 100MB
const MaxFileSize = 104857600

func NewFile()ClientFile{
	return &clientFileMng{}
}

type clientFileMng struct {
	filePath string
	file *os.File
}



func (f *clientFileMng) Open(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	f.file = file
	f.filePath = filePath
	return nil
}

func (f *clientFileMng) CheckFileSize() error {
	fi, err := f.file.Stat()
	if err != nil {
		return err
	}

	log.Printf("file size is %v\n", fi.Size())

	if fi.Size() > MaxFileSize {
		return errors.New("file size too big")
	}

	return nil
}

func (f *clientFileMng) GetFileSize() int64 {
	fi, err := f.file.Stat()
	if err != nil {
		return 0
	}
	return fi.Size()
}

func (f *clientFileMng) GetFileName() string {
	return filepath.Base(f.filePath)
}

func(f *clientFileMng) Read(b []byte) (n int, err error) {
	return f.file.Read(b)
}

func (f *clientFileMng)Close ()error {
	return f.file.Close()
}