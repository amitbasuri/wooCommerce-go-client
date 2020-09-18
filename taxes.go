package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const TaxRatesEndPoint = "taxes"

type Tax struct {
	Id         int    `json:"id"`
	Country    string `json:"country"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	City       string `json:"city"`
	Rate       string `json:"rate"`
	Name       string `json:"name"`
	Priority   int    `json:"priority"`
	Compound   bool   `json:"compound"`
	Shipping   bool   `json:"shipping"`
	Class      string `json:"class"`
}

func (c *Client) GetTaxes(params url.Values) ([]Tax, error) {

	params["consumer_key"] = []string{c.Key}
	params["consumer_secret"] = []string{c.Secret}

	res, err := c.Get(TaxRatesEndPoint, params)
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

	allTaxes := make([]Tax, 0)
	err = json.Unmarshal(bodyBytes, &allTaxes)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError, err.Error())
	}

	countryCode := params["countryCode"][0]
	state := params["state"][0]

	taxes := make([]Tax, 0)

	for _, tax := range allTaxes {
		if tax.Country == countryCode && tax.State == state {
			taxes = append(taxes, tax)
		}
	}

	return taxes, nil
}
