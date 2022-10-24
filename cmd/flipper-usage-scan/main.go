package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Xfers/flipper-usage/internal/config"
	"github.com/Xfers/flipper-usage/pkg/flipper"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: flipper-usage-scan -f FILE [options] scan_folder\n")
	flag.PrintDefaults()
}

func main() {
	startTime := time.Now()
	var opts config.Options

	flag.StringVar(&opts.FlipperFlagFile, "f", "", "flipper flags file")
	flag.StringVar(&opts.FileSuffix, "s", ".rb", "the file suffix to scan")
	flag.StringVar(&opts.CsvOutFile, "o", "", "the CSV file to store scan results")
	flag.Usage = usage

	flag.Parse()

	args := flag.Args()
	if len(opts.FlipperFlagFile) < 1 || len(args) < 1 {
		usage()
		os.Exit(1)
	}
	opts.ScanFolder = args[0]

	err := flipper.ScanFlipperFlags(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Scan flipper flags fail, reason: %s\n", err.Error())
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Execute duration: %.2f seconds.\n", elapsed.Seconds())
	os.Exit(0)
}
