package flipper

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/Xfers/flipper-usage/internal/config"
	"golang.org/x/sync/errgroup"
)

func readFlipperFlags(flagFile string) ([]string, error) {
	dev, err := os.Open(flagFile)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(dev)
	scanner.Split(bufio.ScanLines)
	flags := make([]string, 0, 256)

	for scanner.Scan() {
		flags = append(flags, strings.TrimSpace(scanner.Text()))
	}

	return flags, nil
}

func ScanFlipperFlags(opts config.Options) error {
	// read flipper flags from file
	flags, err := readFlipperFlags(opts.FlipperFlagFile)
	if err != nil {
		return nil
	}

	// read source files and build source file caches
	store := NewSourceFileStore(opts.ScanFolder, opts.FileSuffix)
	err = store.build()
	if err != nil {
		return nil
	}

	analyzer := NewFlipperAnalyzer()
	var usageCount atomic.Int32

	fmt.Printf("Check %d flipper flags, %d files\n", len(flags), len(store.files))

	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(-1)
	usageChan := make(chan *FlipperUsage, 1024)

	// start goroutines
	for _, flag := range flags {
		flag := flag
		for _, sourceFile := range store.files {
			sourceFile := sourceFile
			group.Go(func() error {
				usage := NewFlipperUsage()
				usage.Flag = flag
				if err := usage.CheckFlag(store, sourceFile); err != nil {
					return err
				}

				select {
				// pass processed usage results to channel
				case usageChan <- usage:
					return nil
				// capture if any error occurs.
				case <-ctx.Done():
					return ctx.Err()
				}
			})
		}
	}

	// create a goroutine to close out the channel when the first error occurs or when all tasks finished.
	go func() {
		group.Wait()
		close(usageChan)
	}()

	for usage := range usageChan {
		usageCount.Add(1)
		analyzer.Analyze(usage)
	}

	// check if any errors
	err = group.Wait()
	if err != nil {
		return err
	}

	fmt.Printf("Number of flipper usage analyzed: %d\n", usageCount.Load())

	err = WriteUsageStats(opts.CsvOutFile, analyzer.Results)
	if err != nil {
		return err
	}

	return nil
}
