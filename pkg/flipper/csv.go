package flipper

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func CsvHeaders() []string {
	return []string{"Flipper Flag", "Occurrence", "Authors", "Last Commit Time", "Source Paths"}
}

func WriteUsageStats(csvOutFile string, results map[string]*FlipperUsageStats) error {
	csvDev := os.Stdout
	if len(csvOutFile) > 0 {
		dev, err := os.Create(csvOutFile)
		if err != nil {
			return err
		}
		csvDev = dev
		defer csvDev.Close()

		fmt.Printf("Write results to %s\n", csvOutFile)
	}

	writer := csv.NewWriter(csvDev)
	defer writer.Flush()

	writer.Write(CsvHeaders())
	for flag, result := range results {
		authors := make([]string, 0, len(result.Authors))
		sourceFiles := make([]string, 0, len(result.SourceFiles))

		for author := range result.Authors {
			authors = append(authors, author)
		}
		for file := range result.SourceFiles {
			sourceFiles = append(sourceFiles, file)
		}

		timeStr := ""
		if !result.LastCommitTime.IsZero() {
			timeStr = result.LastCommitTime.Format(time.RFC3339)
		}
		writer.Write([]string{
			flag,
			strconv.Itoa(result.Occurrence),
			strings.Join(authors, ";"),
			timeStr,
			strings.Join(sourceFiles, ";"),
		})
	}

	return nil
}

func MergeUsageStats(csvFile string, usageStats map[string]*FlipperUsageStats) error {
	file, err := os.Open(csvFile)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	content, err := reader.ReadAll()
	if err != nil {
		return err
	}
	content = content[1:]
	for _, row := range content {
		if _, ok := usageStats[row[0]]; !ok {
			usageStats[row[0]] = NewFlipperUsageStats()
		}
		stat := usageStats[row[0]]

		occurrence, err := strconv.Atoi(row[1])
		if err != nil {
			return err
		}
		if occurrence == 0 {
			continue
		}
		stat.Occurrence += occurrence
		authors := splitCsvColumn(row[2])
		for _, author := range authors {
			stat.Authors[author] = true
		}

		lastCommitTime, err := time.Parse(time.RFC3339, row[3])
		if err != nil {
			return err
		}

		if lastCommitTime.After(stat.LastCommitTime) {
			stat.LastCommitTime = lastCommitTime
		}

		sourceFiles := splitCsvColumn(row[4])
		for _, file := range sourceFiles {
			stat.SourceFiles[file] = true
		}
	}

	return nil
}

func splitCsvColumn(data string) []string {
	tokens := strings.Split(data, ";")
	results := make([]string, len(tokens))
	if len(tokens) < 1 {
		return []string{}
	}
	for i, token := range tokens {
		results[i] = strings.TrimSpace(token)
	}

	return results
}
