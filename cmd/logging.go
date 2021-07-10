package cmd

import (
	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
)

var zlog *zap.Logger

func init() {
	logging.Register("github.com/hypha-dao/daoctl/cmd", &zlog)
}

// var zlog *zap.SugaredLogger

// func InitLogger() {

// 	config := zap.NewDevelopmentConfig()
// 	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
// 	config.EncoderConfig.TimeKey = "timestamp"
// 	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	logger, _ := config.Build()
// 	zap.ReplaceGlobals(logger)

// 	defer logger.Sync()
// 	zlog = logger.Sugar()
// }

// func getEncoder() zapcore.Encoder {
// 	encoderConfig := zap.NewProductionEncoderConfig()
// 	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
// 	return zapcore.NewConsoleEncoder(encoderConfig)
// }
