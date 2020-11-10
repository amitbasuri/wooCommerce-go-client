package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Shipping struct {
	ProductVariantId string            `json:"product_variant_id"`
	SourceId         string            `json:"source_id"`
	Quantity         uint              `json:"quantity"`
	Name             string            `json:"name"`
	Price            float64           `json:"price"`
	VariantNameValue map[string]string `json:"variant_name_value"`
	CountryCode      string            `json:"country_code"`
}

type CalculateShippingParam struct {
	Country       string `json:"country"`
	ReturnMethods bool   `json:"return_methods"`
}

type ShippingResponse struct {
	Cost string `json:"cost"`
	Html string `json:"html"`
}

const shippingEndpoint = "calculate/shipping"

func (c *Client) CalculateShipping(shipping Shipping) (*ShippingResponse, error) {

	calculateShippingParam := CalculateShippingParam{
		Country:       shipping.CountryCode,
		ReturnMethods: true,
	}

	params, err := json.Marshal(calculateShippingParam)
	if err != nil {
		return nil, err
	}


	response, err := c.Post(shippingEndpoint, string(params), nil)
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

	responseMap := make(map[string]interface{})
	if err := json.Unmarshal(bodyBytes, &responseMap); err != nil {
		return nil, NewError(err, http.StatusInternalServerError, "error unmarshalling order response")
	}

	if val, ok := responseMap["cost"]; ok {
		shippingResponse.Cost = val.(string)
	}

	if val, ok := responseMap["html"]; ok {
		shippingResponse.Cost = val.(string)
	}

	return shippingResponse, nil
}
