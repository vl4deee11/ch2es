package es

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/olivere/elastic/v7"
)

type Writer struct {
	cli   *elastic.Client
	index string
	blksz int
}

func (w *Writer) Init(cfg *Conf) error {
	cliEs, err := elastic.NewClient(
		elastic.SetURL(w.httpF(cfg.Host, cfg.Port)),
	)
	if err != nil {
		return err
	}
	w.cli = cliEs
	w.index = cfg.Index
	w.blksz = cfg.BlkSz
	return nil
}

func (w *Writer) httpF(h string, p int) string {
	return fmt.Sprintf("http://%s:%d", h, p)
}

func (w *Writer) Write(ch chan map[string]interface{}, eCh chan error, wg *sync.WaitGroup) {
	log.Println("start new elasticsearch writer")
	defer func() {
		wg.Done()
		log.Println("writer is stop")
	}()

	bulk := w.cli.Bulk()
	for d := range ch {
		bulk = bulk.Add(elastic.NewBulkIndexRequest().Index(w.index).Doc(d))
		if bulk.NumberOfActions() >= w.blksz {
			log.Println("dump new buffer with length =", bulk.NumberOfActions())
			if _, err := bulk.Do(context.Background()); err != nil {
				eCh <- err
				break
			}
			bulk.Reset()
		}
	}
	log.Println("chan is close, dump new buffer with length =", bulk.NumberOfActions())
	_, err := bulk.Do(context.Background())
	if err != nil {
		log.Println(err)
	}
}