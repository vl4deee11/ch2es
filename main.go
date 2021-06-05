package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	log.Println(`This tool is not support user and passwords`)
	cfg := new(conf)
	cfg.parse()

	chHTTPD, err := cfg.getChHTTPD()
	if err != nil {
		log.Fatal(err)
	}
	stat, _ := ioutil.ReadFile("stats")
	offset := 0
	idx, err := strconv.Atoi(string(stat))
	if err == nil {
		offset = idx
	}

	log.Println("ch2es start")
	log.Printf("start offset %d\n", offset)
	esHTTPD, err := cfg.getEsHTTPD()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan map[string]interface{})
	defer close(ch)

	for i := 0; i < cfg.ThreadsNum; i++ {
		go esHTTPD.BulkDumper(ch, cfg.EsIndex, cfg.EsBlkSz)
	}

	var exReason error = nil
	for offset < cfg.MaxOffset {
		data, err := chHTTPD.Get(fmt.Sprintf(
			"select * from %s.%s order by %s limit %d offset %d format JSON",
			cfg.ChDB,
			cfg.ChTable,
			cfg.ChOrderField,
			cfg.ChStepSz,
			offset,
		))
		if err != nil {
			exReason = err
			break
		}
		if len(data) == 0 {
			break
		}

		for i := range data {
			ch <- data[i]
		}
		offset += cfg.ChStepSz
		if err := ioutil.WriteFile("stats", []byte(fmt.Sprintf("%d", offset)), 0600); err != nil {
			exReason = err
			break
		}
		log.Println("current offset =", offset)
	}
	if exReason != nil {
		log.Print(exReason)
	}
	log.Println("successfully transferred")
}
