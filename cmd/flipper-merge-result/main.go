package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Xfers/flipper-usage/pkg/flipper"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: flipper-merge-result -o output_file csv_files...\n")
	flag.PrintDefaults()
}

func main() {
	var outCsvFile string
	flag.StringVar(&outCsvFile, "o", "", "the output CSV file")
	flag.Usage = usage

	flag.Parse()

	args := flag.Args()
	if len(outCsvFile) < 1 || len(args) < 2 {
		usage()
		os.Exit(1)
	}

	usageStats := make(map[string]*flipper.FlipperUsageStats, 1024)
	for _, csvFile := range args {
		fmt.Printf("Merge %s\n", csvFile)
		err := flipper.MergeUsageStats(csvFile, usageStats)
		if err != nil {
			fmt.Printf("Merge failed, reason: %s\n", err.Error())
			os.Exit(1)
		}
	}

	err := flipper.WriteUsageStats(outCsvFile, usageStats)
	if err != nil {
		fmt.Printf("Write CSV file failed, reason: %s\n", err.Error())
		os.Exit(1)
	}
}
