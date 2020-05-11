package cmd

import (
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

// SetLogger ...
func SetLogger(l *zap.Logger) {
	zlog = l
}
