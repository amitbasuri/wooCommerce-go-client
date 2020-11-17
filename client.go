package wooCommerce

import (
	"bytes"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type Client interface {
	GetTaxes(params url.Values) ([]Tax, error)
	GetTaxesPaginated(params url.Values) (*GetTaxesResponse, error)
	GetOrder(id int) (*Order, error)
	CreateOrder(o *Order) (*Order, error)
	SystemStatus() (*SystemStatus, error)
	GetSettings(key SettingsKey) (*Setting, error)
	QueryProductVariations(productId int, params url.Values) (*QueryProductsVariationResponse, error)
	QueryProducts(params url.Values) (*QueryProductsResponse, error)
	CreateWebhook(w *Webhook) error
	DeleteWebhook(id int, force bool) error
	CalculateShipping(shippingCart ShippingCart) (*ShippingResponse, error)
	QueryWebhooks(deliveryUrl *string) ([]Webhook, error)
}

type clientImpl struct {
	HostURL             string
	Key                 string
	Secret              string
	NextQueryPageRegexp *regexp.Regexp
	HTTPClient          *http.Client
}

func NewClient(hostUrl string, key string, secret string) Client {

	re := regexp.MustCompile(`\<(.*)\>;.(rel="next")`)

	newClient := &clientImpl{
		Key:                 key,
		Secret:              secret,
		NextQueryPageRegexp: re,
		HostURL:             hostUrl,
		HTTPClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}

	newClient.setHostURL()

	return newClient
}

func (c *clientImpl) setHostURL() {
	c.HostURL = c.HostURL + wpAPIPrefix
}

func (c *clientImpl) sendRequest(request *http.Request) (*http.Response, error) {

	request.Header.Set(contentTypeHeader, applicationJson)
	request.SetBasicAuth(c.Key, c.Secret)

	return c.HTTPClient.Do(request)
}

func (c *clientImpl) getURL(endpoint string, parameters url.Values) (string, error) {
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

func (c *clientImpl) get(endpoint string, parameters url.Values) (*http.Response, error) {

	if parameters == nil {
		parameters = make(map[string][]string)
	}

	parameters["consumer_key"] = []string{c.Key}
	parameters["consumer_secret"] = []string{c.Secret}

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

func (c *clientImpl) post(endpoint string, json string, cookies []*http.Cookie) (*http.Response, error) {

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

func (c *clientImpl) delete(endpoint string, body []byte) (*http.Response, error) {

	payload := bytes.NewReader(body)

	parameters := url.Values{"consumer_key": []string{c.Key}, "consumer_secret": []string{c.Secret}}

	endpoint, err := c.getURL(endpoint, parameters)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("DELETE", endpoint, payload)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(request)
}

type HeaderKey string

const (
	HeaderKeyTotalPages HeaderKey = "X-Wp-Totalpages"
	HeaderKeyTotal      HeaderKey = "X-Wp-Total"
	HeaderKeyLink       HeaderKey = "Link"
	HeaderKeySource     HeaderKey = "X-Wc-Webhook-Source"
)

type QueryParam string

const (
	QueryParamPage    QueryParam = "page"
	QueryParamPerPage QueryParam = "per_page"
	QueryParamOffset  QueryParam = "offset"
)
