package wooCommerce

type Order struct {
	ID              int             `json:"id"`
	ParentID        int             `json:"parent_id"`
	Number          string          `json:"number"`
	OrderKey        string          `json:"order_key"`
	Status          OrderStatus     `json:"status"`
	Total           string          `json:"total"`
	LineItems       []LineItem      `json:"line_items"`
	TransactionID   string          `json:"transaction_id"`
	BillingAddress  Address         `json:"billing"`
	ShippingAddress ShippingAddress `json:"shipping"`
	Currency        string          `json:"currency"`
	CreatedAtGMT    string          `json:"date_created_gmt"`
	DiscountTotal   string          `json:"discount_total"`
	ShippingTotal   string          `json:"shipping_total"`
	TotalTax        string          `json:"total_tax"`
	MetaDataList    []MetaData      `json:"meta_data"`
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
	ID          int    `json:"id"`
	ProductID   int    `json:"product_id"`
	VariationID int    `json:"variation_id"`
	SubTotal    string `json:"subtotal"`
	Total       string `json:"total"`
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
