package es

import (
	"ch2es/es/converter"
	"ch2es/log"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

type Writer struct {
	cli         *elastic.Client
	index       string
	idField     string
	bulksz      int
	qTimeout    time.Duration
	docConv     converter.Converter
	jsonCleaner func(string) string
}

func NewWriter(cfg *Conf) (*Writer, error) {
	w := &Writer{
		qTimeout: time.Duration(cfg.QueryTimeoutSec) * time.Second,
		index:    cfg.Index,
		idField:  cfg.IDField,
		bulksz:   cfg.BulkSz,
	}
	cfg.BuildHTTP()

	cliEs, err := elastic.NewClient(
		elastic.SetURL(cfg.URL),
		elastic.SetBasicAuth(cfg.User, cfg.Pass),
	)
	if err != nil {
		return nil, err
	}
	w.cli = cliEs

	switch converter.T(cfg.ConverterT) {
	case converter.Nested:
		w.docConv = converter.NewNested(cfg.NCC)
	default:
		w.docConv = converter.NewNull()
	}

	w.jsonCleaner = w.getJSONCleaner(cfg.DotReplacer)
	return w, nil
}

func (w *Writer) Write(ch chan map[string]interface{}, eCh chan error, wg *sync.WaitGroup) {
	log.Info("start new elasticsearch writer")
	defer func() {
		wg.Done()
		log.Info("writer is stop")
	}()

	bulk := w.cli.Bulk()
	for d := range ch {
		d = w.docConv.Convert(d, w.jsonCleaner)
		bReq := elastic.NewBulkIndexRequest()
		if w.idField != "" {
			bReq = bReq.OpType("index").Id(fmt.Sprint(d[w.idField]))
		}
		bulk = bulk.Add(bReq.Index(w.index).Doc(d))
		if bulk.NumberOfActions() >= w.bulksz {
			log.Info(fmt.Sprintf("dump new buffer with length [%d]", bulk.NumberOfActions()))
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
	log.Info(fmt.Sprintf("chan is close, dump new buffer with length [%d]", bulk.NumberOfActions()))
	ctx, cancel := context.WithTimeout(context.Background(), w.qTimeout)
	r, err := bulk.Do(ctx)
	if err != nil {
		log.Err(err)
	}
	if r != nil && r.Errors {
		log.Err(fmt.Errorf("error happened in bulk request"))
	}
	cancel()
}

func (w *Writer) getJSONCleaner(r string) func(k string) string {
	if r != "" {
		return func(k string) string {
			return strings.ReplaceAll(k, ".", r)
		}
	}

	return func(k string) string {
		return k
	}
}

//nolint:unused //set TODO
// TODO: parse dots from fields name and generate sub-maps
func (w *Writer) rebuildJSON(d map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for k := range d {
		if subNames := strings.Split(k, "."); len(subNames) > 1 {
			fmt.Println(subNames)
			cm, ok := m[subNames[0]]
			if !ok {
				m[subNames[0]] = make(map[string]interface{})
				cm = m[subNames[0]]
			}
			tm := cm.(map[string]interface{})
			for i := 1; i < len(subNames); i++ {
				_, ok = tm[subNames[i]]
				if !ok {
					tm[subNames[i]] = make(map[string]interface{})
					tm = tm[subNames[i]].(map[string]interface{})
				}
			}

			tm[subNames[len(subNames)-1]] = d[k]
		} else {
			m[subNames[0]] = d[k]
		}
	}
	return m
}
