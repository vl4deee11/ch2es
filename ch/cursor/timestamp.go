package cursor

import (
	"bytes"
	"ch2es/log"
	"fmt"
)

type timestamp struct {
	min       int
	curr      int
	max       int
	step      int
	field     string
	tempQBuff *bytes.Buffer
}

type TimestampCursorConf struct {
	StepSec int    `desc:"step in seconds"`
	Field   string `desc:"field"`
	Min     int    `desc:"min unix timestamp"`
	Max     int    `desc:"max unix timestamp"`
}

func NewTimestamp(cfg *TimestampCursorConf, fields, db, table, condition string) Cursor {
	c := &timestamp{
		min:       cfg.Min,
		curr:      cfg.Min + cfg.StepSec,
		max:       cfg.Max,
		step:      cfg.StepSec,
		field:     cfg.Field,
		tempQBuff: bytes.NewBufferString(fmt.Sprintf("select %s from %s.%s", fields, db, table)),
	}

	if condition != "" {
		c.tempQBuff.WriteString(fmt.Sprintf(" where %s and ", condition))
	} else {
		c.tempQBuff.WriteString(" where ")
	}

	log.Info(
		fmt.Sprintf(
			"[CLICKHOUSE TEMPLATE QUERY]: %s (toUnixTimestamp(%s) >= [n] and toUnixTimestamp(%s) < [n+step])  format JSON",
			c.tempQBuff.String(),
			c.field,
			c.field,
		),
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
		fmt.Sprintf(
			"(toUnixTimestamp(%s) >= %d and toUnixTimestamp(%s) < %d)  format JSON",
			c.field,
			c.min,
			c.field,
			c.curr,
		),
	)

	log.Progress(fmt.Sprintf("current time range = start: %d, end: %d", c.min, c.curr))
	c.min += c.step
	c.curr += c.step

	return buff
}
