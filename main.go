package main

import (
	_ "cloud/internal/logic"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"cloud/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
