package ch

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

type jsonFile struct {
	r         *bufio.Reader
	file      *os.File
	count     int
	startLine int
}

const MaxLineSize = 250000 * 1024

type JSONFileCursorConf struct {
	Line int    `desc:"start line in file"`
	File string `desc:"path to file"`
}

func NewJSONFileCursor(cfg *Conf) (Cursor, error) {
	file, err := os.Open(cfg.JFC.File)
	if err != nil {
		return nil, err
	}
	c := &jsonFile{
		file:      file,
		count:     0,
		r:         bufio.NewReaderSize(file, MaxLineSize),
		startLine: cfg.JFC.Line,
	}
	log.Println(fmt.Sprintf("go to %d line", c.startLine))
	for c.count < c.startLine-1 {
		c.count++
		_, isP, err := c.r.ReadLine()
		if isP {
			return nil, fmt.Errorf("too many tokens")
		}
		if err != nil {
			_ = c.file.Close()
			return nil, err
		}
	}
	return c, nil
}

func (c *jsonFile) Next() *bytes.Buffer {
	c.count++

	buff := bytes.NewBuffer(nil)

	b, isP, err := c.r.ReadLine()
	if isP {
		log.Println("too many tokens")
		return nil
	}
	if err != nil {
		log.Println(err)
		if err := c.file.Close(); err != nil {
			log.Println(err)
		}
		return nil
	}

	buff.Write(b)

	log.Println("current line in file =", c.count)
	return buff
}
