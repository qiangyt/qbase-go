package comm

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

func FallbackFilePath(fallbackDir string, url string) string {
	sumBytes := sha256.Sum256([]byte(url))
	sumText := fmt.Sprintf("%x", sumBytes)
	return filepath.Join(fallbackDir, sumText)
}

func HasFallbackFile(fallbackDir string, fs afero.Fs, url string) (bool, error) {
	fallbackFilePath := FallbackFilePath(fallbackDir, url)
	return FileExists(fs, fallbackFilePath)
}

func ReadFallbackFile(fallbackDir string, fs afero.Fs, url string) (string, []byte, error) {
	fallbackFilePath := FallbackFilePath(fallbackDir, url)

	exists, err := FileExists(fs, fallbackFilePath)
	if err != nil {
		return "", nil, err
	}
	if !exists {
		return "", nil, nil
	}
	bytes, err := ReadFileBytes(fs, fallbackFilePath)
	return fallbackFilePath, bytes, err
}

func WriteFallbackFile(fallbackDir string, fs afero.Fs, url string, bytes []byte) (string, error) {
	fallbackFilePath := FallbackFilePath(fallbackDir, url)

	exists, err := FileExists(fs, fallbackFilePath)
	if err != nil {
		return "", err
	}
	if exists {
		err = RemoveFile(fs, fallbackFilePath)
		if err != nil {
			return "", err
		}
	}
	return fallbackFilePath, WriteFile(fs, fallbackFilePath, bytes)
}
