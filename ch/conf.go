package ch

import (
	"ch2es/common"
)

type Conf struct {
	*common.HTTPConf

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
	OFC *OffsetCursorConf `desc:"Offset cursor config"`

	// timestamp cursor
	TSC *TSCursorConf `desc:"Timestamp cursor config"`

	// json file cursor
	JFC *JSONFileCursorConf `desc:"JSON file cursor config"`
}

func (c *Conf) Print() {
	common.PrintFromDesc("[CLICKHOUSE HTTP CONFIG]:", *c.HTTPConf)

	switch cursorT(c.CursorT) {
	case offsetCursor:
		common.PrintFromDesc("[CLICKHOUSE OFFSET CURSOR CONFIG]:", *c.OFC)
	case timeStampCursor:
		common.PrintFromDesc("[CLICKHOUSE TIMESTAMP CURSOR CONFIG]:", *c.TSC)
	case fileCursor:
		common.PrintFromDesc("[CLICKHOUSE FILE CURSOR CONFIG]:", *c.JFC)
	}

	common.PrintFromDesc("[CLICKHOUSE CONFIG]:", *c)
	common.PrintFromDesc("[CLICKHOUSE CONFIG]:", c.URLParams)
}
