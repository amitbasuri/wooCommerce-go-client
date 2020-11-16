package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SettingsKey string

const SettingsKeyWeightUnit = SettingsKey("wc/v3/settings/products/woocommerce_weight_unit")
const SettingsKeyDimensionUnit = SettingsKey("wc/v3/settings/products/woocommerce_dimension_unit")

type Setting struct {
	ID          string      `json:"id"`
	Label       string      `json:"label"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Default     string      `json:"default"`
	Options     interface{} `json:"options"`
	Tip         string      `json:"tip"`
	Value       string      `json:"value"`
	GroupID     string      `json:"group_id"`
}

func (c *clientImpl) GetSettings(key SettingsKey) (*Setting, error) {
	params := url.Values{}
	params["consumer_key"] = []string{c.Key}
	params["consumer_secret"] = []string{c.Secret}

	res, err := c.get(string(key), params)
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

	var setting Setting
	err = json.Unmarshal(bodyBytes, &setting)
	if err != nil {
		return nil, NewError(err, http.StatusInternalServerError, err.Error())
	}

	return &setting, nil
}
