package flipper

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type SourceFileStore struct {
	files      []string
	cacheMap   sourceCacheMap
	scanFolder string
	fileSuffix string
}

type sourceCache struct {
	mu     sync.Mutex
	reader *bytes.Reader
}

type sourceCacheMap map[string]*sourceCache

func NewSourceFileStore(scanFolder string, fileSuffix string) *SourceFileStore {
	var instance SourceFileStore
	instance.scanFolder = scanFolder
	instance.fileSuffix = fileSuffix
	instance.files = make([]string, 0, 4096)

	return &instance
}

func (s *SourceFileStore) build() error {
	// retrieve source files
	if stat, err := os.Stat(s.scanFolder); err != nil || !stat.IsDir() {
		return fmt.Errorf("scan directory doesn't exist: %s", s.scanFolder)
	}

	err := filepath.Walk(s.scanFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() || filepath.Ext(path) != s.fileSuffix {
			return nil
		}
		s.files = append(s.files, path)
		return nil
	})

	if err != nil {
		return err
	}

	// build source cache
	fmt.Printf("Building %d source file caches.\n", len(s.files))
	s.cacheMap = make(sourceCacheMap, len(s.files))
	for _, file := range s.files {
		dev, err := os.Open(file)
		if err != nil {
			return err
		}
		defer dev.Close()

		buf, err := io.ReadAll(dev)
		if err != nil {
			return err
		}
		cache := new(sourceCache)
		cache.reader = bytes.NewReader(buf)
		s.cacheMap[file] = cache
	}

	return nil
}
