package flipper

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

type AttrSet map[string]bool

type FlipperUsage struct {
	SourcePath     string
	Flag           string
	Occurrence     int
	Authors        AttrSet
	LastCommitTime time.Time
}

func NewFlipperUsage() *FlipperUsage {
	var instance FlipperUsage
	instance.Authors = make(AttrSet, 1)
	return &instance
}

func (f *FlipperUsage) CheckFlag(store *SourceFileStore, filePath string) error {
	cache, ok := store.cacheMap[filePath]
	if !ok {
		return fmt.Errorf("can't find source caches for %s", filePath)
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()

	reader := cache.reader
	reader.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	relPath, err := filepath.Rel(store.scanFolder, filePath)
	if err != nil {
		return err
	}

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		if strings.Contains(scanner.Text(), f.Flag) == false {
			continue
		}
		f.Occurrence++

		log, err := GetGitShortLog(lineNo, store.scanFolder, relPath)
		if err != nil {
			return err
		}
		if _, ok := f.Authors[log.author]; !ok {
			f.Authors[log.author] = true
		}
		if log.authorDate.After(f.LastCommitTime) {
			f.LastCommitTime = log.authorDate
		}
	}
	if f.Occurrence > 0 {
		f.SourcePath = relPath
	}

	return nil
}

type FlipperUsageStats struct {
	Occurrence     int
	Authors        AttrSet
	SourceFiles    AttrSet
	LastCommitTime time.Time
}

func NewFlipperUsageStats() *FlipperUsageStats {
	var stats FlipperUsageStats
	stats.Occurrence = 0
	stats.Authors = make(AttrSet, 1)
	stats.SourceFiles = make(AttrSet, 1)
	return &stats
}

type FlipperAnalyzer struct {
	Results map[string]*FlipperUsageStats
}

func NewFlipperAnalyzer() *FlipperAnalyzer {
	var instance FlipperAnalyzer
	instance.Results = make(map[string]*FlipperUsageStats, 4096)
	return &instance
}

func (f *FlipperAnalyzer) Analyze(usage *FlipperUsage) {
	if _, ok := f.Results[usage.Flag]; !ok {
		f.Results[usage.Flag] = NewFlipperUsageStats()
	}

	if usage.Occurrence < 1 {
		return
	}

	result := f.Results[usage.Flag]

	result.Occurrence += usage.Occurrence
	result.SourceFiles[usage.SourcePath] = true
	for author := range usage.Authors {
		result.Authors[author] = true
	}
	if usage.LastCommitTime.After(result.LastCommitTime) {
		result.LastCommitTime = usage.LastCommitTime
	}
}
