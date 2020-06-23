# Expressions

### Background

Expressions all return a boolean (true/false) and allow you to evaluate conditions on a column. For brevity, test functions are implicitely passed the column being evaluated as the first argument. For example:

```yaml
status: contains("# NOTE")
# Converted to:
status: contains(status, "# NOTE")
```

You can still be explicit when referencing the column of interest, and you can also combine test functions using different columns:

```yaml
status: contains("# NOTE") && contains(color, "red")
# Converted to:
status: contains(status, "# NOTE") && contains(time, "# NOTE")
```

### Operators

`still` uses [Knetic/govaluate](https://github.com/Knetic/govaluate) to evaluate expressions. See [the manual](https://github.com/Knetic/govaluate/blob/master/MANUAL.md) for more detail on operators. The following operators are supported.

* Modifiers: `+` `-` `/` `*` `&` `|` `^` `**` `%` `>>` `<<`
* Comparators: `>` `>=` `<` `<=` `==` `!=` `=~` `!~`
* Logical ops: `||` `&&`
* Numeric constants, as 64-bit floating point (`12345.678`)
* String constants (single quotes: `'foobar'`)
* Date constants (__single quotes__, using any RFC3339, ISO8601, ruby date, or unix date; date parsing is automatically tried with any string constant)
* Boolean constants: `true` `false`
* Parenthesis to control order of evaluation `(` `)`
* Arrays (anything separated by `,` within parenthesis: `(1, 2, 'foo')`)
* Prefixes: `!` `-` `~`
* Ternary conditional: `?` `:`
* Null coalescence: `??`

#### Dates

Single quoted dates are parsed...

### Basic

##### `is`

Tests whether a value matches.

```yaml
color: is("red")
```

##### `not`

Tests whether a value does not match.

```yaml
is_passed: not("fail")
```

### Sets

##### `any`

Tests whether a value matches any of passed arguments.

```yaml
color: any("red", "blue", "green")
```

You can also use an array with the `IN` operator, but you must specify the column name:

```yaml
color: color IN ("red", "blue", "green")
```

##### `unique`

Tests whether a column is unique.

```yaml
items: unique()
# You can also test that a set of columns are unique.
color: unique(size, weight) # remember that color is implicit
```

!!! note

    `unique` does not work well on large datasets. It stores a hash digest of the arguments to test for uniqueness.

##### `is_subset_list`

```
is_subset_list(group_set, column_delimiter)
```

* `group_set` - A comma-delimited set of values.
* `column_delimiter` - The delimiter for the specified column you are testing.

Tests whether a delimited list (nested data) is a subset of the specified delimited list. For example:


```
letters: is_subset_list("A,B,C", ",") # If letters='A,B' --> TRUE
letters: is_subset_list("A,C", ",")   # If letters='A'B' --> FALSE
```

### Numbers

##### `range`

```
range(lower, upper)
```

Tests whether a value falls between `lower` and `upper` inclusive.

```yaml
rating: range(0,10)
```

##### `is_positive`

Tests whether a value is positive.

##### `is_negative`

Tests whether a value is negative.

### Strings

##### `contains`

```
contains(substring)
```

Tests for the presence of a substring in a value.

##### `regex`

```js
regex(expression)
```

Tests whether a value matches a regular expression.

You can also use `=~` or `!~` regex comparators.

```js
(colname =~ "L[0-9]+")
```

##### `uppercase`

Tests whether a value is all uppercase.

##### `lowercase`

Tests whether a value is all lowercase.

##### `replace(string, find, replace)`

Replace `find` with `replace` in `string`.

##### `length(low, high = None)`

Tests for string length in a given column.

```yaml
# Test for exact length match
column: length(10)

# Test for minimum length
column: length(10, "*")

# Test for range of lengths
column: length(10, 20)
```

##### `is_url`

Tests whether a string is a valid URL

#### Types

##### `is_numeric`
##### `is_int`

Test for numeric or integer values.

##### `is_bool`

Tests that column contains `true`, `TRUE`, `false`, or `FALSE`

#### Dates

##### `is_date`

Checks whether a value is date-like using strict criteria. `is_date` will fail on ambiguous date strings. For example, `02/03/2020` is interpretted differently in Europe vs. the US but `2020-02-03` is not.

##### `is_date_relaxed`

Checks whether a value is date-like with potential ambiguity. `02/03/2020` would pass.

##### `is_date_format(format)`

Check whether a column matches a specified date format. Formats can be specified as any date like this:

```
September 17, 2012, 10:10:09
oct. 7, 70
8/8/1965 12:00:00 AM
2006-01-02T15:04:05+0000
2014-04-26
```

```
date_of_birth: is_date_format("June 10, 1987")
```

See [araddon/dateparse](https://github.com/araddon/dateparse/blob/master/example/main.go#L12) for more examples.

__Important!__ You need to escape date values containing dashes using brackets (`[]`) or a double backslash (`\\`). For example `2020-02-10` is escaped like this:

```yaml
collection_date: is_date_format("[2020-02-10]")
collection_date: is_date_format("2020\\-02\\-10")
```

#### Files

##### `file_exists`

Checks whether a file exists

##### `file_min_size(fsize)`

Checks whether a file is a minimum size

```yaml
photo: file_min_size("1MB")
```

`fsize` is a size string such as `100 mb`, `1G`, or `500`.

##### `mimetype(type)`

Validates the mimetype for a given file. See [gabriel-vasile/mimetype](https://github.com/gabriel-vasile/mimetype/blob/master/supported_mimes.md) for available mimetypes.

```yaml
photo: mimetype("image/jpeg")
```