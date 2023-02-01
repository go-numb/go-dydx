package private

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-numb/go-dydx/helpers"
)

func (p *Private) get(endpoint string, params url.Values) ([]byte, error) {
	return p.request(http.MethodGet, helpers.GenerateQueryPath(endpoint, params), "")
}

func (p *Private) post(endpoint string, data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return p.request(http.MethodPost, endpoint, string(b))
}

func (p *Private) delete(endpoint string, params url.Values) ([]byte, error) {
	return p.request(http.MethodDelete, helpers.GenerateQueryPath(endpoint, params), "")
}

func (p *Private) request(method, endpoint string, data string) ([]byte, error) {
	isoTimestamp := generateNowISO()
	requestPath := fmt.Sprintf("/v3/%s", endpoint)
	headers := map[string]string{
		"DYDX-SIGNATURE":  p.Sign(requestPath, method, isoTimestamp, data),
		"DYDX-API-KEY":    p.ApiKeyCredentials.Key,
		"DYDX-TIMESTAMP":  isoTimestamp,
		"DYDX-PASSPHRASE": p.ApiKeyCredentials.Passphrase,
	}
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
		return nil, fmt.Errorf("uri: %v , status code: %d", requestPath, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	p.Logger.Printf("uri: %s, response body: %s", requestPath, b)
	return b, err
}

func (p *Private) execute(method string, requestPath string, headers map[string]string, data string) (*http.Response, error) {
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
		Timeout: 5 * time.Second,
	}
	return c.Do(req)

}

func generateNowISO() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
}

func (p *Private) Sign(requestPath, method, isoTimestamp, body string) string {
	message := fmt.Sprintf("%s%s%s%s", isoTimestamp, method, requestPath, body)
	secret, _ := base64.URLEncoding.DecodeString(p.ApiKeyCredentials.Secret)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
