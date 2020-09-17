package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
	Id  int    `json:"id"`
	Src string `json:"src"`
}

type ProductAttribute struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Variation bool     `json:"variation"`
	Options   []string `json:"options"`
}

type Store struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ShopName string `json:"shop_name"`
	Url      string `json:"url"`
}

type TaxStatus string

const TaxStatusTaxable = TaxStatus("taxable")

type Status string

const StatusPublish = Status("publish")

type Product struct {
	ID               int                `json:"id"`
	Name             string             `json:"name"`
	Slug             string             `json:"slug"`
	Permalink        string             `json:"permalink"`
	ProductType      ProductType        `json:"type"`
	Status           Status             `json:"status"`
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
	ManageStock      bool               `json:"manage_stock"`
	StockQuantity    int                `json:"stock_quantity"`
	StockStatus      StockStatus        `json:"stock_status"`
	Store            Store              `json:"store"`
	Dimensions       Dimensions         `json:"dimensions"`
	Weight           string             `json:"weight"`
}

type Dimensions struct {
	Length string `json:"length"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type StockStatus string

//stock_status	string	Controls the stock status of the product.
//Options: instock, outofstock, onbackorder. Default is instock.

const StockStatusInstock = StockStatus("instock")
const StockStatusOutofstock = StockStatus("outofstock")
const StockStatusOnbackorder = StockStatus("onbackorder")

type QueryProductsResponse struct {
	Products []Product
	NextPage string
}

func (c *Client) QueryProducts(params url.Values) (*QueryProductsResponse, error) {

	params["consumer_key"] = []string{c.Key}
	params["consumer_secret"] = []string{c.Secret}
	vendors := params["vendor"]

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
		return nil, NewError(err, http.StatusInternalServerError, err.Error())
	}

	if len(vendors) > 0 && len(vendors[0]) > 0 {
		products = filterByVendor(products, vendors[0])
	}

	response := &QueryProductsResponse{
		Products: products,
		NextPage: nextPage(res.Header, c.NextQueryPageRegexp),
	}
	return response, nil

}

func filterByVendor(products []Product, vendor string) []Product {
	filteredProducts := make([]Product, 0)
	for _, product := range products {
		if strings.ToLower(product.Store.ShopName) == strings.ToLower(vendor) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	return filteredProducts
}
