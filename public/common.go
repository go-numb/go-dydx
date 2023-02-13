package public

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/huangjosh/go-dydx/helpers"
)

func (p *Public) get(endpoint string, params url.Values) ([]byte, error) {
	return p.request(http.MethodGet, helpers.GenerateQueryPath(endpoint, params), "")
}

func (p *Public) post(endpoint string, data interface{}) ([]byte, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return p.request(http.MethodPost, endpoint, string(d))
}

func (p *Public) delete(endpoint string, params url.Values) ([]byte, error) {
	return p.request(http.MethodGet, helpers.GenerateQueryPath(endpoint, params), "")
}

func (p *Public) request(method, endpoint string, data string) ([]byte, error) {
	requestPath := fmt.Sprintf("/v3/%s", endpoint)
	headers := map[string]string{}
	res, err := p.execute(method, requestPath, headers, data)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Rate Limit
	p.RateLimit.Remaining = res.Header.Get("RateLimit-Remaining")
	p.RateLimit.Reset = res.Header.Get("RateLimit-Reset")
	p.RateLimit.RetryAfter = res.Header.Get("Retry-After")
	p.RateLimit.Limit = res.Header.Get("RateLimit-Limit")

	if res.StatusCode < 200 || res.StatusCode > 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		p.Logger.Printf("uri: %s, code: %d, err msg: %s", requestPath, res.StatusCode, buf.String())
		return nil, fmt.Errorf("uri:%v , status cod e: %d", requestPath, res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	p.Logger.Printf("uri: %s,response body:%s", requestPath, b)
	return b, err
}

func (p *Public) execute(method string, requestPath string, headers map[string]string, data string) (*http.Response, error) {
	requestPath = fmt.Sprintf("%s%s", p.Host, requestPath)
	req, err := http.NewRequest(method, requestPath, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "go-dydx")

	c := &http.Client{
		Timeout: time.Second * 5,
	}
	return c.Do(req)
}
