package main

import (
	"github.com/khwong-c/wtcode/server"
	"github.com/khwong-c/wtcode/tooling/di"
	"github.com/khwong-c/wtcode/tooling/log"
)

var logger = log.NewLogger("main")

func main() {
	injector := di.CreateInjector(true, false)
	svr := di.InvokeOrProvide(injector, server.CreateServer)
	logger.Info("Server created", "addr", svr.Addr, "cfg", svr.Config)
	_ = injector.Shutdown()
}
