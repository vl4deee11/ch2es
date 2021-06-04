package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
)

type ES struct {
	Content string `json:"content"`
}

type conf struct {
	ThreadsNum   int
	MaxOffset    int
	EsURL        string
	EsIndex      string
	EsBlkSz      int
	ChURL        string
	ChOrderField string
	ChTable      string
	ChStepSz     int
}

func main() {
	cfg := new(conf)
	flag.StringVar(&cfg.ChURL, "ch-url", "", "Clickhouse dsn URL (str)")
	flag.StringVar(&cfg.ChOrderField, "ch-order", "", "Clickhouse order field (str)")
	flag.StringVar(&cfg.ChTable, "ch-table", "", "Clickhouse table (str)")
	flag.StringVar(&cfg.EsURL, "es-url", "", "Elastic search dsn URL (str)")
	flag.StringVar(&cfg.EsIndex, "es-idx", "", "Elastic search index (str)")
	flag.IntVar(&cfg.MaxOffset, "max-offset", 0, "Max offset in clickhouse table (int)")
	flag.IntVar(&cfg.ChStepSz, "step", 0, "Step size (int)")
	flag.IntVar(&cfg.EsBlkSz, "es-blksz", 0, "Elastic search bulk insert size (int)")
	flag.IntVar(&cfg.ThreadsNum, "tn", 0, "Threads number for parallel bulk inserts (int)")
	flag.Parse()
	stat, _ := ioutil.ReadFile("stats")
	offset := 0
	idx, err := strconv.Atoi(string(stat))
	if err == nil {
		offset = idx
	}
	log.Println("ch2es start")
	log.Printf("start offset %d\n", offset)
	cliEs, err := elastic.NewClient(
		elastic.SetURL(cfg.EsURL),
		elastic.SetHealthcheckInterval(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan *ES)
	defer close(ch)
	for i := 0; i < cfg.ThreadsNum; i++ {
		go dumper(ch, cliEs, cfg)
	}
	for offset < cfg.MaxOffset {
		q := fmt.Sprintf("select * from %s order by %s limit %d offset %d format JSON", cfg.ChTable, cfg.ChOrderField, cfg.ChStepSz, offset)
		resp, err := http.Post(cfg.ChURL, "application/json", bytes.NewReader([]byte(q)))
		if err != nil {
			log.Fatal(err)
		}
		m := map[string]interface{}{}
		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &m); err != nil {
			log.Fatal(err)
		}
		data := m["data"].([]interface{})
		if len(data) == 0 {
			break
		}
		for i := range data {
			j, err := json.Marshal(data[i])
			if err != nil {
				log.Fatal(err)
			}
			ch <- &ES{string(j)}
		}
		offset += cfg.ChStepSz
		if err := ioutil.WriteFile("stats", []byte(fmt.Sprintf("%d", offset)), 0644); err != nil {
			log.Fatal(err)
		}
		log.Println("current offset =", offset)
	}
	log.Println("successfully transferred")
}
