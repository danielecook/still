# Functions

Functions operate at the **row** level.

##### `if_else`

```js
if_else(condition, expr_if_true, expr_if_false)
// A ternary operator is also supported.
(condition ? expr_if_true : expr_if_false)
```

* `condition` - An expression that evaluates to `TRUE` or `FALSE`.
* `expr_if_true` - Value to return if `condition=-TRUE`
* `expr_if_false` - Value to return if `condition==FALSE`

The `if_else` function evaluates an expression and returns `expr_if_true` or `expr_if_false` depending on the evaluation of `condition`.

```js
flavor: if_else(flavor == "Chocolate", sugar > 50, sugar < 50)
```

##### `max`

Returns the max of arguments passed

```yaml
orders: max(1,2,3) == 3         # Returns 3 == 3; true
items: max(col1, col2, col3)    # Returns the max value for the given row of col1-3.
inventory: max(col1)            # This does not return the max for an entire column; col1 is a scaler value.
```
 
##### `min`

Returns the minimum of arguments passed

##### `to_lower`

Converts a string to lowercase.

##### `to_upper`

Converts a string to uppercase.

##### `replace`

```
replace(value, find, replace)
```

##### `count`

```js
count(column)
```

Returns the number of times the passed value(s) have been observed.

```yaml
color: count(color) <= 20 # Fails if a value is observed more than 20x times.
configuration: count(color, size) <= 10 # Fails if the combination of values is observed more than 10 times.
```

!!! note

    `count` does not work well on large datasets. It stores a hash digest of the arguments to test for uniqueness.

##### null coalescence

`??` can be used to set a default value.

```
(colname ?? 1) == 1 # returns TRUE if colname==NA/nil
```

## Two-Pass Functions

Two-pass functions allow for more advanced expressions to be evaluated, but they require a first-pass through the file to collect information.

##### `group_count`

```
group_count(group_column, count_column, eq_value)
```

`group_count` will group data by the `group_column`, and count the number of occurences of `eq_value` in the `count_column`.

__example__

```
family_id: is_int()
person: group_count(person, is_head_of_household, true) == 1
is_head_of_household: is_bool()
```

In the example above, we are checking to see that only one person in a family is set to `true` for the column `is_head_of_household`. If you have missing data in your grouping column you may need use an if_else statement to conditionally validate a row as true:

```
family_id: is_int()
person: if_else(is_missing(), true, group_count(person, is_head_of_household, true) == 1)
is_head_of_household: is_bool()
```