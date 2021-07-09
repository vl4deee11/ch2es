package util

import "fmt"

type HTTPConf struct {
	Protocol string `desc:"protocol"`
	Host     string `desc:"host"`
	Port     int    `desc:"port"`
	URL      string
}

func (c *HTTPConf) BuildHTTP() {
	c.URL = fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
}
