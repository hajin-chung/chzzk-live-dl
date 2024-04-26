package main

import (
	"fmt"
	"net/http"
	"os"
)

type Credentials struct {
	NID_AUT string `json:"NID_AUT"`
	NID_SES string `json:"NID_SES"`
}

type ChzzkClient struct {
	*http.Client
}

var client *ChzzkClient

func InitClient() {
	client = &ChzzkClient{
		Client: &http.Client{},
	}
}

func (c *ChzzkClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *ChzzkClient) Do(req *http.Request) (*http.Response, error) {
	NID_AUT := os.Getenv("NID_AUT")
	NID_SES := os.Getenv("NID_SES")
	req.Header.Set("Cookie", fmt.Sprintf("NID_AUT=%s; NID_SES=%s;", NID_AUT, NID_SES))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	return c.Client.Do(req)
}
