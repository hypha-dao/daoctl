package main

import (
	"github.com/dfuse-io/logging"
	"github.com/hypha-dao/daoctl/cmd"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

//lint:ignore U1000 leveraged at runtime
var tracer = logging.ApplicationLogger("daoctl", "daoctl", &zlog)

func main() {
	cmd.Execute()
}
