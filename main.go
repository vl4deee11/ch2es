package main

import (
	"log"
)

func main() {
	cfg := new(conf)
	cfg.parse()

	chTrManager, err := cfg.getChTrManager()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ch2es start")

	esTrManager, err := cfg.getEsTrManager()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan map[string]interface{})
	defer close(ch)

	for i := 0; i < cfg.ThreadsNum; i++ {
		go esTrManager.BulkDumper(ch)
	}

	chTrManager.Run(ch, cfg.MaxOffset)
	log.Println("successfully transferred")
}
