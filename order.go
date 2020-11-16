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
	Key   MetaDataKey `json:"key,omitempty"`
	Value string      `json:"value,omitempty"`
}

type MetaDataKey string

const MetaDataKeyCheckoutUrl = MetaDataKey("checkout_url")

type OrderStatus string

const OrderStatusCompleted = OrderStatus("completed")
const OrderStatusOnHold = OrderStatus("on-hold")

type LineItem struct {
	ID          int           `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	ProductID   int           `json:"product_id,omitempty"`
	VariationID int           `json:"variation_id,omitempty"`
	Quantity    int           `json:"quantity,omitempty"`
	TaxClass    string        `json:"tax_class,omitempty"`
	Subtotal    string        `json:"subtotal,omitempty"`
	SubtotalTax string        `json:"subtotal_tax,omitempty"`
	Total       string        `json:"total,omitempty"`
	TotalTax    string        `json:"total_tax,omitempty"`
	Taxes       []interface{} `json:"taxes,omitempty"`
	MetaData    []interface{} `json:"meta_data,omitempty"`
	Sku         *string       `json:"sku,omitempty"`
	Price       float64       `json:"price,omitempty"`
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
	Compound         bool          `json:"compound,omitempty"`
	ID               int64         `json:"id,omitempty"`
	RateID           int64         `json:"rate_id,omitempty"`
	RateCode         string        `json:"rate_code,omitempty"`
	Label            string        `json:"label,omitempty"`
	TaxTotal         string        `json:"tax_total,omitempty"`
	ShippingTaxTotal string        `json:"shipping_tax_total,omitempty"`
	MetaData         []interface{} `json:"meta_data,omitempty"`
}

type ShippingLines struct {
	ID          int64         `json:"id,omitempty"`
	MethodTitle string        `json:"method_title,omitempty"`
	MethodID    string        `json:"method_id,omitempty"`
	Total       string        `json:"total,omitempty"`
	TotalTax    string        `json:"total_tax,omitempty"`
	Taxes       []interface{} `json:"taxes,omitempty"`
	MetaData    []interface{} `json:"meta_data,omitempty"`
}

const OrdersEndpoint = "wc/v3/orders"

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

	res, err := c.Post(OrdersEndpoint, string(params), nil)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError)
	}

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
