# Conway's Game of Life in VimComplete
Active cells are represented by '#' and dead cells are represented by '·'.
Board is conditionally expanded if any any active cells are present on the padding layer.

State of a cell in the next generation based on number of live neighbours is computed using a mapping table which can be freely edited in the source code for an alternate game of life ruleset.

## building
```bash
vc build src -m src/main.vc -o gol.vcx
```

## running
```bash
vc run gol.vcx <board.txt> <generations>
```

For example:
```bash
vc run gol.vcx sample_input.txt 2
```

For a more dynamic visualization, remove the `ZQ` at the end of main.vc and run the built binary with the debugger
```bash
vc debug gol.vcx sample_input.txt 1
```
Press 'c' to skip the source code execution.

Then, to go one generation forward, use:
```
:lua VC.run({"@s"})
```
