package ch

import "os"

type StdInCursorConf struct {
	Line int `desc:"start line in file"`
}

func NewStdInCursor(cfg *StdInCursorConf) (Cursor, error) {
	return newIOReaderTempCursor(cfg.Line, os.Stdin)
}
