# Still

`still` is a program for validating tabular data from CSV, TSV, and Excel.

## Quick Start

Generate a schema. Directives start with `@` and refer to global options. These are followed by column names and test expressions.

__`cars.schema`__

```yaml
@separater TAB
mpg: is_numeric()
cyl: range(2,8)
hp: is_positive() && range(10, 500)
vs: is(0) || is(1)
am: any(0, 1)
```

Then run the command line tool:

```bash
still validate cars.schema cars.tsv
```
