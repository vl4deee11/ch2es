package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
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
	rCh := make(chan string)
	eCh := make(chan error)

	for i := 0; i < cfg.ThreadsNum; i++ {
		go writer.Write(wCh)
	}

	for i := 0; i < cfg.ThreadsNum; i++ {
		go reader.Read(rCh, wCh, eCh)
	}

	var wg *sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case err := <- eCh:
			close(rCh)
			close(wCh)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}
	}()


	offset := reader.ReadInitialOffset()

	log.Println("start offset =", offset)
	for offset < cfg.MaxOffset {
		rCh <- fmt.Sprintf("kffset %d format JSON", offset)
		offset += cfg.ChConf.Limit
		if err := ioutil.WriteFile("stats", []byte(fmt.Sprintf("%d", offset)), 0600); err != nil {
			eCh <- err
		}
		log.Println("current offset =", offset)
	}
	log.Println("successfully transferred")
}
