package ch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Reader struct {
	dsnURL string
	cli    *http.Client
	lim    int

	tempQBuff    *bytes.Buffer
	tempQBuffLen int
	qTimeout     time.Duration
}

func (r *Reader) Init(cfg *Conf) error {
	cfg.BuildHTTP()
	r.qTimeout = time.Duration(cfg.QueryTimeoutSec) * time.Second
	r.cli = &http.Client{
		Timeout: time.Duration(cfg.ConnTimeoutSec) * time.Second,
	}
	r.dsnURL = r.httpF(cfg)
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

	return r.ping(cfg)
}

func (r *Reader) ping(cfg *Conf) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.qTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/ping", cfg.URL),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.cli.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status = %d", resp.StatusCode)
	}

	_ = resp.Body.Close()
	return nil
}

func (r *Reader) httpF(cfg *Conf) string {
	url := fmt.Sprintf("%s/%s", cfg.URL, cfg.DB)
	params := make([]string, 0)

	v := reflect.ValueOf(cfg.URLParams)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		p, ok := t.Field(i).Tag.Lookup("param")
		if !ok {
			continue
		}
		field := v.Field(i).Interface()
		if field == "" {
			continue
		}
		if field == "" {
			continue
		}
		params = append(params, fmt.Sprintf("%s=%s", p, field))
	}
	if len(params) > 0 {
		url += "?" + strings.Join(params, "&")
	}
	return url
}

func (r *Reader) get(buff *bytes.Buffer) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.qTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		r.dsnURL,
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
		log.Println(string(body))
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return bodyM["data"].([]interface{}), nil
}

func (r *Reader) Read(
	rCh chan string,
	wCh chan map[string]interface{},
	eCh chan error,
	wg *sync.WaitGroup,
) {
	log.Println("start new clickhouse reader")
	defer func() {
		wg.Done()
		log.Println("reader is stop")
	}()
	buff := bytes.NewBuffer(r.tempQBuff.Bytes())
	for q := range rCh {
		_, err := buff.WriteString(q)
		if err != nil {
			eCh <- err
			break
		}

		data, err := r.get(buff)
		if err != nil {
			eCh <- err
			break
		}
		if len(data) == 0 {
			eCh <- nil
			break
		}

		for i := range data {
			m := data[i].(map[string]interface{})
			for k := range m {
				kk := strings.ReplaceAll(k, ".", "_")
				if kk != k {
					m[kk] = m[k]
					delete(m, k)
				}
			}

			wCh <- m
		}
		buff.Truncate(r.tempQBuffLen)
	}
}
