package batch

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"cryptor/internal/file"
	"cryptor/pkg/decrypt"
	"cryptor/pkg/encrypt"
)

type ProcessResult struct {
	Result string
	Err    error
	Target string
}

func GatherTargets(directoryPath, mode string) ([]string, error) {
	var targets []string
	isEncrypt := mode == "encrypt"

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isEncrypt {
			if filepath.Ext(path) == ".encrypted" {
				return nil
			}
		} else {
			if filepath.Ext(path) != ".encrypted" {
				return nil
			}
		}

		targets = append(targets, path)
		return nil
	})

	return targets, err
}

func ProcessDirectoryParallel(directoryPath, password, mode string, removeOriginal bool) (processed, errors, total int) {
	targets, err := GatherTargets(directoryPath, mode)
	if err != nil {
		return 0, 1, 0
	}

	total = len(targets)
	if total == 0 {
		return 0, 0, 0
	}

	workers := file.GetOptimalWorkers()
	if total < workers {
		workers = total
	}

	jobs := make(chan string, workers)
	results := make(chan ProcessResult, total)

	var wg sync.WaitGroup
	var processedCount int64
	var errorCount int64

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range jobs {
				var result string
				var err error

				if mode == "encrypt" {
					result, err = encrypt.EncryptPath(target, password, removeOriginal)
				} else {
					result, err = decrypt.DecryptFile(target, password, removeOriginal)
				}

				if err != nil {
					atomic.AddInt64(&errorCount, 1)
					results <- ProcessResult{Target: target, Err: err}
				} else {
					atomic.AddInt64(&processedCount, 1)
					results <- ProcessResult{Target: target, Result: result, Err: nil}
				}
			}
		}()
	}

	go func() {
		for _, target := range targets {
			jobs <- target
		}
		close(jobs)
	}()

	completed := 0
	for completed < total {
		<-results
		completed++
	}

	wg.Wait()
	close(results)

	return int(atomic.LoadInt64(&processedCount)), int(atomic.LoadInt64(&errorCount)), total
}

