# Schema

Schemas begin with a set of `directives` which define global rules and settings.

This is followed by rules for each column.

## Directives

Directives are assigned using a `@` prefix and apply global settings.

### `na_values`

Values to treat as missing or `NA`. Use `""` for empty cells.

```yaml
@na_values NA NULL ""
```

### `@separater`

Sets the separater/delimiter for a data file. Do not quote the delimiter. Use `TAB` or `\t` for tab-delimited data.

```yaml
@separater: TAB
# comma-delimited
@separater: ,
```

`@sep` also works.

## Column Definitions

Column Definitions consist of a column and expressions for testing conditions on that column. For example, the following tests that the color column must be equal to red, blue or green.

```yaml
color: any("red", "blue", "green")
```
