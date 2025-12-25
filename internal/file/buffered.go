package file

import (
	"bufio"
	"io"
	"os"
	"runtime"
)

func WriteFileBuffered(filename string, data []byte, perm os.FileMode) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 256*1024)
	defer writer.Flush()

	_, err = writer.Write(data)
	return err
}

func CopyFileBuffered(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 256*1024)
	return io.CopyBuffer(dst, src, buf)
}

func GetOptimalWorkers() int {
	workers := runtime.NumCPU() * 2
	if workers > 16 {
		workers = 16
	}
	if workers < 4 {
		workers = 4
	}
	return workers
}


