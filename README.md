# COT Excel Processor

The COT Excel Processor is a tool that allows users to download and process COT Excel reports. It extracts specific 
columns of interest and exports to CSV.

## Features

- Download COT reports between specified years.
- Process and filter data columns from `.xls` files.
- Export processed data to `.csv` files.

## Getting Started

### Prerequisites

Ensure you have Go installed. This project was built with Go 1.16 but should work with other recent versions.

### Installation

Clone the repository:

```shell
git clone https://github.com/paulschick/cot-downloader.git
cd cot-downloader
```

Build the project:

```shell
go build .
```

This will create an executable named `cot-downloader` in your directory.

## Usage

Just process existing files:

```shell
./cot-downloader
```

Download reports for specific years and process:

```shell
./cot-downloader -download -startYear=2022 -endYear=2023
```

Specifying custom directories:

```shell
./cot-downloader -downloadDir="./customData" -outputDir="./customOutput"
```

### Command Line Options

- `-env`: Use a `.env` file instead of command line arguments. Defaults to true. This will ignore all other command line arguments.
- `-download`: Download reports before processing. Defaults to true. If false, you must have `.xls` files in the download directory.
- `-startYear`: Start year for report dates. Defaults to current year.
- `-endYear`: End year for report dates. Defaults to current year.
- `-downloadDir`: Directory for downloaded reports. Defaults to `./data/xls`.
- `-outputDir`: Directory for CSV exports. Defaults to `./data/csv`.
