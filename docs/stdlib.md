# VimComplete Stdlib
Most stdlib modules load a macro into '@0', although some like `std:file_diff_assert`
are literal scripts performing operations on the current buffer.

- Parameters are passed in number registers in ascending order starting from @1
- Some macros can't complete all their operations due to the limitations of macro
execution model, and hence set post execution macros which must be called for correct
functioning
- Post-execution macros, if present, should be in @1
- Multiple return values can be returned through number registers much like parameters.
If a post-execution macro is required, then return values start from @2
- Stdlib macros should be in stored in @0 and should ensure that they don't
overwrite that register during execution to ensure macros can be repeated

**Usage Example**
```
\r std:add

o<value>ESC$v0"1y"_dd
o<value>ESC$v0"2y"_dd

@0

\\ return value in @1
```

## std:add

**Parameters**
- `1`: first operand
- `2`: second operand

**Returns:** Sum of the operands in register `@1`

Implemented using vim's `<C-a>`, which increments the number under cursor.
The second operand is used as the count for the motion `<C-a>` unless 0 is used, 
which is handled seperately.

Currently only works with natural numbers as negative integers can't be used as
counts.

```
\r std:add

o25^[$v0"1y"_dd
o15^[$v0"2y"_dd

@0

\\ 40
o^R1 
```

## std:assert_absent

**Parameters**
- `1`: Regex pattern that should have no matches in current buffer

Sets @1 to a no-op unless the pattern has a match, in which case it is set to 
a recursive ZQ macro which closes all windows.

`@1` must be called after execution if assertion fail is supposed to crash the program.
If some other behaviour is required then the contents of @1 can be matched, to detect
success or failure, and then operated on as required.

`@1` := "0" if ASSERTION HOLDS else "ZQ@1"

**Note:** Before checking for matches, the macro creates a new line at the end of the file
with `Go00`. Because the actual last line of the program has changes, if any search patterns
deal with end of the file with `\%$`, then they have to be adjusted accordingly.

For example, if the assertion was to check if the last line of the file isn't "end", then instead
of the search pattern being `^end\%$`, it would have have to be changed to `^end\n00\%$`.

```
\r std:assert_absent

\\ no empty lines present
o^\s*$
$v0"1y"_dd

@0
\\ post-execution handler
@1
```

## std:assert_present

**Parameters**
- `1`: Regex pattern that should have atleast one match

Sets @1 to a recursive ZQ macro, which closes all windows unless the pattern has a match, 
in which case it is set to a no-op.

`@1` must be called after execution if assertion fail is supposed to crash the program.
If some other behaviour is required then the contents of @1 can be matched, to detect
success or failure, and then operated on as required.

`@1` := "0" if ASSERTION HOLDS else "ZQ@1"

**Note:** Before checking for matches, the macro creates a new line at the end of the file
with `Go00`. Because the actual last line of the program has changes, if any search patterns
deal with end of the file with `\%$`, then they have to be adjusted accordingly.

For example, if the assertion was to check if the last line of the file is "end", then instead
of the search pattern being `^end\%$`, it would have have to be changed to `^end\n00\%$`.

```
\r std:assert_absent

\\ first line is init
o\%^init$
$v0"1y"_dd

@0
\\ post-execution handler
@1
```

## std:copy_register

**Parameters**
- `1`: source register name
- `2`: destination register name

Copies contents of registers specified in `@1` to register specified `@2`

```
oregister x
"xdd

\r std:copy_register

ax"1dl
ay"2dl

@0

\\ register x
"yp
```

## std:delete_pattern

**Operates on Buffer**

**Parameters**
- `1`: Regex pattern to delete from the buffer

Deletes all matches for regex pattern stored in `@1` from the current buffer.

Sets a post-execution cleanup macro in @1 which must be used to prevent unintended
changes to current buffer.

```
\r std:delete_pattern

\\ delete all empty lines
o^\s*$
$v0"1y"_dd

@0
\\ cleanup macro
@1
```

## std:file_diff_assert

**Parameters**
- `1`: File to check current buffer against

Asserts that the current buffer and the file in `@1` are equivalent. Sets `@1` to
a no-op if the files are equivalent; otherwise, it is set to a recursive ZQ macro
which closes all windows and exits.

This module doesn't load a macro for use; instead it directly runs a script.

Thus, the parameters must be set before the module is called.

```
oFAIL compiler:preprocess matchesTargetFile"zdd

\\ Set the parameter before the module is called
otest-data/preprocessed.vc
$v0"1y"_dd

\r std:file_diff_assert

\\ Post-execution handler
@1

oPASS compiler:preprocess matchesTargetFile"ydd
```

**Note:** Not completely reliable currently. Can give false positives in some
edge cases.

## std:left_trim

**Operates on Buffer**

Removes leading from all lines in the current buffer

```
\r std:left_trim

@0
```

## std:right_trim

**Operates on Buffer**

Removes trailing spaces from all lines in the current buffer

```
\r std:right_trim

@0
```

## std:search_and_replace

**Operates on Buffer**

**Parameters**
- `1`: Regex pattern to delete from the buffer
- `2`: Literal string to replace matched strings with

Replaces all matches for regex pattern stored in `@1` with the literal string `@2`

Sets a post-execution cleanup macro in @1 which must be used to prevent unintended
changes to current buffer.

```
\r std:search_and_replace

\\ lines containing the word eggs
o^.*eggs.*$
$v0"1y"_dd

o[redacted]
$v0"2y"_dd

@0
\\ cleanup macro
@1
```

## std:search_and_transform

**Operates on Buffer**

**Parameters**
- `1`: Regex pattern to delete from the buffer
- `2`: Macro to transform matched text with

Runs the macro `@2` over all matches for regex pattern stored in `@1`

Sets a post-execution cleanup macro in @1 which must be used to prevent unintended
changes to current buffer.

```
\r std:search_and_replace

\\ A literal \ or "
o[\\"]
$v0"1y"_dd 

\\ Escape by adding a '\' before
oi\l
$v0"2y"_dd 

@0
\\ cleanup macro
@1
```

## std:subtract

**Parameters**
- `1`: first operand
- `2`: second operand

**Returns:** Difference of @1 & @2 stored in @1

Implemented using vim's `<C-x>`, which decrements the number under cursor.
The second operand is used as the count for the motion `<C-x>` unless 0 is used, 
which is handled seperately.

Currently only works with natural numbers as negative integers can't be used as
counts.

```
\r std:subtract

o25^[$v0"1y"_dd
o15^[$v0"2y"_dd

@0

\\ 10
o^R1 
```

## std:trim

**Operates on Buffer**

Removes leading and trailing spaces from all lines in the current buffer

```
\r std:trim

@0
```
