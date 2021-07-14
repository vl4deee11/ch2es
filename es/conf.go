package es

import (
	"ch2es/es/converter"
	"ch2es/util"
)

type Conf struct {
	*util.HTTPConf
	User            string `desc:"user"`
	Pass            string `desc:"password"`
	Index           string `desc:"index"`
	BulkSz          int    `desc:"bulk size"`
	IDField         string `desc:"id field"`
	QueryTimeoutSec int    `desc:"query timeout sec"`
	ConverterT      int    `desc:"converter type"`
	DotReplacer     string `desc:"change dots symbol"`

	NCC *converter.NestedConverterConf `desc:"nested converter config"`
}

func (c *Conf) Print() {
	if converter.T(c.ConverterT) == converter.Nested {
		util.PrintFromDesc("[ELASTICSEARCH NESTED CONVERTER CONFIG]:", *c.NCC)
	}

	util.PrintFromDesc("[ELASTICSEARCH CONFIG]:", *c.HTTPConf)
	util.PrintFromDesc("[ELASTICSEARCH CONFIG]:", *c)
}
