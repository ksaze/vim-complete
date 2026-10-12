vim.api.nvim_create_autocmd("VimLeave", {
  callback = function()
    local path = vim.env.VC_RESULT
    if not path then
      return
    end

    local result = vim.json.encode({
      exit_code = tonumber(vim.fn.getreg("x")) or 0,
      stdout = vim.fn.getreg("y"),
      stderr = vim.fn.getreg("z"),
    })
    vim.fn.writefile({ result }, path)
  end,
})

vim.o.swapfile = false
vim.opt.autoindent = false
vim.opt.wrapscan = false

VC = {}

-- Shared prologue for both VC.run implementations.
function VC.init()
  vim.cmd("bnext")
  vim.cmd("bprevious")
  vim.fn.setreg("c", table.concat(vim.v.argv, " ", 3))
  vim.fn.setreg("x", "0")
  vim.fn.setreg("y", "")
  vim.fn.setreg("z", "")
end

-- Errors the interpreter tolerates (failed searches etc.).
VC.ignored = { E486 = true, E384 = true, E385 = true }

-- Runs one instruction; rethrows unless the error is ignorable.
function VC.exec(cmd)
  local ok, err = pcall(vim.cmd, "normal! " .. cmd)
  if not ok then
    local code = tostring(err):match("(E%d+)")
    if not VC.ignored[code] then
      vim.fn.setreg("x", "1")
      vim.fn.setreg("z", tostring(err))
      vim.cmd("qall!")
    end
  end
end

function VC.unresolved(module)
  -- Report through the registers so the message reaches vc via the result
  -- file (and doesn't scribble over the TUI in debug mode).
  vim.fn.setreg("x", "1")
  vim.fn.setreg(
    "z",
    string.format(
      "VC Link Error\n\nUnresolved dependency: %s\n\n"
        .. "This is an intermediate object file that must be linked before execution.\n",
      module
    )
  )
  vim.cmd("qall!")
end
