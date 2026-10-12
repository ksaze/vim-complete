dofile(
  vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":p:h")
    .. "/common.lua"
)

function VC.run(program)
  VC.init()
  local step = true

  for i, cmd in ipairs(program) do
    vim.api.nvim_echo(
      { { string.format("[%3d] %s", i, cmd), "Question" } },
      false,
      {}
    )
    VC.exec(cmd)
    vim.cmd("redraw")

    if step then
      local key = vim.fn.nr2char(vim.fn.getchar())
      if key == "q" then
        print("Execution aborted.")
        return
      elseif key == "c" then
        step = false
      end
    end
  end

  print("Execution finished.")
end
