package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SystemStatus struct {
	Environment Environment    `json:"environment"`
	Settings    SystemSettings `json:"settings"`
}

type Environment struct {
	HomeUrl string `json:"home_url"`
	SiteUrl string `json:"site_url"`
}

type SystemSettings struct {
	ApiEnabled bool   `json:"api_enabled"`
	Currency   string `json:"currency"`
}

const SystemStatusEndpoint = "system_status"

func (c *Client) SystemStatus() (*SystemStatus, error) {
	params := url.Values{}
	params["consumer_key"] = []string{c.Key}
	params["consumer_secret"] = []string{c.Secret}

	res, err := c.Get(SystemStatusEndpoint, params)
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

	response := &SystemStatus{}
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError, err.Error())
	}

	return response, nil

}
