package ch

import (
	"bytes"
	"fmt"
	"log"
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

func NewOffsetCursor(cfg *Conf) Cursor {
	c := &offset{
		l:         cfg.OFC.Limit,
		o:         cfg.OFC.Offset,
		maxOffset: cfg.OFC.MaxOffset,
	}
	c.tempQBuff = bytes.NewBufferString(fmt.Sprintf("select %s from %s.%s", cfg.Fields, cfg.DB, cfg.Table))

	if cfg.Condition != "" {
		c.tempQBuff.WriteString(fmt.Sprintf(" where %s ", cfg.Condition))
	}

	if cfg.OFC.OrderField != "" {
		c.tempQBuff.WriteString(fmt.Sprintf(" order by %s ", cfg.OFC.OrderField))
	}

	c.tempQBuff.WriteString(fmt.Sprintf(" limit %d ", cfg.OFC.Limit))
	log.Printf("[CLICKHOUSE TEMPLATE QUERY]: %s  offset [n] format JSON", c.tempQBuff.String())

	return c
}

func (c *offset) Next() *bytes.Buffer {
	buff := bytes.NewBuffer(c.tempQBuff.Bytes())
	buff.WriteString(fmt.Sprintf(" offset %d format JSON", c.o))

	log.Println("current offset =", c.o)
	if c.o >= c.maxOffset {
		return nil
	}
	c.o += c.l

	return buff
}
