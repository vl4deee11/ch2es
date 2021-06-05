package ch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPD struct {
	url string
}

func (d *HTTPD) SetHTTP(h string, p int) {
	// TODO: add tls and auth
	d.url = d.httpF(h, p)
}

func (d *HTTPD) httpF(h string, p int) string {
	return fmt.Sprintf("http://%s:%d", h, p)
}

func (d *HTTPD) Get(q string) ([]map[string]interface{}, error) {
	resp, err := http.Post(d.url, "application/json", bytes.NewReader([]byte(q)))
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
