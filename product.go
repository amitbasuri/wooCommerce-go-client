package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const ProductsEndpoint = "products"

type ProductType string

const (
	ProductTypeSimple   = ProductType("simple")
	ProductTypeVariable = ProductType("variable")
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Image struct {
	Id  string `json:"id"`
	Src string `json:"src"`
}

type ProductAttribute struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Variation bool     `json:"variation"`
	Options   []string `json:"options"`
}

type TaxStatus string

const TaxStatusTaxable = TaxStatus("taxable")

type Product struct {
	ID               int                `json:"id"`
	Name             string             `json:"name"`
	Slug             string             `json:"slug"`
	Permalink        string             `json:"permalink"`
	ProductType      ProductType        `json:"type"`
	Status           string             `json:"status"`
	Description      string             `json:"description"`
	ShortDescription string             `json:"short_description"`
	Sku              string             `json:"sku"`
	Price            string             `json:"price"`
	RegularPrice     string             `json:"regular_price"`
	SalePrice        string             `json:"sale_price"`
	Tags             []Tag              `json:"tags"`
	Images           []Image            `json:"images"`
	Attributes       []ProductAttribute `json:"attributes"`
	TaxStatus        TaxStatus          `json:"tax_status"`
}

type QueryProductsResponse struct {
	Products []Product
	NextPage string
}

func (c *Client) QueryProducts(params url.Values) (*QueryProductsResponse, error) {

	res, err := c.Get(ProductsEndpoint, params)
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

	var products []Product
	err = json.Unmarshal(bodyBytes, &products)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError, "error unmarshal response")
	}

	response := &QueryProductsResponse{
		Products: products,
		NextPage: nextPage(res.Header, c.NextQueryPageRegexp),
	}
	return response, nil

}
