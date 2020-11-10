package wooCommerce

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

var testClient = NewClient("https://shoptypewoo.wpcomstaging.com/",
	"ck_b22be12d33b3bee1365fb2776aaff11d6c9d7c9a",
	"cs_11d03e4028aaec811ef45dd1b246250e030fb517")

//var testClient = NewClient("https://www.adamscbd.com/",
//	"ck_8ba6fba964c8883cfa4deb2a80cb670ac7ad1cc8",
//	"cs_4c6ac24e46437855b8b5eb99118a3b27ff19f61f")

//curl https://shoptypewoo.wpcomstaging.com/wp-json/wc/v3/orders/345 \
//-u ck_0df9f84a48f84e0447e546b2fca6a38a60f2edc2:cs_dcb9d713e2a695e320bf3c8195e6db12dc82dfd8 -d samples/order_upadted.json
//
//curl -X PUT https://shoptypewoo.wpcomstaging.com/wp-json/wc/v3/orders/345 \
//-u "ck_b22be12d33b3bee1365fb2776aaff11d6c9d7c9a:cs_11d03e4028aaec811ef45dd1b246250e030fb517" \
//-H "Content-Type: application/json" \
//-d '{
//"transaction_id": "12345"
//}'

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

func TestClient_SystemStatus(t *testing.T) {
	s, err := testClient.SystemStatus()
	assert.NoError(t, err)
	assert.Equal(t, "https://shoptypewoo.wpcomstaging.com", s.Environment.SiteUrl)
}

func TestClient_GetOrder(t *testing.T) {
	s, err := testClient.GetOrder(543)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	//fmt.Printf("%+v", s)
	//assert.Equal(t, "https://shoptypewoo.wpcomstaging.com", s.Environment.SiteUrl)
}

//curl  https://adamscbd.com/wp-json/wc/v3/products?consumer_key=ck_8ba6fba964c8883cfa4deb2a80cb670ac7ad1cc8&consumer_secret=cs_4c6ac24e46437855b8b5eb99118a3b27ff19f61f
//
//
//
//-H "Content-Type: application/json"
//
//-u "ck_8ba6fba964c8883cfa4deb2a80cb670ac7ad1cc8"
//
//
//curl 'https://shoptypewoo.wpcomstaging.com/wp-json/wc/v3/orders/345?consumer_key=ck_0df9f84a48f84e0447e546b2fca6a38a60f2edc2&consumer_secret=cs_dcb9d713e2a695e320bf3c8195e6db12dc82dfd8'
//
//
//-d samples/order_upadted.json

//https://www.adamscbd.com/wp-json/wc/v3/products?
//// consumer_key=ck_8ba6fba964c8883cfa4deb2a80cb670ac7ad1cc8&
//// consumer_secret=cs_4c6ac24e46437855b8b5eb99118a3b27ff19f61f

func TestClient_GetSettings(t *testing.T) {

	p, err := testClient.GetSettings(SettingsKeyWeightUnit)
	assert.NoError(t, err)
	//for _,v := range p {
	//	fmt.Printf("+%v\n", v)
	//}
	fmt.Printf("+%v\n", p)
}

func TestClient_CalculateShipping(t *testing.T) {
	shipping := Shipping{
		ProductVariantId: "",
		SourceId:         "614",
		Quantity:         1,
		CountryCode:      "US",
	}

	_, err := testClient.CalculateShipping(shipping)

	assert.NoError(t, err)
}

func TestClient_CreateOrder(t *testing.T) {

	order := Order{}
	err := json.Unmarshal([]byte("{ \"payment_method\": \"bacs\",\n        \"payment_method_title\": \"Direct Bank Transfer\",\n        \"set_paid\": true,\n        \"billing\": {\n            \"first_name\": \"John\",\n            \"last_name\": \"Doe\",\n            \"address_1\": \"969 Market\",\n            \"address_2\": \"\",\n            \"city\": \"San Francisco\",\n            \"state\": \"CA\",\n            \"postcode\": \"94103\",\n            \"country\": \"US\",\n            \"email\": \"john.doe@example.com\",\n            \"phone\": \"(555) 555-5555\"\n        },\n        \"shipping\": {\n            \"first_name\": \"John\",\n            \"last_name\": \"Doe\",\n            \"address_1\": \"969 Market\",\n            \"address_2\": \"\",\n            \"city\": \"San Francisco\",\n            \"state\": \"CA\",\n            \"postcode\": \"94103\",\n            \"country\": \"US\"\n        },\n        \"line_items\": [\n            {\n                \"product_id\": 287,\n                \"quantity\": 2\n            }\n        ]}"), &order)


	_, err = testClient.CreateOrder(&order)

	assert.NoError(t, err)
}
