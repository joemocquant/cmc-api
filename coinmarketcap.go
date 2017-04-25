// Coinmarket cap API implementation.
//
// API Doc: http://coinmarketcap.com/api/
//
// Please limit requests to no more than 10 per minute.
// Endpoints update every 5 minutes.
package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var conf *configuration

type Client struct {
	httpClient *http.Client
	throttle   <-chan time.Time
}

type configuration struct {
	CoinmarketcapConf `json:"coinmarketcap"`
}

type CoinmarketcapConf struct {
	APIUrl               string `json:"api_url"`
	HTTPClientTimeoutSec int    `json:"httpclient_timeout_sec"`
	MaxRequestsMin       int    `json:"max_requests_min"`
	LogLevel             string `json:"log_level"`
}

func init() {

	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	content, err := ioutil.ReadFile("conf.json")

	if err != nil {
		log.WithField("error", err).Fatal("loading configuration")
	}

	if err := json.Unmarshal(content, &conf); err != nil {
		log.WithField("error", err).Fatal("loading configuration")
	}

	switch conf.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}

// NewClient returns a newly configured client
func NewClient() *Client {

	reqInterval := 60 * time.Second / time.Duration(conf.MaxRequestsMin)

	client := http.Client{
		Timeout: time.Duration(conf.HTTPClientTimeoutSec) * time.Second,
	}

	return &Client{&client, time.Tick(reqInterval)}
}

// Do prepares and executes api call requests.
func (c *Client) do(endpoint string, params map[string]string) ([]byte, error) {

	url := buildUrl(endpoint, params)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v (API command: %s)",
			err, params["command"])
	}

	req.Header.Add("Accept", "application/json")

	type result struct {
		resp *http.Response
		err  error
	}

	done := make(chan result)
	go func() {
		<-c.throttle
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	res := <-done

	if res.err != nil {
		return nil, fmt.Errorf("http.Client.Do: %v", res.err)
	}

	defer res.resp.Body.Close()

	body, err := ioutil.ReadAll(res.resp.Body)
	if err != nil {
		return body, fmt.Errorf("ioutil.readAll: %v", err)
	}

	if res.resp.StatusCode != 200 {
		return body, fmt.Errorf("status code: %s (API endpoint: %s)",
			res.resp.Status, endpoint)
	}

	return body, nil
}

func buildUrl(endpoint string, params map[string]string) string {

	u := conf.APIUrl + "/" + endpoint + "?"

	var parameters []string
	for k, v := range params {
		parameters = append(parameters, k+"="+url.QueryEscape(v))
	}

	return u + strings.Join(parameters, "&")
}
