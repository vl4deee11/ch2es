package cursor

import (
	"bufio"
	"bytes"
	"ch2es/log"
	"fmt"
	"io"
)

type ioTemp struct {
	r         *bufio.Reader
	rc        io.ReadCloser
	count     int
	startLine int
}

const MaxLineSize = 250000 * 1024

func newIOReaderTemp(startLine int, rc io.ReadCloser) (Cursor, error) {
	c := &ioTemp{
		rc:        rc,
		count:     0,
		r:         bufio.NewReaderSize(rc, MaxLineSize),
		startLine: startLine,
	}

	log.Info(fmt.Sprintf("go to %d line", c.startLine))
	for c.count < c.startLine-1 {
		c.count++
		_, isP, err := c.r.ReadLine()
		if isP {
			return nil, fmt.Errorf("too many tokens")
		}
		if err != nil {
			_ = c.rc.Close()
			return nil, err
		}
	}

	return c, nil
}

func (c *ioTemp) Next() *bytes.Buffer {
	c.count++

	buff := bytes.NewBuffer(nil)

	b, isP, err := c.r.ReadLine()
	if isP {
		log.Err(fmt.Errorf("too many tokens"))
		return nil
	}

	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			log.Info("data is end")
		} else {
			log.Err(err)
		}
		if err := c.rc.Close(); err != nil {
			log.Err(err)
		}
		return nil
	}

	buff.Write(b)

	log.Progress(fmt.Sprintf("current line = %d", c.count))
	return buff
}
