package utils

import (
	"QuakeAPI/log"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Http interface {
	DoGet(url string, data map[string]string, headers map[string]string) []byte
	DoPost(url string, data map[string]string, headers map[string]string) []byte
}

type HttpClient struct {
}

func (h *HttpClient) DoGet(url string, data map[string]string, headers map[string]string) []byte {
	return doRequest("GET", url, data, headers)
}

func (h *HttpClient) DoPost(url string, data map[string]string, headers map[string]string) []byte {
	return doRequest("POST", url, data, headers)
}

func doRequest(
	method string,
	url string,
	data map[string]string,
	headers map[string]string) []byte {
	client := http.Client{Timeout: 30 * time.Second}

	var body []byte
	if rawBody, ok := data["_raw_body"]; ok {
		body = []byte(rawBody)
	} else {
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			log.Log("json marshal error:"+err.Error(), log.ERROR)
			return nil
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Log("http error:"+err.Error(), log.ERROR)
		return nil
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Log("request error:"+err.Error(), log.ERROR)
		return nil
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Log("read error:"+err.Error(), log.ERROR)
			return nil
		}
		return body
	}
	return nil
}
