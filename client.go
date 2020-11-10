package wooCommerce

import (
	"bytes"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type Client struct {
	HostURL             string
	Key                 string
	Secret              string
	NextQueryPageRegexp *regexp.Regexp
	HTTPClient          *http.Client
}

func NewClient(hostUrl string, key string, secret string, isCoCart bool) *Client {
	var version string

	if isCoCart {
		version = "cocart/v1/"
	} else {
		version = "wc/v3/"
	}

	re := regexp.MustCompile(`\<(.*)\>;.(rel="next")`)

	newClient := &Client{
		Key:                 key,
		Secret:              secret,
		NextQueryPageRegexp: re,
		HostURL:             hostUrl,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}

	newClient.setHostURL(version)

	return newClient
}

func (c *Client) setHostURL(version string) {
	c.HostURL = c.HostURL + wpAPIPrefix + version
}

func (c *Client) sendRequest(request *http.Request) (*http.Response, error) {

	request.Header.Set(contentTypeHeader, applicationJson)
	request.SetBasicAuth(c.Key, c.Secret)

	return c.HTTPClient.Do(request)
}

func (c *Client) getURL(endpoint string, parameters url.Values) (string, error) {
	wcURL, err := url.Parse(c.HostURL + endpoint)
	if err != nil {
		return "", err
	}

	slug, err := url.ParseRequestURI(wcURL.String())
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(wcURL.RawQuery)
	if err != nil {
		return "", err
	}

	for key, params := range parameters {
		for _, param := range params {
			q.Add(key, param)
		}
	}

	slug.RawQuery = q.Encode()

	return slug.String(), nil
}

func setHeaderCookie(request *http.Request, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
}

func (c *Client) Get(endpoint string, parameters url.Values) (*http.Response, error) {

	endpoint, err := c.getURL(endpoint, parameters)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(request)
}

func (c *Client) Post(endpoint string, json string, cookies []*http.Cookie) (*http.Response, error) {

	payload := bytes.NewReader([]byte(json))

	parameters := url.Values{"consumer_key": []string{c.Key}, "consumer_secret": []string{c.Secret}}

	endpoint, err := c.getURL(endpoint, parameters)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return nil, err
	}

	if cookies != nil {
		setHeaderCookie(request, cookies)
	}

	return c.sendRequest(request)
}

type HeaderKey string

const (
	HeaderKeyTotalPages HeaderKey = "X-Wp-Totalpages"
	HeaderKeyTotal      HeaderKey = "X-Wp-Total"
	HeaderLink          HeaderKey = "Link"
)

type QueryParam string

const (
	QueryParamPage    QueryParam = "page"
	QueryParamPerPage QueryParam = "per_page"
	QueryParamOffset  QueryParam = "offset"
)
