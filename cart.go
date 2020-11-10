package wooCommerce

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const addToCartEndpoint = "cocart/v1/add-item"

type addToCartParam struct {
	ProductId   string `json:"product_id"`
	Quantity    string `json:"quantity"`
	VariationID int    `json:"variation_id"`
}

func (c *Client) addToCart(shipping Shipping) (*http.Response, error) {

	addToCartParam := addToCartParam{
		ProductId:   shipping.SourceId,
		Quantity:    strconv.Itoa(int(shipping.Quantity)),
	}

	if len(shipping.ProductVariantId) > 0 {
		variantID, err := strconv.Atoi(shipping.ProductVariantId)
		if err != nil {
			return nil, err
		}

		addToCartParam.VariationID = variantID
	}

	params, err := json.Marshal(&addToCartParam)
	if err != nil {
		return nil, err
	}

	return c.Post(addToCartEndpoint, string(params), nil)
}
