package cursor

import "os"

type StdInCursorConf struct {
	Line int `desc:"start line in file"`
}

func NewStdin(cfg *StdInCursorConf) (Cursor, error) {
	return newIOReaderTemp(cfg.Line, os.Stdin)
}
