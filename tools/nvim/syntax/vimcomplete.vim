if exists("b:current_syntax")
  finish
endif

syntax match vimcompleteComment "^\s*\\.*$"

syntax match vimcompleteMacro "@[a-zA-Z0-9]"

syntax match vimcompleteMark "'[a-zA-Z]"

syntax match vimcompleteRegister '"[a-zA-Z0-9]'

syntax match vimcompleteSpecialKey "<[A-Za-z0-9-]\+>"

syntax match vimcompleteCount "\<[0-9]\+"

" Module imports, e.g. \r stdlib.add
syntax match vimcompleteModuleDirective "^\\r\>"
syntax match vimcompleteModuleName "\%(^\\r\s\+\)\@<=[A-Za-z0-9_.-]\+"

" Standard macro scaffolding

" A line that is *only* "o" -- opens the scratch line
syntax match vimcompleteScratchOpen "^o$"

" Leading "a" that starts an append-fragment line
syntax match vimcompleteMacroAppend "^\s*a"

" Fixed wrapper around the compile/yank line: $v0"X ydd
syntax match vimcompleteMacroCompile "^\$v0\ze\""
syntax match vimcompleteMacroCompile "y\%(\"_\)\?dd$"

highlight default link vimcompleteComment       Comment
highlight default link vimcompleteMacro         Function
highlight default link vimcompleteMark          Constant
highlight default link vimcompleteRegister      Special
highlight default link vimcompleteSpecialKey    Statement
highlight default link vimcompleteCount         Number
highlight default link vimcompleteScratchOpen   PreProc
highlight default link vimcompleteMacroAppend   PreProc
highlight default link vimcompleteMacroCompile  PreProc
highlight default link vimcompleteModuleDirective Include
highlight default link vimcompleteModuleName      Type

let b:current_syntax = "vimcomplete"
