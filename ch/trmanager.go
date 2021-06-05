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

type TrManager struct {
	url string
	cli *http.Client
	lim int

	tempQBuff    *bytes.Buffer
	tempQBuffLen int
}

func (m *TrManager) Init(cfg *Conf) {
	// TODO: add tls and auth
	m.cli = &http.Client{
		Timeout: time.Duration(cfg.ConnTimeout) * time.Second,
	}
	m.url = m.httpF(cfg.Host, cfg.Port)
	m.tempQBuff = bytes.NewBufferString(fmt.Sprintf(
		"select %s from %s.%s where %s order by %s limit %d",
		cfg.Fields,
		cfg.DB,
		cfg.Table,
		cfg.Condition,
		cfg.OrderField,
		cfg.Limit,
	))
	m.tempQBuffLen = m.tempQBuff.Len()
	m.lim = cfg.Limit
}

func (m *TrManager) httpF(h string, p int) string {
	return fmt.Sprintf("http://%s:%m", h, p)
}

func (m *TrManager) get(q string) ([]map[string]interface{}, error) {
	_, err := m.tempQBuff.WriteString(q)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", m.url, m.tempQBuff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := m.cli.Do(req)
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

	m.tempQBuff.Truncate(m.tempQBuffLen)
	return bodyM["data"].([]map[string]interface{}), nil
}

func (m *TrManager) readInitialOffset() int {
	f, _ := ioutil.ReadFile("stats")
	r, err := strconv.Atoi(string(f))
	if err == nil {
		return r
	}
	return 0
}

func (m *TrManager) log(msg string) {
	log.Printf("[Clickhouse Tranfer Manager]: %s", msg)
}

func (m *TrManager) logFatal(err error) {
	log.Fatalf("[Clickhouse Tranfer Manager]: %s", err.Error())
}

func (m *TrManager) Run(ch chan map[string]interface{}, max int) {
	offset := m.readInitialOffset()

	m.log(fmt.Sprintf("start offset = %d", offset))
	for offset < max {
		data, err := m.get(
			fmt.Sprintf("offset %d format JSON", offset),
		)
		if err != nil {
			m.logFatal(err)
		}
		if len(data) == 0 {
			break
		}

		for i := range data {
			ch <- data[i]
		}
		offset += m.lim
		if err := ioutil.WriteFile("stats", []byte(fmt.Sprintf("%d", offset)), 0600); err != nil {
			m.logFatal(err)
		}
		m.log(fmt.Sprintf("current offset = %d", offset))
	}
}
