package wooCommerce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)



func (c *Client) addToCart() {
	var orderResponse Order

	params, err := json.Marshal(orderResponse)
	if err != nil {

	}

	res, err := c.Post(OrdersEndpoint, string(params))
	if err != nil {

	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {

	}

	if res.StatusCode != http.StatusCreated {
		var resErr Error
		err = json.Unmarshal(bodyBytes, &resErr)
		if err != nil {

		}

	}

	if err := json.Unmarshal(bodyBytes, &orderResponse); err != nil {

	}


}
