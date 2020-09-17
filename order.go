package wooCommerce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Order struct {
	ID                 int             `json:"id"`
	ParentID           int             `json:"parent_id"`
	Number             string          `json:"number"`
	OrderKey           string          `json:"order_key"`
	Status             OrderStatus     `json:"status,omitempty"`
	Total              string          `json:"total"`
	LineItems          []LineItem      `json:"line_items"`
	TransactionID      string          `json:"transaction_id"`
	BillingAddress     Address         `json:"billing"`
	ShippingAddress    ShippingAddress `json:"shipping"`
	Currency           string          `json:"currency,omitempty"`
	CreatedAtGMT       string          `json:"date_created_gmt"`
	DiscountTotal      string          `json:"discount_total"`
	ShippingTotal      string          `json:"shipping_total"`
	TotalTax           string          `json:"total_tax"`
	MetaDataList       []MetaData      `json:"meta_data"`
	PaymentMethod      string          `json:"payment_method"`
	PaymentMethodTitle string          `json:"payment_method_title"`
	SetPaid            bool            `json:"set_paid"`
	TaxLines           []TaxLines      `json:"tax_lines,omitempty"`
	ShippingLines      []ShippingLines `json:"shipping_lines,omitempty"`
}

type MetaData struct {
	Key   MetaDataKey `json:"key"`
	Value string      `json:"value"`
}

type MetaDataKey string

const MetaDataKeyCheckoutUrl = MetaDataKey("checkout_url")

type OrderStatus string

const OrderStatusCompleted = OrderStatus("completed")

type LineItem struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	ProductID   int           `json:"product_id"`
	VariationID int           `json:"variation_id"`
	Quantity    int           `json:"quantity"`
	TaxClass    string        `json:"tax_class"`
	Subtotal    string        `json:"subtotal"`
	SubtotalTax string        `json:"subtotal_tax"`
	Total       string        `json:"total"`
	TotalTax    string        `json:"total_tax"`
	Taxes       []interface{} `json:"taxes"`
	MetaData    []interface{} `json:"meta_data"`
	Sku         *string       `json:"sku"`
	Price       float64       `json:"price"`
}

type Address struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
}

type ShippingAddress struct {
	Address
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type TaxLines struct {
	Compound         bool          `json:"compound"`
	ID               int64         `json:"id"`
	RateID           int64         `json:"rate_id"`
	RateCode         string        `json:"rate_code"`
	Label            string        `json:"label"`
	TaxTotal         string        `json:"tax_total"`
	ShippingTaxTotal string        `json:"shipping_tax_total"`
	MetaData         []interface{} `json:"meta_data"`
}

type ShippingLines struct {
	ID          int64         `json:"id"`
	MethodTitle string        `json:"method_title"`
	MethodID    string        `json:"method_id"`
	Total       string        `json:"total"`
	TotalTax    string        `json:"total_tax"`
	Taxes       []interface{} `json:"taxes"`
	MetaData    []interface{} `json:"meta_data"`
}

const OrdersEndpoint = "orders"

func (c *Client) GetOrder(id int) (*Order, error) {
	params := url.Values{}
	params["consumer_key"] = []string{c.Key}
	params["consumer_secret"] = []string{c.Secret}

	res, err := c.Get(fmt.Sprintf("%s/%d", OrdersEndpoint, id), params)
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

	var order Order
	err = json.Unmarshal(bodyBytes, &order)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError, err.Error())
	}

	return &order, nil
}

func (c *Client) CreateOrder(o *Order) (*Order, error) {

	var orderResponse Order

	params, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	res, err := c.Post(OrdersEndpoint, string(params))
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError)
	}

	fmt.Println(res.Status)

	if res.StatusCode != http.StatusCreated {
		var resErr Error
		err = json.Unmarshal(bodyBytes, &resErr)
		if err != nil {
			return nil, NewError(err, res.StatusCode)
		}
		return nil, NewError(resErr, res.StatusCode, resErr.Message)
	}

	if err := json.Unmarshal(bodyBytes, &orderResponse); err != nil {
		return nil, NewError(err, http.StatusInternalServerError, "error unmarshalling order response")
	}

	return &orderResponse, nil
}
