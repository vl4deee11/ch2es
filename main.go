package main

import (
	"bytes"
	"ch2es/ch"
	"ch2es/es"
	"log"
	"sync"
	"time"
)

func main() {
	cfg := new(conf)
	cfg.parse()

	log.Println("=========== START ===========")

	reader, cursor, err := ch.NewReader(cfg.ChConf)
	if err != nil {
		log.Fatal("FATAL ERROR = ", err)
	}

	writer, err := es.NewWriter(cfg.EsConf)
	if err != nil {
		log.Fatal("FATAL ERROR = ", err)
	}

	wCh := make(chan map[string]interface{})
	rCh := make(chan *bytes.Buffer)

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

	start := time.Now()
	defer func() {
		end(&rwg, &wwg, wCh, rCh, eCh)
		log.Printf("Elapsed: [%s]", time.Since(start))
	}()

	for {
		select {
		case err := <-eCh:
			log.Println("FATAL ERROR = ", err)
			return
		default:
			b := cursor.Next()
			if b == nil {
				log.Println("stop, cursor is end")
				return
			}
			rCh <- b
		}
	}
}

func end(
	rwg, wwg *sync.WaitGroup,
	wCh chan map[string]interface{},
	rCh chan *bytes.Buffer,
	eCh chan error,
) {
	close(rCh)
	rwg.Wait()
	close(wCh)
	wwg.Wait()
	close(eCh)
	log.Println("=========== END ===========")
}
