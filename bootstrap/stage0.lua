-- Seed compiler for vim complete
-- Removes comments and empty lines
-- Encloses the remaining lines in a VC.run() invocation provided by the runtime

local argv = vim.fn.argv()

-- argv[1] = input .vc
-- argv[2] = output .vcb
local input = argv[1]
local output = argv[2]

if not input or not output then
  error(
    "Usage: nvim --headless -u NONE -S bootstrap/stage0.lua -- input.vc output.lua"
  )
end

-- Custom %q replacement: leaves control bytes as literal bytes,
-- except CR (decimal 13 / 0x0D, ^M), which is kept as an escape sequence
-- so it can't be literally misread as \r.
local function quote(s)
  local out = { '"' }
  local n = #s
  for i = 1, n do
    local c = s:sub(i, i)
    local b = s:byte(i)
    if c == '"' then
      out[#out + 1] = '\\"'
    elseif c == "\\" then
      out[#out + 1] = "\\\\"
    elseif b == 13 then
      local nb = s:byte(i + 1)
      if nb and nb >= 48 and nb <= 57 then
        out[#out + 1] = "\\013" -- pad so next char isn't read as part of the escape
      else
        out[#out + 1] = "\\13"
      end
    else
      out[#out + 1] = c -- literal byte, control chars included
    end
  end
  out[#out + 1] = '"'
  return table.concat(out)
end

local f = assert(io.open(input, "rb"))
local data = f:read("*a")
f:close()

local program = {}

local f = assert(io.open(input, "rb"))
local data = f:read("*a")
f:close()

for line in (data .. "\n"):gmatch("(.-)\n") do
  -- Trim only spaces/tabs for comment detection.
  local trimmed = line:gsub("^[ \t]+", "")

  -- Ignore comments and blank lines.
  if trimmed ~= "" and not trimmed:match("^\\\\") then
    table.insert(program, trimmed)
  end
end

local out = {}

table.insert(
  out,
  [[
lua << EOF
VC.run({]]
)

for _, line in ipairs(program) do
  -- table.insert(out, string.format("%q", line) .. ",")
  table.insert(out, quote(line) .. ",")
end

table.insert(
  out,
  [[
})
EOF
]]
)

local f = assert(io.open(output, "wb"))
f:write(table.concat(out, "\n"))
f:close()

print(("Compiled %s -> %s\n"):format(input, output))
vim.fn.setreg("z", "")

vim.cmd("qall!")
