package main

import (
	"ch2es/ch"
	"ch2es/es"
	"flag"
)

type conf struct {
	ThreadsNum int
	MaxOffset  int
	EsHost     string
	EsPort     int
	EsIndex    string
	EsBlkSz    int

	ChHost       string
	ChUser       string
	ChPass       string
	ChPort       int
	ChDB         string
	ChOrderField string
	ChTable      string
	ChStepSz     int
}

func (c *conf) parse() {
	flag.StringVar(&c.ChHost, "ch-host", "0.0.0.0", "Clickhouse host (str)")
	flag.IntVar(&c.ChPort, "ch-port", 8123, "Clickhouse http host (int)")
	flag.StringVar(&c.ChOrderField, "ch-order", "", "Clickhouse order field (str)")
	flag.StringVar(&c.ChDB, "ch-db", "default", "Clickhouse db name (str)")
	flag.StringVar(&c.ChTable, "ch-table", "", "Clickhouse table (str)")
	flag.IntVar(&c.ChStepSz, "step", 0, "Step size (int)")

	flag.StringVar(&c.EsHost, "es-host", "0.0.0.0", "Elastic search host (str)")
	flag.IntVar(&c.EsPort, "es-port", 9200, "Elastic search port (int)")
	flag.StringVar(&c.EsIndex, "es-idx", "", "Elastic search index (str)")
	flag.IntVar(&c.EsBlkSz, "es-blksz", 0, "Elastic search bulk insert size (int)")

	flag.IntVar(&c.MaxOffset, "max-offset", 0, "Max offset in clickhouse table (int)")

	flag.IntVar(&c.ThreadsNum, "tn", 0, "Threads number for parallel bulk inserts (int)")
	flag.Parse()
}

func (c *conf) getChHTTPD() (*ch.HTTPD, error) {
	d := new(ch.HTTPD)
	d.SetHTTP(c.ChHost, c.ChPort)
	return d, nil
}

func (c *conf) getEsHTTPD() (*es.HTTPD, error) {
	d := new(es.HTTPD)
	if err := d.SetHTTP(c.EsHost, c.EsPort); err != nil {
		return nil, err
	}

	return d, nil
}
