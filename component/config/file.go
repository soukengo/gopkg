package config

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type FileConfig struct {
	Dir string
}

type FileSource struct {
	err        error
	configType string
	file       *os.File
}

func NewFileSource(filePath string) Source {
	configType := strings.TrimPrefix(path.Ext(filePath), ".")
	s := &FileSource{configType: configType}
	fi, err := os.Stat(filePath)
	if err != nil {
		s.err = err
		return s
	}
	if fi.IsDir() {
		err = fmt.Errorf("%s is not a directory", filePath)
		return s
	}
	file, err := os.Open(filePath)
	if err != nil {
		s.err = err
		return s
	}
	s.file = file
	return s
}

func (s *FileSource) From() (io.Reader, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.file, nil
}

func (s *FileSource) ConfigType() string {
	return s.configType
}
func (s *FileSource) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
