# StageSync CLI

StageSync CLI is a command-line tool written in Go that synchronizes data between different stages of a database. 
It uses a YAML configuration file to specify the source and target databases, as well as the tables and columns to sync. 
The tool is designed to handle complex scenarios, such as syncing only specific rows, not deleting rows in the target that don't exist in the source, and more. 
It's a powerful tool for managing data propagation across different stages of your application.

# Usage

## Installation

### Install CLI from latest Github Release

[Latest Release](https://github.com/kilianstallz/stage-sync/releases/latest)

### Install using Go
```sh
  go install github.com/kilianstallz/stage-sync/cmd/stage-sync@latest
```

## Configuration

Define a new configuration by creating a new `config.yaml` file like in the [Example](/config_example.md).
Configure tables like in [the Table Description](/docs/table_config.md).


