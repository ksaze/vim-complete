-- Activate with, from the project root:
--   :source nvim/vimcomplete.lua

local script_path = debug.getinfo(1, "S").source:sub(2) -- strip leading '@'
local plugin_dir = vim.fn.fnamemodify(script_path, ":h") -- .../nvim

-- Make ./syntax and ./ftplugin (next to this file) discoverable
vim.opt.runtimepath:prepend(plugin_dir)

-- Register the filetype for *.vc files
vim.filetype.add({
  extension = {
    vc = "vimcomplete",
  },
})

local function apply_to_current_buffer()
  if vim.bo.filetype ~= "vimcomplete" then
    vim.bo.filetype = "vimcomplete"
  end
  vim.cmd("syntax on")
end

-- Apply immediately in case a .vm file is already open
apply_to_current_buffer()

-- Keep applying to any .vc files opened later in this session
vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
  pattern = "*.vc",
  callback = apply_to_current_buffer,
})

-- Manual re-trigger, e.g. if highlighting ever looks stale
vim.api.nvim_create_user_command("VmSyntax", apply_to_current_buffer, {})

print("Vim Complete (.vm) syntax support loaded -- :VmSyntax to reapply")
