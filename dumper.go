package main

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
)

func dumper(ch chan *ES, cliEs *elastic.Client, cfg *conf) {
	log.Println("bulk dumper is up")
	bulk := cliEs.Bulk()
	for d := range ch {
		bulk = bulk.Add(elastic.NewBulkIndexRequest().Index(cfg.EsIndex).Doc(d))
		if bulk.NumberOfActions() >= cfg.EsBlkSz {
			log.Printf("dump new buffer with length = %d", cfg.EsBlkSz)
			if _, err := bulk.Do(context.Background()); err != nil {
				log.Fatal(err)
			}
			bulk.Reset()
		}
	}
	log.Printf("dump new buffer with length = %d", bulk.NumberOfActions())
	_, err := bulk.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	bulk.Reset()
}
