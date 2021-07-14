package cursor

import (
	"os"
)

type JSONFileCursorConf struct {
	Line int    `desc:"start line in file"`
	File string `desc:"path to file"`
}

func NewJSONFile(cfg *JSONFileCursorConf) (Cursor, error) {
	file, err := os.Open(cfg.File)
	if err != nil {
		return nil, err
	}
	return newIOReaderTemp(cfg.Line, file)
}
