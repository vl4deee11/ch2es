package main

import (
	"ch2es/ch"
	"ch2es/es"
	"ch2es/util"
	"flag"
)

type conf struct {
	ThreadsNum int `desc:"threads number"`
	EsConf     *es.Conf

	ChConf *ch.Conf
}

func (c *conf) parse() {
	c.ChConf = &ch.Conf{
		HTTPConf: new(util.HTTPConf),
		TSC:      new(ch.TSCursorConf),
		OFC:      new(ch.OffsetCursorConf),
		JFC:      new(ch.JSONFileCursorConf),
		StdinC:   new(ch.StdInCursorConf),
	}
	c.EsConf = &es.Conf{HTTPConf: new(util.HTTPConf)}

	// CLICKHOUSE
	flag.StringVar(&c.ChConf.Protocol, "ch-protocol", "http", "[Clickhouse] protocol")
	flag.StringVar(&c.ChConf.Host, "ch-host", "0.0.0.0", "[Clickhouse] host")
	flag.IntVar(&c.ChConf.Port, "ch-port", 8123, "[Clickhouse] http host")
	flag.StringVar(&c.ChConf.Fields, "ch-fields", "*", "[Clickhouse] fields for transfer ex: f_1,f_2,f_3")
	flag.StringVar(&c.ChConf.Condition, "ch-cond", "", "[Clickhouse] where condition")
	flag.StringVar(&c.ChConf.DB, "ch-db", "default", "[Clickhouse] db name")
	flag.StringVar(&c.ChConf.Table, "ch-table", "", "[Clickhouse] table")
	flag.StringVar(&c.ChConf.DotReplacer, "ch-dot-replacer", "", "[Clickhouse] Replacer for dots in fields if need")
	flag.StringVar(&c.ChConf.URLParams.User, "ch-user", "", "[Clickhouse] db username")
	flag.StringVar(&c.ChConf.URLParams.Pass, "ch-pass", "", "[Clickhouse] db password")
	flag.IntVar(&c.ChConf.ConnTimeoutSec, "ch-conn-timeout", 20, "[Clickhouse] connect timeout in sec")
	flag.IntVar(&c.ChConf.QueryTimeoutSec, "ch-query-timeout", 60, "[Clickhouse] query timeout in sec")
	flag.IntVar(&c.ChConf.CursorT, "ch-cursor", 0, "[Clickhouse] cursor type. Available 0 (offset cursor), 1 (timestamp cursor), 2 (json file cursor), 3 (stdin cursor)")

	// CLICKHOUSE timestamp cursor
	flag.IntVar(&c.ChConf.TSC.StepSec, "ch-tsc-step", 0, "[Clickhouse timestamp cursor] step in sec. Use only if --ch-cursor=1")
	flag.IntVar(&c.ChConf.TSC.Min, "ch-tsc-min", 0, "[Clickhouse timestamp cursor] start time format unix timestamp. Use only if --ch-cursor=1")
	flag.IntVar(&c.ChConf.TSC.Max, "ch-tsc-max", 0, "[Clickhouse timestamp cursor] end time format unix timestamp. Use only if --ch-cursor=1")
	flag.StringVar(&c.ChConf.TSC.Field, "ch-tsc-field", "", "[Clickhouse timestamp cursor] field. Should be datetime type or timestamp. Use only if --ch-cursor=1")

	// CLICKHOUSE offset cursor
	flag.StringVar(&c.ChConf.OFC.OrderField, "ch-ofc-order", "", "[Clickhouse offset cursor] order field. Use only if --ch-cursor=0 (by default)")
	flag.IntVar(&c.ChConf.OFC.Limit, "ch-ofc-limit", 100, "[Clickhouse offset cursor] limit. Use only if --ch-cursor=0 (by default)")
	flag.IntVar(&c.ChConf.OFC.Offset, "ch-ofc-offset", 0, "[Clickhouse offset cursor] start offset. Use only if --ch-cursor=0 (by default)")
	flag.IntVar(&c.ChConf.OFC.MaxOffset, "ch-ofc-max-offset", 0, "[Clickhouse offset cursor] max offset in clickhouse table. Use only if --ch-cursor=0 (by default)")

	// CLICKHOUSE json file cursor
	flag.StringVar(&c.ChConf.JFC.File, "ch-jfc-file", "", "[Clickhouse json file cursor] path to file with data formatted JSONEachRow. Use only if --ch-cursor=2")
	flag.IntVar(&c.ChConf.JFC.Line, "ch-jfc-line", 0, "[Clickhouse json file cursor] start line in file with data formatted JSONEachRow. Use only if --ch-cursor=2")

	// CLICKHOUSE stdin cursor
	flag.IntVar(&c.ChConf.StdinC.Line, "ch-stdinc-line", 0, "[Clickhouse stdin cursor] start line in stdin with data formatted JSONEachRow. Use only if --ch-cursor=3")

	// ELASTIC
	flag.StringVar(&c.EsConf.Protocol, "es-protocol", "http", "[Elasticsearch] protocol")
	flag.StringVar(&c.EsConf.IDField, "es-id-field", "", "[Elasticsearch] id field")
	flag.StringVar(&c.EsConf.Host, "es-host", "0.0.0.0", "[Elasticsearch] search host")
	flag.IntVar(&c.EsConf.Port, "es-port", 9200, "[Elasticsearch] search port")
	flag.IntVar(&c.EsConf.QueryTimeoutSec, "es-query-timeout", 60, "[Elasticsearch] search query timeout in sec")
	flag.StringVar(&c.EsConf.User, "es-user", "", "[Elasticsearch] search username")
	flag.StringVar(&c.EsConf.Pass, "es-pass", "", "[Elasticsearch] search password")
	flag.StringVar(&c.EsConf.Index, "es-idx", "", "[Elasticsearch] search index")
	flag.IntVar(&c.EsConf.BlkSz, "es-blksz", 0, "[Elasticsearch] search bulk insert size")

	// COMMON
	flag.IntVar(&c.ThreadsNum, "tn", 0, "[Common] Threads number for parallel insert and read")
	flag.Parse()
	c.print()
}

func (c *conf) print() {
	c.ChConf.Print()
	c.EsConf.Print()
	util.PrintFromDesc("[COMMON CONFIG]", *c)
}
