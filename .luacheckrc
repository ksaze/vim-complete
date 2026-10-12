std = "luajit"
cache = true
codes = true

exclude_files = {
  ".luarocks/",
}

ignore = {
  "111", -- Setting an undefined global variable.
  "211", -- Unused local variable.
  "411", -- Redefining a local variable.
  "121", -- Setting a read-only global variable.
  "122", -- Setting a read-only field.
}

read_globals = {
  "vim",
  "VC",
  "describe",
  "it",
  "bit",
  "assert",
  "before_each",
  "after_each",
}
