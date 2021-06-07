package es

import (
	"ch2es/common"
)

type Conf struct {
	*common.HTTPConf
	User  string `desc:"user"`
	Pass  string `desc:"password"`
	Index string `desc:"index"`
	BlkSz int    `desc:"bulk size"`
}

func (c *Conf) Print() {
	common.PrintFromDesc("[ELASTICSEARCH CONFIG]:", *c.HTTPConf)
	common.PrintFromDesc("[ELASTICSEARCH CONFIG]:", *c)
}
