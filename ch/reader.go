package ch

import (
	"bytes"
	"ch2es/ch/cursor"
	"ch2es/log"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Reader struct {
	dsnURL   string
	cli      *http.Client
	qTimeout time.Duration
	cursorT  cursor.T
}

func NewReader(cfg *Conf) (*Reader, cursor.Cursor, error) {
	r := &Reader{
		cursorT: cursor.T(cfg.CursorT),
	}

	switch r.cursorT {
	case cursor.JSONFile:
		cur, err := cursor.NewJSONFile(cfg.JFC)
		return r, cur, err
	case cursor.Stdin:
		cur, err := cursor.NewStdin(cfg.StdinC)
		return r, cur, err
	case cursor.Offset, cursor.Timestamp:
		cfg.BuildHTTP()
		r.qTimeout = time.Duration(cfg.QueryTimeoutSec) * time.Second
		r.cli = &http.Client{
			Timeout: time.Duration(cfg.ConnTimeoutSec) * time.Second,
		}
		r.dsnURL = r.httpF(cfg)
		var cur cursor.Cursor
		if cfg.CursorT == 0 {
			cur = cursor.NewOffset(cfg.OFC, cfg.Fields, cfg.DB, cfg.Table, cfg.Condition)
		} else {
			cur = cursor.NewTimestamp(cfg.TSC, cfg.Fields, cfg.DB, cfg.Table, cfg.Condition)
		}

		err := r.ping(cfg)
		return r, cur, err
	}
	return nil, nil, fmt.Errorf("bad Cursor")
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

func (r *Reader) Read(ch chan *bytes.Buffer, wCh chan map[string]interface{}, eCh chan error, wg *sync.WaitGroup) {
	log.Info("start new clickhouse reader")
	switch r.cursorT {
	case cursor.JSONFile, cursor.Stdin:
		r.readIO(ch, wCh, eCh, wg)
	case cursor.Offset, cursor.Timestamp:
		r.readHTTP(ch, wCh, eCh, wg)
	}
}

func (r *Reader) readHTTP(ch chan *bytes.Buffer, wCh chan map[string]interface{}, eCh chan error, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		log.Info("reader is stop")
	}()

	for q := range ch {
		data, err := r.get(q)
		if err != nil {
			eCh <- err
			break
		}

		for i := range data {
			wCh <- data[i].(map[string]interface{})
		}
	}
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
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return bodyM["data"].([]interface{}), nil
}

func (r *Reader) readIO(ch chan *bytes.Buffer, wCh chan map[string]interface{}, eCh chan error, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		log.Info("reader is stop")
	}()

	for b := range ch {
		data := map[string]interface{}{}
		if err := json.Unmarshal(b.Bytes(), &data); err != nil {
			eCh <- err
			break
		}
		wCh <- data
	}
}
