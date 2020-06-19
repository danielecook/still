# still

`still` is a program for validating tabular data from CSV, TSV, and Excel.

## Quick Start

Generate a schema. Directives start with `@` and refer to global options. These are followed by column names and test expressions.

__`cars.schema`__
```
@separater TAB
mpg: is_numeric()
cyl: range(2,8)
hp: is_positive() && range(10, 500)
vs: is(0) || is(1)
am: any(0, 1)
```

Then run the command line tool:

```
still validate cars.schema cars.tsv
```

# Test Expressions

## Vector Functions

## Utility Functions

Utility functions operate at the row level.

##### `if_else(condition, expr_if_true, expr_if_false)`

The `if_else` function evaluates an expression and returns `expr_if_true` or `expr_if_false` depending on the evaluation of `condition`.

```
column: if_else(column == "AB1", height > 50, height < 50)
```

#####  `max(...)`

Returns the max of arguments passed; Operates at the __row__ level.

```
orders: max(1,2,3) == 3         # Returns 3 == 3; true
items: max(col1, col2, col3)    # Returns the max value for the given row of col1-3.
inventory: max(col1)            # This does not return the max for an entire column; col1 is a scaler value.
```

__`min(...)`__

Returns the minimum of arguments passed; Operates at the __row__ level.

## Test Functions

### Background

Test functions all return a boolean (true/false) and allow you to evaluate conditions on a column. For brevity, they are implicitely passed the specified column when none is specified. For example:

```
status: contains("# NOTE")
# Converted to:
status: contains(status, "# NOTE")
```

You can still be explicit when referencing the column of interest, and you can also combine test functions using different columns:

```
status: contains("# NOTE") && contains(color, "red")
# Converted to:
status: contains(status, "# NOTE") && contains(time, "# NOTE")
```

### Basic

__`is(value)`__

Tests whether a value matches.

```
color: is("red")
```

__`not(value)`__

Tests whether a value does not match.

```
is_passed: not("fail")
```

### Sets

__`any(...)`__

Tests whether a value matches any of passed arguments.

```
color: any("red", "blue", "green")
```

### Numbers

__`range(lower, upper)`__

Tests whether a value falls between `lower` and `upper` inclusive.

```
rating: range(0,10)
```

__`is_positive()`__
__`is_negative()`__

Tests whether a value is positive or negative.


### Strings

__`contains(substr)`__

Tests for the presence of a substring in a value.

__`regex(expression)`__

Tests whether a value matches a regular expression.

__`uppercase()`__

Tests whether a value is all uppercase.

__`lowercase()`__

Tests whether a value is all lowercase.

__`length(low, high = None)`__

Tests for string length in a given column.

```
# Test for exact length match
column: length(10)

# Test for minimum length
column: length(10, "*")

# Test for range of lengths
column: length(10, 20)
```

#### Types

__`is_numeric()`__
__`is_int()`__

Test for numeric or integer values.

__`is_bool()`__

Tests that column contains `true`, `TRUE`, `false`, or `FALSE`

#### Dates

__`is_date()`__

Checks whether a value is date-like using strict criteria. `is_date()` will fail on ambiguous date strings. For example, `02/03/2020` is interpretted differently in Europe vs. the US but `2020-02-03` is not.

__`is_date_relaxed()`__

Checks whether a value is date-like with potential ambiguity. `02/03/2020` would pass.

__`is_date_format(format)`__

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

```
collection_date: is_date_format("[2020-02-10]")
collection_date: is_date_format("2020\\-02\\-10")
```

## Notes

`still` is largely influenced by [csv-validate](http://digital-preservation.github.io/csv-validator/), but offers more flexibility when validating tabular data.
