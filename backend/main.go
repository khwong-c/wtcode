package main

import (
	"os"

	"github.com/khwong-c/wtcode/server"
	"github.com/khwong-c/wtcode/tooling/di"
	"github.com/khwong-c/wtcode/tooling/log"
)

var logger = log.NewLogger("main")

func main() {
	injector := di.CreateInjector(true, false)
	svr := di.InvokeOrProvide(injector, server.CreateServer)
	logger.Info("Server created", "addr", svr.Addr)

	svr.Serve()
	err := injector.ShutdownOnSignals(os.Interrupt, os.Kill)
	if err != nil {
		logger.Error("Injector shutdown error", "err", err)
	}
}
