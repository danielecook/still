# Schema

Schemas have the following structure:

* `directives` - specify global rules and settings.
* `column rules` - Define a set of expressions on which to evaluate each column.
* `data` __optional__ - a yaml-formatted dataset can be appended to the end of a schema to define value sets or other data used in validation.

## Directives

Directives are assigned using a `@` prefix and apply global settings.

### `@na_values`

Use `@na_values` to specify values to treat as `NA`. See [Handling Missing Data](### Handling missing data) for more details.

```yaml
@na_values NA
```

### `@empty_values`

Use `@empty_values` to specify values to treat as `empty`.

Values to treat as missing or `NA`. Use `""` for empty cells. See [Handling Missing Data](### Handling missing data) for more details.

```yaml
@empty_values "" NULL
```

### `@separater`

Sets the separater/delimiter for a data file. Do not quote the delimiter. Use `TAB` or `\t` for tab-delimited data.

```yaml
@separater: TAB
# comma-delimited
@separater: ,
```

`@sep` also works.

## Column Rules

Column Rules consist of a column name and expressions to test for each column. For example, the following tests that the color column must be equal to `red`, `blue` or `green`.

```yaml
color: any("red", "blue", "green")
```

## Data Providers

Certain rules are easier to specify if you need to compare a column against a larger set of data.

### YAML Data

Adding a dashline (`---`) signals the beginning of the data section of the schema. Any content below the dashline is parsed as YAML and can be accessed in expressions using its key. For example:

```yaml
color_values:
  - red
  - blue
  - green

flavor_values:
  - chocolate
  - vanilla
  - strawberry
```

Column rules might be specified like this:

```yaml
color: any(color_values)
flavor: any(flavor_values)
```

#### Functions supporting data providers

* `any`

## Missing Data

There are two types of missing data that `still` manages: 

### `empty`

`empty` data indicate missing data that may exist, but is currently unknown. This is similar to the `NULL` type in R. By default, `empty` values are defined as `""`.

### `NA`

`NA` values indicate known missing data. This is similar to the `NA` type in R. By default, `NA` values are defined as `NA`.

### Handling missing data

The difference between `NA` and `empty` is important for data validation purposes. There are cases where you might permit `empty` data, but not `NA` or vice versa. `@na_values` and `@empty_values` can be used to set values treated as `NA` and `empty`, respectively.


## Comments

You can add comments to your schema file using `//` or `/* block */`. For example:

```yaml
// This is a comment
color: any("blue", "red", "green") // expression for the color column
/*
    Using a block comment is fun
*/
```
