package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	cfg := new(conf)
	cfg.parse()

	reader, err := cfg.getReader()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ch2es start")

	writer, err := cfg.getWriter()
	if err != nil {
		log.Fatal(err)
	}

	wCh := make(chan map[string]interface{})
	defer close(wCh)

	rCh := make(chan string)
	defer close(rCh)

	for i := 0; i < cfg.ThreadsNum; i++ {
		go writer.Write(wCh)
	}

	for i := 0; i < cfg.ThreadsNum; i++ {
		go reader.Read(rCh, wCh)
	}

	offset := reader.ReadInitialOffset()

	log.Println("start offset =", offset)
	for offset < cfg.MaxOffset {
		rCh <- fmt.Sprintf("offset %d format JSON", offset)
		offset += cfg.ChConf.Limit
		if err := ioutil.WriteFile("stats", []byte(fmt.Sprintf("%d", offset)), 0600); err != nil {
			log.Fatal(err)
		}
		log.Println("current offset =", offset)
	}
	log.Println("successfully transferred")
}
