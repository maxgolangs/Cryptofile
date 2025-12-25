package file

import (
	"fmt"
	"os"
	"path/filepath"

	"cryptor/internal/obfuscation"
)

func EnsureUniquePath(path string, isDir bool) string {
	candidate := path
	counter := 1

	for {
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}

		if isDir {
			candidate = fmt.Sprintf("%s_%d", path, counter)
		} else {
			ext := filepath.Ext(path)
			base := filepath.Base(path)
			dir := filepath.Dir(path)
			baseWithoutExt := base[:len(base)-len(ext)]
			candidate = filepath.Join(dir, fmt.Sprintf("%s_%d%s", baseWithoutExt, counter, ext))
		}
		counter++
	}
}

func IsEncryptedCandidate(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}

	if filepath.Ext(path) != ".encrypted" {
		return false
	}

	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	magicStr := obfuscation.GetMagic()
	magic := make([]byte, len(magicStr))
	if n, err := f.Read(magic); err != nil || n != len(magicStr) {
		return false
	}

	return string(magic) == magicStr
}

