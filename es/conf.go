package es

import (
	"ch2es/util"
)

type Conf struct {
	*util.HTTPConf
	User            string `desc:"user"`
	Pass            string `desc:"password"`
	Index           string `desc:"index"`
	BlkSz           int    `desc:"bulk size"`
	IDField         string `desc:"id field"`
	QueryTimeoutSec int    `desc:"query timeout sec"`
}

func (c *Conf) Print() {
	util.PrintFromDesc("[ELASTICSEARCH CONFIG]:", *c.HTTPConf)
	util.PrintFromDesc("[ELASTICSEARCH CONFIG]:", *c)
}
