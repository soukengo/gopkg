package env

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	goBuildFlags = []string{"go_build", "go-build", "Test", "___Benchmark", "GoLand/___"}
)

func RootDir() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	for _, flag := range goBuildFlags {
		if strings.Contains(path, flag) {
			return os.Getwd()
		}
	}
	return getParentDirectory(filepath.Dir(path)), nil
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(directory string) string {
	return substr(directory, 0, strings.LastIndex(directory, string(os.PathSeparator)))
}
