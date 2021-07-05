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
	dsnURL      string
	cli         *http.Client
	qTimeout    time.Duration
	cursorT     cursorT
	dotReplacer string
}

func NewReader(cfg *Conf) (*Reader, Cursor, error) {
	r := new(Reader)
	r.dotReplacer = cfg.DotReplacer
	r.cursorT = cursorT(cfg.CursorT)

	switch r.cursorT {
	case fileCursor:
		cur, err := NewJSONFileCursor(cfg.JFC)
		return r, cur, err
	case offsetCursor, timeStampCursor:
		cfg.BuildHTTP()
		r.qTimeout = time.Duration(cfg.QueryTimeoutSec) * time.Second
		r.cli = &http.Client{
			Timeout: time.Duration(cfg.ConnTimeoutSec) * time.Second,
		}
		r.dsnURL = r.httpF(cfg)
		var cur Cursor
		if cfg.CursorT == 0 {
			cur = NewOffsetCursor(cfg)
		} else {
			cur = NewTimestampCursor(cfg)
		}

		err := r.ping(cfg)
		return r, cur, err
	}
	return nil, nil, fmt.Errorf("bad cursor")
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

func (r *Reader) Read(
	ch chan *bytes.Buffer,
	wCh chan map[string]interface{},
	eCh chan error,
	wg *sync.WaitGroup,
) {
	log.Println("start new clickhouse reader")
	switch r.cursorT {
	case fileCursor:
		r.readFile(ch, wCh, eCh, wg)
	case offsetCursor, timeStampCursor:
		r.readHTTP(ch, wCh, eCh, wg)
	}
}

func (r *Reader) readHTTP(
	ch chan *bytes.Buffer,
	wCh chan map[string]interface{},
	eCh chan error,
	wg *sync.WaitGroup,
) {
	defer func() {
		wg.Done()
		log.Println("reader is stop")
	}()

	cleaner := r.getJSONCleaner()
	for q := range ch {
		data, err := r.get(q)
		if err != nil {
			eCh <- err
			break
		}

		for i := range data {
			wCh <- cleaner(data[i].(map[string]interface{}))
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

func (r *Reader) readFile(
	ch chan *bytes.Buffer,
	wCh chan map[string]interface{},
	eCh chan error,
	wg *sync.WaitGroup,
) {
	defer func() {
		wg.Done()
		log.Println("reader is stop")
	}()

	cleaner := r.getJSONCleaner()
	for b := range ch {
		data := map[string]interface{}{}
		if err := json.Unmarshal(b.Bytes(), &data); err != nil {
			eCh <- err
			break
		}
		wCh <- cleaner(data)
	}
}

func (r *Reader) getJSONCleaner() func(d map[string]interface{}) map[string]interface{} {
	if r.dotReplacer != "" {
		return func(d map[string]interface{}) map[string]interface{} {
			m := d
			for k := range m {
				kk := strings.ReplaceAll(k, ".", r.dotReplacer)
				if kk != k {
					m[kk] = m[k]
					delete(m, k)
				}
			}
			return m
		}
	}

	return func(d map[string]interface{}) map[string]interface{} {
		return d
	}
}

//nolint:unused //set TODO
// TODO: parse dots from fields name and generate sub-maps
func (r *Reader) rebuildJSON(d map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for k := range d {
		if subNames := strings.Split(k, "."); len(subNames) > 1 {
			fmt.Println(subNames)
			cm, ok := m[subNames[0]]
			if !ok {
				m[subNames[0]] = make(map[string]interface{})
				cm = m[subNames[0]]
			}
			tm := cm.(map[string]interface{})
			for i := 1; i < len(subNames); i++ {
				_, ok = tm[subNames[i]]
				if !ok {
					tm[subNames[i]] = make(map[string]interface{})
					tm = tm[subNames[i]].(map[string]interface{})
				}
			}

			tm[subNames[len(subNames)-1]] = d[k]
		} else {
			m[subNames[0]] = d[k]
		}
	}
	return m
}
