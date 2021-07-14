package cursor

import (
	"bytes"
	"ch2es/log"
	"fmt"
)

type offset struct {
	o         int
	l         int
	maxOffset int
	tempQBuff *bytes.Buffer
}

type OffsetCursorConf struct {
	MaxOffset  int    `desc:"max offset"`
	OrderField string `desc:"order field"`
	Limit      int    `desc:"limit"`
	Offset     int    `desc:"offset"`
}

func NewOffset(cfg *OffsetCursorConf, fields, db, table, condition string) Cursor {
	c := &offset{
		l:         cfg.Limit,
		o:         cfg.Offset,
		maxOffset: cfg.MaxOffset,
		tempQBuff: bytes.NewBufferString(fmt.Sprintf("select %s from %s.%s", fields, db, table)),
	}

	if condition != "" {
		c.tempQBuff.WriteString(fmt.Sprintf(" where %s ", condition))
	}

	if cfg.OrderField != "" {
		c.tempQBuff.WriteString(fmt.Sprintf(" order by %s ", cfg.OrderField))
	}

	c.tempQBuff.WriteString(fmt.Sprintf(" limit %d ", cfg.Limit))
	log.Info(fmt.Sprintf("[CLICKHOUSE TEMPLATE QUERY]: %s  offset [n] format JSON", c.tempQBuff.String()))

	return c
}

func (c *offset) Next() *bytes.Buffer {
	buff := bytes.NewBuffer(c.tempQBuff.Bytes())
	buff.WriteString(fmt.Sprintf(" offset %d format JSON", c.o))

	log.Progress(fmt.Sprintf("current offset = %d", c.o))
	if c.o >= c.maxOffset {
		return nil
	}
	c.o += c.l

	return buff
}
