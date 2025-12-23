package main

import (
	_ "tele-disk/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"tele-disk/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
