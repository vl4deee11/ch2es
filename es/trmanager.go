package es

import (
	"context"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

type TrManager struct {
	cli   *elastic.Client
	index string
	blksz int
}

func (m *TrManager) Init(cfg *Conf) error {
	cliEs, err := elastic.NewClient(
		elastic.SetURL(m.httpF(cfg.Host, cfg.Port)),
	)
	if err != nil {
		return err
	}
	m.cli = cliEs
	m.index = cfg.Index
	m.blksz = cfg.BlkSz
	return nil
}

func (m *TrManager) httpF(h string, p int) string {
	return fmt.Sprintf("http://%s:%m", h, p)
}

func (m *TrManager) log(msg string) {
	log.Printf("[Elasticsearch Tranfer Manager]: %s", msg)
}

func (m *TrManager) logFatal(err error) {
	log.Fatalf("[Elasticsearch Tranfer Manager]: %s", err.Error())
}

func (m *TrManager) BulkDumper(ch chan map[string]interface{}) {
	m.log("bulk dumper is up")
	bulk := m.cli.Bulk()
	for d := range ch {
		bulk = bulk.Add(elastic.NewBulkIndexRequest().Index(m.index).Doc(d))
		if bulk.NumberOfActions() >= m.blksz {
			m.log(fmt.Sprintf("dump new buffer with length = %d", bulk.NumberOfActions()))
			if _, err := bulk.Do(context.Background()); err != nil {
				m.logFatal(err)
			}
			bulk.Reset()
		}
	}
	m.log(fmt.Sprintf("dump new buffer with length = %d", bulk.NumberOfActions()))
	_, err := bulk.Do(context.Background())
	if err != nil {
		m.logFatal(err)
	}
	bulk.Reset()
}
