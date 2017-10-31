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

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	conf   *configuration
	logger *logrus.Entry
)

type Client struct {
	httpClient *http.Client
	throttle   <-chan time.Time
}

type configuration struct {
	apiConf `json:"coinmarketcap_api"`
}

type apiConf struct {
	APIUrl               string `json:"api_url"`
	HTTPClientTimeoutSec int    `json:"httpclient_timeout_sec"`
	MaxRequestsMin       int    `json:"max_requests_min"`
	LogLevel             string `json:"log_level"`
}

func init() {

	customFormatter := new(prefixed.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true
	customFormatter.ForceFormatting = true
	logrus.SetFormatter(customFormatter)

	logger = logrus.WithField("prefix", "[api:coinmarketcap]")

	content, err := ioutil.ReadFile("conf.json")

	if err != nil {
		logger.WithField("error", err).Fatal("loading configuration")
	}

	if err := json.Unmarshal(content, &conf); err != nil {
		logger.WithField("error", err).Fatal("loading configuration")
	}

	switch conf.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
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
		return nil, fmt.Errorf("http.NewRequest: %v (API endpoint: %s)",
			err, endpoint)
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

	u := conf.APIUrl + "/" + endpoint

	var parameters []string
	for k, v := range params {
		parameters = append(parameters, k+"="+url.QueryEscape(v))
	}

	if len(parameters) > 0 {
		return u + "?" + strings.Join(parameters, "&")
	}

	return u
}
