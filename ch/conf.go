package ch

import (
	"ch2es/ch/cursor"
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
	CursorT         int    `desc:"Cursor type"`
	Fields          string `desc:"fields"`
	Condition       string `desc:"condition"`

	// offset Cursor
	OFC *cursor.OffsetCursorConf `desc:"offset cursor config"`

	// timestamp Cursor
	TSC *cursor.TimestampCursorConf `desc:"timestamp cursor config"`

	// json file Cursor
	JFC *cursor.JSONFileCursorConf `desc:"json file cursor config"`

	// stdin Cursor
	StdinC *cursor.StdInCursorConf `desc:"stdin cursor config"`
}

func (c *Conf) printDB() {
	util.PrintFromDesc("[CLICKHOUSE CONFIG]:", *c)
	util.PrintFromDesc("[CLICKHOUSE CONFIG]:", c.URLParams)
	util.PrintFromDesc("[CLICKHOUSE HTTP CONFIG]:", *c.HTTPConf)
}

func (c *Conf) Print() {
	switch cursor.T(c.CursorT) {
	case cursor.Offset:
		c.printDB()
		util.PrintFromDesc("[CLICKHOUSE OFFSET CURSOR CONFIG]:", *c.OFC)
	case cursor.Timestamp:
		c.printDB()
		util.PrintFromDesc("[CLICKHOUSE TIMESTAMP CURSOR CONFIG]:", *c.TSC)
	case cursor.JSONFile:
		util.PrintFromDesc("[CLICKHOUSE JSON FILE CURSOR CONFIG]:", *c.JFC)
	case cursor.Stdin:
		util.PrintFromDesc("[CLICKHOUSE STDIN CURSOR CONFIG]:", *c.StdinC)
	}

}
