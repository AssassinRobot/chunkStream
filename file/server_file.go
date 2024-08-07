package file

import (
	"os"
	"path/filepath"
)

type ServerFile interface {
	SetFile(fileName, path string) error
	Write(chunk []byte) error
	Close() error
}

func NewServerFile()ServerFile{
	return &serverFileMng{}
}

type serverFileMng struct {
	filePath   string
	outputFile *os.File
}

func (f *serverFileMng) SetFile(fileName, path string) error {
	if f.filePath != "" {
		return nil
	}
	
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	f.filePath = filepath.Join(path, fileName)

	file, createErr := os.Create(f.filePath)
	if createErr != nil {
		return err
	}

	f.outputFile = file

	return nil
}

func (f *serverFileMng) Write(chunk []byte) error {
	if f.outputFile == nil {
		return nil
	}
	_, err := f.outputFile.Write(chunk)
	return err
}

func (f *serverFileMng) Close() error {
	return f.outputFile.Close()
}
