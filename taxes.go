package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const TaxRatesEndPoint = "wc/v3/taxes"

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

type getTaxesResponse struct {
	Taxes    []Tax
	NextPage string
}

func (c *clientImpl) GetTaxes(params url.Values) ([]Tax, error) {

	allTaxes := make([]Tax, 0)

	// get all taxes
	next := "1"
	for next != "" {

		res, err := c.GetTaxesPaginated(url.Values{
			string(QueryParamPage): []string{next},
		})
		if err != nil {
			break
		}
		allTaxes = append(allTaxes, res.Taxes...)
		next = res.NextPage
	}

	countryCode := params["countryCode"][0]
	state := params["state"][0]

	taxes := filterTaxesByCountryAndState(allTaxes, countryCode, state)

	return taxes, nil
}

func (c *clientImpl) GetTaxesPaginated(params url.Values) (*getTaxesResponse, error) {

	params["consumer_key"] = []string{c.Key}
	params["consumer_secret"] = []string{c.Secret}

	res, err := c.get(TaxRatesEndPoint, params)
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

	taxes := make([]Tax, 0)
	err = json.Unmarshal(bodyBytes, &taxes)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError, err.Error())
	}

	response := &getTaxesResponse{
		Taxes:    taxes,
		NextPage: nextPage(res.Header, c.NextQueryPageRegexp),
	}
	return response, nil
}

func filterTaxesByCountryAndState(taxesList []Tax, countryCode, state string) []Tax {

	taxes := make([]Tax, 0)

	for _, tax := range taxesList {
		if tax.Country == countryCode && tax.State == state {
			taxes = append(taxes, tax)
		}
	}

	return taxes
}
