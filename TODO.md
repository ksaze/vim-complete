# TODO

Suggested areas to contribute. These are just soft suggestions—feel free to work on anything else.


- [compiler] Implement a segment stage which segments source code such that each line 
to just one mode (normal, visual, insert) or just one logical motion unit. 
This can make analysis and optimisation easier; however splitting a line this isn't 
necessarily equivalent if the line contains a motion like 'f' which might fail. 
Thus line index information must be stored and used during optimisation phase. 
After which the line index information can be used to combine lines which can 
be without chaning program semantics.

- [compiler] Suggest and document possible optimisations in `optimise.vc`. Post in 
discussions for review before creating a pull request. Review unimplemented optmisations
in the same file and implement them.

- [linker] Move the current linker to bootstrap stage and work on a cleaner linker which
uses stdlib and is split into multiple modules which can be unit tested independently. A 
potential addition for a mode which doesn't link stdlib files could be added for the cli
library build mode.

- [tests/stdlib] `search_and_replace`, `search_and_transform`, and `file_assert_diff` 
currently don't have any unit tests. Existing unit tests could be also be improved.

- [stdlib] A possible macro for simplifying unit test implementations

- [stdlib] Script for performing operations on a scratch buffer which could be used 
like python context managers.
```
\r scratch
  ...
  ...
ZQ or \r scratch_close
```

- [examples] Any other interesting programs or improve the current game of life
program by adding support for a config file which allows changing ruleset

- [meta-programming] Introduce meta-programming VimComplete standards which can
make programs more readable

- [tools] docs and fmt utilty for vc programs. `\\#` could be used as a prefix for
docs.

- [cli] The go CLI was just a chore for me to complete so that I could return to
working on VimComplete programs. Both the overall design and code quality are terrible, 
so they could use some rework.
