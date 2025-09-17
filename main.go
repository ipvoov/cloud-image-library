package main

import (
	_ "cloud/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"cloud/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
