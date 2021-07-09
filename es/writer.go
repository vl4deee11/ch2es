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
	idField  string
	blksz    int
	qTimeout time.Duration
}

func NewWriter(cfg *Conf) (*Writer, error) {
	w := new(Writer)
	cfg.BuildHTTP()

	cliEs, err := elastic.NewClient(
		elastic.SetURL(cfg.URL),
		elastic.SetBasicAuth(cfg.User, cfg.Pass),
	)
	if err != nil {
		return nil, err
	}
	w.qTimeout = time.Duration(cfg.QueryTimeoutSec) * time.Second
	w.cli = cliEs
	w.index = cfg.Index
	w.blksz = cfg.BlkSz
	w.idField = cfg.IDField
	return w, nil
}

func (w *Writer) Write(ch chan map[string]interface{}, eCh chan error, wg *sync.WaitGroup) {
	log.Println("start new elasticsearch writer")
	defer func() {
		wg.Done()
		log.Println("writer is stop")
	}()

	bulk := w.cli.Bulk()
	for d := range ch {
		bReq := elastic.NewBulkIndexRequest()
		if w.idField != "" {
			bReq = bReq.OpType("index").Id(fmt.Sprint(d[w.idField]))
		}
		bulk = bulk.Add(bReq.Index(w.index).Doc(d))
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
	r, err := bulk.Do(ctx)
	if err != nil {
		log.Println(err)
	}
	if r != nil && r.Errors {
		log.Println("error happened in bulk request")
	}
	cancel()
}
