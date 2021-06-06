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

	log.Println("=========== START ===========")

	writer, err := cfg.getWriter()
	if err != nil {
		log.Fatal(err)
	}

	wCh := make(chan map[string]interface{})
	rCh := make(chan string)

	//nolint:gomnd // 2 use for non blocking write for writer and reader
	eCh := make(chan error, 2*cfg.ThreadsNum)

	var wwg sync.WaitGroup
	for i := 0; i < cfg.ThreadsNum; i++ {
		wwg.Add(1)
		go writer.Write(wCh, eCh, &wwg)
	}

	var rwg sync.WaitGroup
	for i := 0; i < cfg.ThreadsNum; i++ {
		rwg.Add(1)
		go reader.Read(rCh, wCh, eCh, &rwg)
	}

	defer func() {
		end(&rwg, &wwg, wCh, rCh, eCh)
	}()

	offset := reader.ReadInitialOffset()

	log.Println("start offset =", offset)
	if offset >= cfg.MaxOffset {
		log.Println("incorrect config offset >= maxOffset, check stats file")
		return
	}

	for {
		select {
		case err := <-eCh:
			log.Println(err)
			return
		default:
			rCh <- fmt.Sprintf("offset %d format JSON", offset)
			offset += cfg.ChConf.Limit
			if err := ioutil.WriteFile("stats", []byte(fmt.Sprintf("%d", offset)), 0600); err != nil {
				log.Println(err)
				return
			}
			log.Println("current offset =", offset)
			if offset >= cfg.MaxOffset {
				return
			}
		}
	}
}

func end(
	rwg *sync.WaitGroup,
	wwg *sync.WaitGroup,
	wCh chan map[string]interface{},
	rCh chan string,
	eCh chan error,
) {
	close(rCh)
	rwg.Wait()
	close(wCh)
	wwg.Wait()
	close(eCh)
	log.Println("=========== END ===========")
}
