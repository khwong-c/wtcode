package log

import (
	"github.com/inconshreveable/log15"
	"github.com/mattn/go-colorable"
)

func init() {
	log15.Root().SetHandler(
		log15.StreamHandler(colorable.NewColorableStdout(), log15.TerminalFormat()),
	)
}

func NewLogger(pkg string) log15.Logger {
	return log15.Root().New("package", pkg)
}

func SetLogLevel(maxLvl log15.Lvl) {
	log15.Root().SetHandler(
		log15.LvlFilterHandler(maxLvl, log15.Root().GetHandler()),
	)
}
