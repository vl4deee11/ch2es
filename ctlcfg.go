package main

import (
	"ch2es/ch"
	"ch2es/es"
	"flag"
)

type conf struct {
	ThreadsNum int
	MaxOffset  int
	EsConf     *es.Conf

	ChConf *ch.Conf
}

func (c *conf) parse() {
	c.ChConf = new(ch.Conf)
	c.EsConf = new(es.Conf)

	flag.StringVar(&c.ChConf.Host, "ch-host", "0.0.0.0", "Clickhouse host (str)")
	flag.IntVar(&c.ChConf.Port, "ch-port", 8123, "Clickhouse http host (int)")
	flag.StringVar(&c.ChConf.OrderField, "ch-order", "", "Clickhouse order field (str)")
	flag.StringVar(
		&c.ChConf.Fields,
		"ch-fields",
		"*",
		"Clickhouse clickhouse fields for transfer ex: f_1,f_2,f_3 (str)",
	)

	flag.StringVar(&c.ChConf.Condition, "ch-cond", "1", "Clickhouse clickhouse where condition (str)")
	flag.StringVar(&c.ChConf.DB, "ch-db", "default", "Clickhouse db name (str)")
	flag.StringVar(&c.ChConf.Table, "ch-table", "", "Clickhouse table (str)")
	flag.IntVar(&c.ChConf.Limit, "ch-limit", 0, "Clickhouse limit (int)")
	flag.IntVar(&c.ChConf.ConnTimeout, "ch-timeout", 0, "Clickhouse connect timeout in seconds (int)")

	flag.StringVar(&c.EsConf.Host, "es-host", "0.0.0.0", "Elastic search host (str)")
	flag.IntVar(&c.EsConf.Port, "es-port", 9200, "Elastic search port (int)")
	flag.StringVar(&c.EsConf.Index, "es-idx", "", "Elastic search index (str)")
	flag.IntVar(&c.EsConf.BlkSz, "es-blksz", 0, "Elastic search bulk insert size (int)")

	flag.IntVar(&c.MaxOffset, "max-offset", 0, "Max offset in clickhouse table (int)")

	flag.IntVar(&c.ThreadsNum, "tn", 0, "Threads number for parallel bulk inserts (int)")
	flag.Parse()
}

func (c *conf) getReader() (*ch.Reader, error) {
	m := new(ch.Reader)
	if err := m.Init(c.ChConf); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *conf) getWriter() (*es.Writer, error) {
	m := new(es.Writer)
	if err := m.Init(c.EsConf); err != nil {
		return nil, err
	}

	return m, nil
}
