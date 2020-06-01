package wooCommerce

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

var testClient = NewClient("https://shoptypewoo.wpcomstaging.com/",
	"ck_0df9f84a48f84e0447e546b2fca6a38a60f2edc2",
	"cs_dcb9d713e2a695e320bf3c8195e6db12dc82dfd8")

//curl https://shoptypewoo.wpcomstaging.com/wp-json/wc/v3/products/265 \
//-u consumer_key:consumer_secret
func TestClient_QueryProducts(t *testing.T) {

	next := "1"

	for {
		if next == "" || next == "3" {
			break
		}

		p, err := testClient.QueryProducts(url.Values{string(QueryParamPage): []string{next}})
		assert.NoError(t, err)
		if err != nil {
			break
		}
		next = p.NextPage
	}

}

func TestClient_QueryProductVariations(t *testing.T) {

	next := "1"

	for {
		if next == "" || next == "3" {
			break
		}

		p, err := testClient.QueryProductVariations(37, url.Values{string(QueryParamPage): []string{next}})
		assert.NoError(t, err)
		if err != nil {
			break
		}
		next = p.NextPage
	}

}
