package main

import (
	"ch2es/ch"
	"ch2es/common"
	"ch2es/es"
	"flag"
)

type conf struct {
	ThreadsNum int `desc:"threads number"`
	MaxOffset  int `desc:"max offset"`
	EsConf     *es.Conf

	ChConf *ch.Conf
}

func (c *conf) parse() {
	c.ChConf = &ch.Conf{HTTPConf: new(common.HTTPConf)}
	c.EsConf = &es.Conf{HTTPConf: new(common.HTTPConf)}

	flag.StringVar(&c.ChConf.Protocol, "ch-protocol", "http", "Clickhouse protocol (str)")
	flag.StringVar(&c.ChConf.Host, "ch-host", "0.0.0.0", "Clickhouse host (str)")
	flag.IntVar(&c.ChConf.Port, "ch-port", 8123, "Clickhouse http host (int)")
	flag.StringVar(&c.ChConf.OrderField, "ch-order", "", "Clickhouse order field (str)")
	flag.StringVar(&c.ChConf.Fields, "ch-fields", "*", "Clickhouse clickhouse fields for transfer ex: f_1,f_2,f_3 (str)")
	flag.StringVar(&c.ChConf.Condition, "ch-cond", "1", "Clickhouse clickhouse where condition (str)")
	flag.StringVar(&c.ChConf.DB, "ch-db", "default", "Clickhouse db name (str)")
	flag.StringVar(&c.ChConf.Table, "ch-table", "", "Clickhouse table (str)")
	flag.StringVar(&c.ChConf.URLParams.User, "ch-user", "", "Clickhouse db username (str)")
	flag.StringVar(&c.ChConf.URLParams.Pass, "ch-pass", "", "Clickhouse db password (str)")
	flag.IntVar(&c.ChConf.Limit, "ch-limit", 100, "Clickhouse limit (int)")
	flag.IntVar(&c.ChConf.ConnTimeoutSec, "ch-conn-timeout", 20, "Clickhouse connect timeout in sec (int)")
	flag.IntVar(&c.ChConf.QueryTimeoutSec, "ch-query-timeout", 60, "Clickhouse query timeout in sec (int)")

	flag.StringVar(&c.EsConf.Protocol, "es-protocol", "http", "Elastic search protocol (str)")
	flag.StringVar(&c.EsConf.Host, "es-host", "0.0.0.0", "Elastic search host (str)")
	flag.IntVar(&c.EsConf.Port, "es-port", 9200, "Elastic search port (int)")
	flag.IntVar(&c.EsConf.QueryTimeoutSec, "es-query-timeout", 60, "Elastic search query timeout in sec (int)")
	flag.StringVar(&c.EsConf.User, "es-user", "", "Elastic search username (str)")
	flag.StringVar(&c.EsConf.Pass, "es-pass", "", "Elastic search password (str)")
	flag.StringVar(&c.EsConf.Index, "es-idx", "", "Elastic search index (str)")
	flag.IntVar(&c.EsConf.BlkSz, "es-blksz", 0, "Elastic search bulk insert size (int)")

	flag.IntVar(&c.MaxOffset, "max-offset", 0, "Max offset in clickhouse table (int)")

	flag.IntVar(&c.ThreadsNum, "tn", 0, "Threads number for parallel insert and read (int)")
	flag.Parse()
	c.print()
}

func (c *conf) print() {
	c.ChConf.Print()
	c.EsConf.Print()
	common.PrintFromDesc("[COMMON CONFIG]", *c)
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
