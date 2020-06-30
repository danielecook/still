# Schema

Schemas have the following structure:

* `directives` - specify global rules and settings.
* `column rules` - Define a set of expressions on which to evaluate each column.
* `data` __optional__ - a yaml-formatted dataset can be appended to the end of a schema to define value sets or other data used in validation.

## Directives

Directives are assigned using a `@` prefix and apply global settings.

### `@na_values`

Use `@na_values` to specify values to treat as `NA`. See [Missing Data](#missing_data) for more details.

__Default__

```yaml
@na_values NA
```

### `@empty_values`

Use `@empty_values` to specify values to treat as `EMPTY`.

Use `""` for empty cells. See [Missing Data](#missing_data) for more details.

__Default__

```yaml
@empty_values "" NULL
```

### `@ordered`

Require that the columns appear in the same order as specified in the schema. 

```
@ordered
```

### `@fixed`

Require that column names match the schema in the same order.

```
@fixed
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

There are two types of missing data that `still` manages for additional flexibility. However, you can choose to treat all missing data as NA if desired. `NA` values represent "known" missing data. These are similar to `NA` values in R. `EMPTY` can be considered "unknown" missing data

To clarify further, consider a dataset on cars. The column `mpg` for all-electric vehicles would be labeled `NA` ("not applicable") as it does not apply. Another scenerio might be that you know the `name`, `make`, and `mpg` of a new vehicle but not the `color`. This flexibility would allow you to throw an error with an `NA` value for color, but not for an `EMPTY` value.

## Comments

You can add comments to your schema file using `//` or `/* block */`. For example:

```yaml
// This is a comment
color: any("blue", "red", "green") // expression for the color column
/*
    Using a block comment is fun
*/
```
