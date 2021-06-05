package es

import (
	"context"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

type HTTPD struct {
	cli *elastic.Client
}

func (d *HTTPD) SetHTTP(h string, p int) error {
	cliEs, err := elastic.NewClient(
		elastic.SetURL(d.httpF(h, p)),
	)
	if err != nil {
		return err
	}
	d.cli = cliEs
	return nil
}

func (d *HTTPD) httpF(h string, p int) string {
	return fmt.Sprintf("http://%s:%d", h, p)
}

func (d *HTTPD) log(msg string) {
	log.Printf("[Elasticsearch HTTP DRIVER]: %s", msg)
}

func (d *HTTPD) logFatal(err error) {
	log.Fatalf("[Elasticsearch HTTP DRIVER]: %s", err.Error())
}

func (d *HTTPD) BulkDumper(ch chan map[string]interface{}, index string, blksz int) {
	d.log("bulk dumper is up")
	bulk := d.cli.Bulk()
	for m := range ch {
		bulk = bulk.Add(elastic.NewBulkIndexRequest().Index(index).Doc(m))
		if bulk.NumberOfActions() >= blksz {
			log.Printf("dump new buffer with length = %d", blksz)
			if _, err := bulk.Do(context.Background()); err != nil {
				d.logFatal(err)
			}
			bulk.Reset()
		}
	}
	d.log(fmt.Sprintf("dump new buffer with length = %d", bulk.NumberOfActions()))
	_, err := bulk.Do(context.Background())
	if err != nil {
		d.logFatal(err)
	}
	bulk.Reset()
}
