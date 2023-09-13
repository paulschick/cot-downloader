# COT Downloader

A utility to download, extract, and process Commitments of Traders (COT) report data.

## Features

- Download COT reports from specified years.
- Extract data from downloaded archives.
- Rename and process extracted files.
- Can be used as a standalone CLI tool or integrated into other Go projects as a library.

## Getting Started

## As a CLI Tool:

1. Navigate to the project root.
2. Run `go build -o cot-downloader ./cmd/cot-downloader/` to compile the CLI tool.
3. Use the generated cot-downloader binary to interact with the tool. For example:

```shell
./cot-downloader --start 2001 --end 2023
```

### As a Library:

Import the cot-downloader package into your project.

```go
import (
	"github.com/paulschick/cot-downloader/pkg/downloader"
	"github.com/paulschick/cot-downloader/pkg/extractor"
)
```

## Command Line Options

- `--report` - The report type. **Currently only supports Legacy Futures reports**. (Default `legacy`)
- `--start` - The year to start downloading COT reports from. (Default `2001`)
- `--end` - The year to stop downloading COT reports from. (Default `2023`)
- `--output` - The directory to save downloaded files to. (Default: `./output`)
- `--delay` - The delay in milliseconds between each download request. (Default: `500`)
