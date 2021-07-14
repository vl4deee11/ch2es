package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger = nil

type adapterZapT uint8

const (
	infoL adapterZapT = iota
	errorL
	debugL
)

func lvlToZapLvl(lvl int) zapcore.Level {
	switch adapterZapT(lvl) {
	case debugL:
		return zapcore.DebugLevel
	case errorL:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func Err(err error) {
	log.Error(err.Error())
}

func Info(msg string) {
	log.Info(msg)
}

func Progress(msg string) {
	log.Debug(fmt.Sprintf("PROGRESS|%s", msg))
}

func Init(lvl int) error {
	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(lvlToZapLvl(lvl)),
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel:      zapcore.CapitalColorLevelEncoder,
			TimeKey:          "timestamp",
			EncodeTime:       zapcore.RFC3339TimeEncoder,
			ConsoleSeparator: "|",
			LevelKey:         "lvl",
			MessageKey:       "msg",
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}

	log = l
	return nil
}
