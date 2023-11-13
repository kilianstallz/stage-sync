## Tables Configuration

The `tables` section in the `config.yaml` file specifies the tables that you want to sync from the source database to the target database. Each table is represented as a dictionary with the following keys:

### `name`

This key specifies the name of the table in the database.

### `columns`

This key specifies an array of column names that you want to sync. Each item in the array is a string that represents the name of a column.

### `primaryKeys`

This key specifies an array of column names that make up the primary key of the table. Each item in the array is a string that represents the name of a primary key column.

### `onlyWhere`

This key specifies an array of conditions to filter the rows that you want to sync. Each condition is a dictionary with the following keys:

- `name`: The name of the column.
- `value`: The value to compare with.
- `type`: The data type of the column. This can be `string`, `int`, `float`, `bool`, etc.

Here is an example of a `tables` configuration:

```yaml
tables:
  - name: TextTemplates
    columns:
      - Id
      - LanguageId
      - Text
      - TextHeader
      - Format
    primaryKeys:
      - Id
      - LanguageId
    onlyWhere:
      - name: LanguageId
        value: ''
        type: string
  - name: LanguageItems
    columns:
      - LanguageId
      - LanguageKey
      - LanguageValue
      - LanguageItemType
    primaryKeys:
      - LanguageId
      - LanguageKey
