package ch

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
)

type ioReaderTemp struct {
	r         *bufio.Reader
	source    io.ReadCloser
	count     int
	startLine int
}

const MaxLineSize = 250000 * 1024

func newIOReaderTempCursor(startLine int, source io.ReadCloser) (Cursor, error) {
	c := &ioReaderTemp{
		source:    source,
		count:     0,
		r:         bufio.NewReaderSize(source, MaxLineSize),
		startLine: startLine,
	}

	log.Println(fmt.Sprintf("go to %d line", c.startLine))
	for c.count < c.startLine-1 {
		c.count++
		_, isP, err := c.r.ReadLine()
		if isP {
			return nil, fmt.Errorf("too many tokens")
		}
		if err != nil {
			_ = c.source.Close()
			return nil, err
		}
	}

	return c, nil
}

func (c *ioReaderTemp) Next() *bytes.Buffer {
	c.count++

	buff := bytes.NewBuffer(nil)

	b, isP, err := c.r.ReadLine()
	if isP {
		log.Println("too many tokens")
		return nil
	}

	if err != nil {
		log.Println(err)
		if err := c.source.Close(); err != nil {
			log.Println(err)
		}
		return nil
	}

	buff.Write(b)

	log.Println("current line in file =", c.count)
	return buff
}
