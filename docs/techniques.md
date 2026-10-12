# Techniques

Some of the techniques used to implement various features in this repo are documented here
to make it easier to read the source code.

## 1. Macro Definitions

Normally, macros in Vim are defined or set using the syntax `q<reg>...q`. This method
is ineffective for VimComplete programs because for one it cannot accommodate more
complicated macros which define other temporary macros for internal use and more pertinently
recording macros means that a macro must be executed once, which is a side-effect on top of
a pure definition. Writing raw keystrokes to a line and then copying that line into the
required register is preferred.

```
o
  a<keystrokes>
  a<keystrokes>
  a<keystrokes>
$v0"<reg>y"_dd
```

The following scaffold is used for macro definitions in the source code. The lone `o`
followed by indented append chains is purely for readability. The core idea is just to
write keystrokes without affecting the buffer.

The final line for storing the keystrokes into the required register is important. Using a
line mode based motion like `dd` or `yy` appends a carriage return at the end of the macro,
and hence must not be used. The direction in which the line is selected is important as well;
if it were `0v$` instead of `$v0`, then the same carriage return problem would persist.
After the line is selected correctly, both `y` and `d` can be used to store the keystrokes.
Using delete doesn't offer much benefit however because the empty line would still need to be
deleted afterwards.

Writing raw keystrokes like escape characters to the buffer means that <C-v>
must be used. So if a macro wanted to escape out of insert mode, then an additional
<C-v> must be used before the <Esc>.

```
o
  a<keystrokes>
  \\ enter insert mode
  ai<text><C-v><Esc>
$v0"yy"_dd
```

If a helper macro was defined inside the macro definition itself then this idea would apply twice.
Thus an `<Esc>` would now become `<C-v><C-v><C-v><Esc>`.

## 2. Macro Recursion With Motion Failure

If a motion fails during a macro then the entire macro execution is terminated after that point.
For example a `f` search for a character not present in the current line or a `l` at the end of
a line will fail.

A macro like `Gjdd` will not delete the last line.

A motion failure can be used like a base case for a recursive function.

```
o
  a0d^
  aj
  a@a
$v0"ay"_dd

gg
@a
```

This macro applies `0d^` to all lines and automatically stops when it reaches the
end of the buffer because the second line `j` fails.

**Limitations**
If a macro calls another macro in its body, and that macro fails due to motion failure,
it terminates not just the inner macro, but also the outer macro. This makes recursive
macros implemented with motion failure unreliable when used as helpers inside other macros.

Note: This can also be used deliberately like in asserts where the failure should terminate
the entire call stack.

## 3. Iterating over all regex matches

The macro from the previous section essentially iterated over all lines in the current
buffer. The same idea could be replicated for all characters in the current line, because
both provide a sentinel value through EOF/EOL.
The final match of a regex in the buffer provides no such sentinel value, thus it must be
artificially produced by making the final line of the buffer a match for the current regex.

This can be achieved either by manually adding a match to your search term, like if searching
for all lines beginning with 'a' then do `Goa` to ensure that final search will land
on the last line.

A less naive approach can modify the search pattern itself to add another clause for matching
the EOF marker itself (`\%$`). This will work regardless of pattern and thus is used in stdlib
which accepts arguments for patterns and cannot manually produce a match for the regex pattern.

Note: The search must start from the start of the buffer to match everything until EOF. But if
the match is on the first line itself, it might not match because of how search positions
work in Vim. Thus a new line before the actual first line must be created to execute search and
should be cleaned up afterwards.

## 4. Regex Conditionals

The buffer can be setup in a particular way to perform different actions based on whether a search
fails or succeeds. For example, in `std:add`, the macro must check if the second argument is '0', which
is an edge case that doesn't work with the regular <C-a> approach and thus requires separate processing.

The buffer is set up such that the last line is '0' followed by the required macro as the
next word (The macro in this case is a no-op because adding 0 requires no change to the first argument).
The line before that is just the argument `@2` followed by the standard <C-a> + repetition count based approach.
The cursor is set to the second-last line and a search for `^<C-r>2` is executed. If the second argument was 0, then
the search would succeed (Remember that VimComplete has wrapscan off), and the cursor would now be at the last line.
If the second argument was anything else, then the search would fail and the cursor would remain where it was before.
After the search is executed, the macro in the next word can be copied and executed.

Illustrative only — not literal keystroke content:

```
add(25, 30)

v cursor is here
25 30<C-a>
0 0
```

The position where the cursor is set before the search is executed works like the `else` clause, while the one which is matched
by regex becomes the `if` clause. More successive searches could be added to get `else if` clauses.

The limitation of the recursion with motion failure is solved with regex conditionals because the clauses could lead to no-op or
a self-call based on the base condition.
