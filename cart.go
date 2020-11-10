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

func (c *Client) addToCart(shippingCart ShippingCart) ([]*http.Cookie, error) {

	cookies := make([]*http.Cookie, 0)

	for _, cartLine := range shippingCart.CartLines {
		addToCartParam := addToCartParam{
			ProductId: cartLine.SourceId,
			Quantity:  strconv.Itoa(int(cartLine.Quantity)),
		}

		if len(cartLine.ProductVariantId) > 0 {
			variantID, err := strconv.Atoi(cartLine.ProductVariantId)
			if err != nil {
				return nil, err
			}

			addToCartParam.VariationID = variantID
		}

		params, err := json.Marshal(&addToCartParam)
		if err != nil {
			return nil, err
		}

		response, err := c.Post(addToCartEndpoint, string(params), nil)
		if err != nil {
			return nil, err
		}

		cookies = append(cookies, response.Cookies()...)
	}

	return cookies, nil
}
