package wooCommerce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type VariationAttribute struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Option string `json:"option"`
}

type ProductVariation struct {
	ID               int                  `json:"id"`
	Slug             string               `json:"slug"`
	Permalink        string               `json:"permalink"`
	ProductType      ProductType          `json:"type"`
	Status           Status               `json:"status"`
	Description      string               `json:"description"`
	ShortDescription string               `json:"short_description"`
	Sku              string               `json:"sku"`
	Price            string               `json:"price"`
	RegularPrice     string               `json:"regular_price"`
	SalePrice        string               `json:"sale_price"`
	Tags             []Tag                `json:"tags"`
	Image            Image                `json:"image"`
	Attributes       []VariationAttribute `json:"attributes"`
	TaxStatus        TaxStatus            `json:"tax_status"`
	ManageStock      bool                 `json:"manage_stock"`
	StockQuantity    int                  `json:"stock_quantity"`
	StockStatus      StockStatus          `json:"stock_status"`
	Store            Store                `json:"store"`
}

type QueryProductsVariationResponse struct {
	Variations []ProductVariation
	NextPage   string
}

const ProductVariationsEndpoint = "products/%d/variations"

func (c *Client) QueryProductVariations(productId int, params url.Values) (*QueryProductsVariationResponse, error) {

	params["consumer_key"] = []string{c.Key}
	params["consumer_secret"] = []string{c.Secret}

	res, err := c.Get(fmt.Sprintf(ProductVariationsEndpoint, productId), params)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError)
	}

	if res.StatusCode != http.StatusOK {
		var resErr Error
		err = json.Unmarshal(bodyBytes, &resErr)
		if err != nil {
			return nil, NewError(err, res.StatusCode)
		}
		return nil, NewError(resErr, res.StatusCode, resErr.Message)
	}

	var products []ProductVariation
	err = json.Unmarshal(bodyBytes, &products)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError, err.Error())
	}

	response := &QueryProductsVariationResponse{
		Variations: products,
		NextPage:   nextPage(res.Header, c.NextQueryPageRegexp),
	}
	return response, nil

}
