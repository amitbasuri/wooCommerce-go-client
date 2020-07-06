package wooCommerce

import (
	"github.com/tgglv/wc-api-go/client"
	"github.com/tgglv/wc-api-go/options"
	"regexp"
)

type Client struct {
	client.Client
	Key                 string
	Secret              string
	NextQueryPageRegexp *regexp.Regexp
}

func NewClient(hostUrl string, key string, secret string) *Client {
	factory := client.Factory{}
	c := factory.NewClient(options.Basic{
		URL:    hostUrl,
		Key:    key,
		Secret: secret,
		Options: options.Advanced{
			WPAPI:       true,
			WPAPIPrefix: "/wp-json/",
			Version:     "wc/v3",
		},
	})

	re := regexp.MustCompile(`\<(.*)\>;.(rel="next")`)

	return &Client{
		Client:              c,
		Key:                 key,
		Secret:              secret,
		NextQueryPageRegexp: re,
	}

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
