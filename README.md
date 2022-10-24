# Flipper usage analysis utilities

## Prerequisites
* Golang version: 1.16+
* Git CLI: 2.25+

## Install
Via `go install`
```
go install github.com/Xfers/flipper-usage/cmd/flipper-usage-scan@latest
go install github.com/Xfers/flipper-usage/cmd/flipper-merge-result@latest
```

Build from source
```
make
# binaries locates in ./build
```

## Usage
**Generate flipper flags file first**

Use `./script/dump_flipper_flags.sh` to access flipper database and generate flipper flag list to text file.
```
Usage: dump_flipper_flags.sh db_user db_password db_host db_name output.txt

# Example
./script/dump_flipper_flags.sh test password localhost xfersFeatureFlag /tmp/flags.txt
```
> The format of flag list is newline-separated, one line on flag.

### Scan repo: flipper-usage-scan
Scan the give source folder with specific file suffix and flipper flags file
```
Usage: flipper-usage-scan -f FILE [options] scan_folder
  -f string
        flipper flags file
  -o string
        the CSV file to store scan results
  -s string
        the file suffix to scan (default ".rb")
# Example
flipper-usage-scan -f /tmp/flags.txt -o /tmp/analyze_result.csv -s ".rb" /repo
```

### Merge results: flipper-merge-result
Merge multiple analysis results to one file
```
Usage: flipper-merge-result -o output_file csv_files...
  -o string
        the output CSV file

# Example
flipper-merge-result ./merged_results.csv ./project_1.csv ./project_2.csv ./project_3.csv
```
