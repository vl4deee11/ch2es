package ch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPD struct {
	url string
	cli *http.Client
	defaultCtx context.Context
}

func (d *HTTPD) SetHTTP(h string, p int) {
	// TODO: add tls and auth
	d.cli = &http.Client{
		Timeout: 60 * time.Second,
	}
	d.url = d.httpF(h, p)

}

func (d *HTTPD) httpF(h string, p int) string {
	return fmt.Sprintf("http://%s:%d", h, p)
}

func (d *HTTPD) Get(q string) ([]map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(context.Background(), "POST", d.url, bytes.NewBuffer([]byte(q)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := d.cli.Do(req)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return m["data"].([]map[string]interface{}), nil
}
