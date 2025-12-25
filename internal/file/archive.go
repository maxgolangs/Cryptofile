package file

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ArchiveDirectory(source string) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		arcName := filepath.ToSlash(relPath)

		if info.IsDir() {
			arcName += "/"
			_, err := zw.Create(arcName)
			return err
		}

		zf, err := zw.Create(arcName)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = CopyFileBuffered(zf, f)
		return err
	})

	if err != nil {
		zw.Close()
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ExtractDirectory(data []byte, destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	destAbs, err := filepath.Abs(destination)
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		targetPath := filepath.Join(destAbs, f.Name)

		targetAbs, err := filepath.Abs(targetPath)
		if err != nil {
			return err
		}
		if !strings.HasPrefix(targetAbs, destAbs) {
			return fmt.Errorf("обнаружен небезопасный путь в архиве: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		fileMode := f.FileInfo().Mode()
		fileMode = fileMode & 0777

		if fileMode&0700 == 0 {
			fileMode = 0644
		}

		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = CopyFileBuffered(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

