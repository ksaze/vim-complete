dofile(
  vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":p:h")
    .. "/common.lua"
)

function VC.run(program)
  VC.init()
  for _, cmd in ipairs(program) do
    VC.exec(cmd)
  end
end
