package ch

import (
	"ch2es/util"
)

type Conf struct {
	*util.HTTPConf

	URLParams struct {
		User string `param:"user" desc:"http param user"`
		Pass string `param:"password" desc:"http param password"`
	}

	DB              string `desc:"db name"`
	Table           string `desc:"table name"`
	QueryTimeoutSec int    `desc:"query timeout sec"`
	ConnTimeoutSec  int    `desc:"connection timeout sec"`
	CursorT         int    `desc:"cursor type"`
	Fields          string `desc:"fields"`
	Condition       string `desc:"condition"`
	DotReplacer     string `desc:"change dots symbol"`

	// offset cursor
	OFC *OffsetCursorConf `desc:"offset cursor config"`

	// timestamp cursor
	TSC *TSCursorConf `desc:"timestamp cursor config"`

	// json file cursor
	JFC *JSONFileCursorConf `desc:"json file cursor config"`

	// stdin cursor
	StdinC *StdInCursorConf `desc:"stdin cursor config"`
}

func (c *Conf) Print() {
	util.PrintFromDesc("[CLICKHOUSE HTTP CONFIG]:", *c.HTTPConf)

	switch cursorT(c.CursorT) {
	case offsetCursor:
		util.PrintFromDesc("[CLICKHOUSE OFFSET CURSOR CONFIG]:", *c.OFC)
	case timeStampCursor:
		util.PrintFromDesc("[CLICKHOUSE TIMESTAMP CURSOR CONFIG]:", *c.TSC)
	case jsonFileCursor:
		util.PrintFromDesc("[CLICKHOUSE JSON FILE CURSOR CONFIG]:", *c.JFC)
	case stdinCursor:
		util.PrintFromDesc("[CLICKHOUSE STDIN CURSOR CONFIG]:", *c.StdinC)
	}

	util.PrintFromDesc("[CLICKHOUSE CONFIG]:", *c)
	util.PrintFromDesc("[CLICKHOUSE CONFIG]:", c.URLParams)
}
