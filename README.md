# StageSync CLI

## Configuration

- Create a script.yaml file in the root directory of your project.

## Current Pitfalls

- The list of transfer columns must match the order of the ordinal columns in the source table.
- The primary Key columns must be included in the transfer columns.
- Only supported database is postgres
- Table Names need to match the source table name