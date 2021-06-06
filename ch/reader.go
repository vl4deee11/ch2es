package ch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Reader struct {
	url string
	cli *http.Client
	lim int

	tempQBuff    *bytes.Buffer
	tempQBuffLen int
}

func (r *Reader) Init(cfg *Conf) error {
	// TODO: add tls and auth
	r.cli = &http.Client{
		Timeout: time.Duration(cfg.ConnTimeout) * time.Second,
	}
	r.url = r.httpF(cfg.Host, cfg.Port)
	r.tempQBuff = bytes.NewBufferString(fmt.Sprintf(
		"select %s from %s.%s where %s order by %s limit %d ",
		cfg.Fields,
		cfg.DB,
		cfg.Table,
		cfg.Condition,
		cfg.OrderField,
		cfg.Limit,
	))
	r.tempQBuffLen = r.tempQBuff.Len()
	r.lim = cfg.Limit

	req, err := http.NewRequestWithContext(context.Background(), "GET", r.url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.cli.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code = %d", resp.StatusCode)
	}
	return nil
}

func (r *Reader) httpF(h string, p int) string {
	return fmt.Sprintf("http://%s:%d", h, p)
}

func (r *Reader) get(buff *bytes.Buffer) ([]interface{}, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		r.url,
		bytes.NewBuffer(buff.Bytes()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.cli.Do(req)
	if err != nil {
		return nil, err
	}
	bodyM := map[string]interface{}{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &bodyM); err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}


	return bodyM["data"].([]interface{}), nil
}

func (r *Reader) Read(rCh chan string, wCh chan map[string]interface{}) {
	buff := bytes.NewBuffer(r.tempQBuff.Bytes())
	for q := range rCh {
		_, err := buff.WriteString(q)
		if err != nil {
			log.Fatal(err)
		}

		data, err := r.get(buff)
		if err != nil {
			log.Fatal(err)
		}
		if len(data) == 0 {
			break
		}

		for i := range data {
			wCh <- data[i].(map[string]interface{})
		}
		buff.Truncate(r.tempQBuffLen)
	}
}

func (r *Reader) ReadInitialOffset() int {
	f, _ := ioutil.ReadFile("stats")
	n, err := strconv.Atoi(string(f))
	if err == nil {
		return n
	}
	return 0
}
