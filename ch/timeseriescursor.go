package ch

import (
	"bytes"
	"fmt"
	"log"
)

type timestamp struct {
	min       int
	curr      int
	max       int
	step      int
	field     string
	tempQBuff *bytes.Buffer
}

type TSCursorConf struct {
	StepSec int    `desc:"step in seconds"`
	Field   string `desc:"field"`
	Min     int    `desc:"min unix timestamp"`
	Max     int    `desc:"max unix timestamp"`
}

func NewTimestampCursor(cfg *Conf) Cursor {
	c := &timestamp{
		min:   cfg.TSC.Min,
		curr:  cfg.TSC.Min + cfg.TSC.StepSec,
		max:   cfg.TSC.Max,
		step:  cfg.TSC.StepSec,
		field: cfg.TSC.Field,
	}
	c.tempQBuff = bytes.NewBufferString(fmt.Sprintf("select %s from %s.%s", cfg.Fields, cfg.DB, cfg.Table))

	if cfg.Condition != "" {
		c.tempQBuff.WriteString(fmt.Sprintf(" where %s and ", cfg.Condition))
	} else {
		c.tempQBuff.WriteString(" where ")
	}

	log.Printf(
		"[CLICKHOUSE TEMPLATE QUERY]: %s (toUnixTimestamp(%s) >= [n] and toUnixTimestamp(%s) < [n+step])  format JSON",
		c.tempQBuff.String(),
		c.field,
		c.field,
	)
	return c
}

func (c *timestamp) Next() *bytes.Buffer {
	if c.min > c.max {
		return nil
	}
	if c.curr > c.max {
		c.curr = c.max
	}

	buff := bytes.NewBuffer(c.tempQBuff.Bytes())
	buff.WriteString(
		fmt.Sprintf("(toUnixTimestamp(%s) >= %d and toUnixTimestamp(%s) < %d)  format JSON",
			c.field,
			c.min,
			c.field,
			c.curr,
		),
	)

	log.Println(fmt.Sprintf("current time range = start: %d, end: %d", c.min, c.curr))
	c.min += c.step
	c.curr += c.step

	return buff
}
