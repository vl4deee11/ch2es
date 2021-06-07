package es

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

type Writer struct {
	cli      *elastic.Client
	index    string
	blksz    int
	qTimeout time.Duration
}

func (w *Writer) Init(cfg *Conf) error {
	cfg.BuildHTTP()

	cliEs, err := elastic.NewClient(
		elastic.SetURL(cfg.URL),
		elastic.SetBasicAuth(cfg.User, cfg.Pass),
	)
	if err != nil {
		return err
	}
	w.qTimeout = time.Duration(cfg.QueryTimeoutSec) * time.Second
	w.cli = cliEs
	w.index = cfg.Index
	w.blksz = cfg.BlkSz
	return nil
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
			ctx, cancel := context.WithTimeout(context.Background(), w.qTimeout)
			r, err := bulk.Do(ctx)
			if err != nil {
				eCh <- err
				cancel()
				break
			}
			if r.Errors {
				eCh <- fmt.Errorf("error happened in bulk request")
				cancel()
				break
			}
			cancel()
			bulk.Reset()
		}
	}
	log.Println("chan is close, dump new buffer with length =", bulk.NumberOfActions())
	ctx, cancel := context.WithTimeout(context.Background(), w.qTimeout)
	_, err := bulk.Do(ctx)
	if err != nil {
		log.Println(err)
	}
	cancel()
}
