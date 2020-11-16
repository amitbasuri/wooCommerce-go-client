package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ShippingCart struct {
	CountryCode string     `json:"country_code"`
	CartLines   []CartLine `json:"cart_lines"`
}

type CartLine struct {
	ProductVariantId string            `json:"product_variant_id"`
	SourceId         string            `json:"source_id"`
	Quantity         uint              `json:"quantity"`
	Name             string            `json:"name"`
	Price            float64           `json:"price"`
	VariantNameValue map[string]string `json:"variant_name_value"`
}

type CalculateShippingParam struct {
	Country       string `json:"country"`
	ReturnMethods bool   `json:"return_methods"`
}

type ShippingResponse struct {
	Key          string        `json:"key"`
	MethodID     string        `json:"method_id"`
	InstanceID   int           `json:"instance_id"`
	Label        string        `json:"label"`
	Cost         string        `json:"cost"`
	HTML         string        `json:"html"`
	Taxes        []interface{} `json:"taxes"`
	ChosenMethod bool          `json:"chosen_method"`
}

const shippingEndpoint = "cocart/v1/calculate/shipping"

func (c *Client) CalculateShipping(shippingCart ShippingCart) (*ShippingResponse, error) {

	calculateShippingParam := CalculateShippingParam{
		Country:       shippingCart.CountryCode,
		ReturnMethods: true,
	}

	params, err := json.Marshal(calculateShippingParam)
	if err != nil {
		return nil, err
	}

	cookies, err := c.addToCart(shippingCart)
	if err != nil {
		return nil, err
	}

	response, err := c.Post(shippingEndpoint, string(params), cookies)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError)
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError)
	}

	if response.StatusCode != http.StatusOK {
		var resErr Error
		err = json.Unmarshal(bodyBytes, &resErr)
		if err != nil {
			return nil, NewError(err, response.StatusCode)
		}
		return nil, NewError(resErr, response.StatusCode, resErr.Message)
	}

	shippingResponse := &ShippingResponse{}

	responseMap := make(map[string]*ShippingResponse)
	if err := json.Unmarshal(bodyBytes, &responseMap); err != nil {
		return nil, NewError(err, http.StatusInternalServerError, "error unmarshalling order response")
	}

	for _, respMap := range responseMap {
		shippingResponse = respMap
	}

	return shippingResponse, nil
}
