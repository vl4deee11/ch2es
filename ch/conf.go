package ch

import (
	"ch2es/common"
)

type Conf struct {
	*common.HTTPConf
	Fields    string `desc:"fields"`
	Condition string `desc:"condition"`
	URLParams struct {
		User string `param:"user" desc:"http param user"`
		Pass string `param:"password" desc:"http param password"`
	}

	DB              string `desc:"db name"`
	OrderField      string `desc:"order field"`
	Table           string `desc:"table name"`
	Limit           int    `desc:"limit"`
	QueryTimeoutSec int    `desc:"query timeout sec"`
	ConnTimeoutSec  int    `desc:"connection timeout sec"`
}

func (c *Conf) Print() {
	common.PrintFromDesc("[CLICKHOUSE CONFIG]:", *c.HTTPConf)
	common.PrintFromDesc("[CLICKHOUSE CONFIG]:", *c)
	common.PrintFromDesc("[CLICKHOUSE CONFIG]:", c.URLParams)
}
