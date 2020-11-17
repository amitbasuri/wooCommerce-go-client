package wooCommerce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type WebhookTopic string

const WebhookTopicOrderCreated = "order.created"
const WebhookTopicOrderUpdated = "order.updated"
const WebhookTopicOrderDeleted = "order.deleted"

const WebhookTopicProductCreated = "product.created"
const WebhookTopicProductUpdated = "product.updated"
const WebhookTopicProductDeleted = "product.deleted"

type Webhook struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name"`
	Topic       WebhookTopic `json:"topic"`
	DeliveryURL string       `json:"delivery_url"`
}

const WebhookEndpoint = "wc/v3/webhooks"

func (c *clientImpl) CreateWebhook(w *Webhook) error {

	params, err := json.Marshal(w)
	if err != nil {
		return err
	}

	res, err := c.post(WebhookEndpoint, string(params), nil)
	if err != nil {
		return NewError(err, http.StatusInternalServerError)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return NewError(err, http.StatusInternalServerError)
	}

	if res.StatusCode != http.StatusCreated {
		var resErr Error
		err = json.Unmarshal(bodyBytes, &resErr)
		if err != nil {
			return NewError(err, res.StatusCode)
		}
		return NewError(resErr, res.StatusCode, resErr.Message)
	}

	if err := json.Unmarshal(bodyBytes, w); err != nil {
		return NewError(err, http.StatusInternalServerError, "error unmarshalling webhook response")
	}

	return nil
}

func (c *clientImpl) DeleteWebhook(id int, force bool) error {

	body, _ := json.Marshal(map[string]interface{}{"force": force})

	res, err := c.delete(fmt.Sprintf("%s/%d", WebhookEndpoint, id), body)
	if err != nil {
		return NewError(err, http.StatusInternalServerError)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return NewError(err, http.StatusInternalServerError)
	}

	if res.StatusCode != http.StatusOK {
		var resErr Error
		err = json.Unmarshal(bodyBytes, &resErr)
		if err != nil {
			return NewError(err, res.StatusCode)
		}
		return NewError(resErr, res.StatusCode, resErr.Message)
	}

	return nil
}

func (c *clientImpl) QueryWebhooks(deliveryUrl *string) ([]Webhook, error) {

	allWebhooks := make([]Webhook, 0)

	// get all webhooks
	next := "1"
	for next != "" {

		res, nextPage, err := c.queryWebhooksPaginated(url.Values{
			string(QueryParamPage): []string{next},
		})

		if err != nil {
			return nil, err
		}

		if deliveryUrl != nil {
			for _, webhook := range res {
				if webhook.DeliveryURL == *deliveryUrl {
					allWebhooks = append(allWebhooks, webhook)
				}
			}
		}

		next = nextPage
	}

	return allWebhooks, nil
}

func (c *clientImpl) queryWebhooksPaginated(params url.Values) ([]Webhook, string, error) {

	res, err := c.get(WebhookEndpoint, params)
	if err != nil {
		return nil, "", NewError(err, http.StatusInternalServerError)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", NewError(err, http.StatusInternalServerError)
	}

	if res.StatusCode != http.StatusOK {
		var resErr Error
		err = json.Unmarshal(bodyBytes, &resErr)
		if err != nil {
			return nil, "", NewError(err, res.StatusCode)
		}
		return nil, "", NewError(resErr, res.StatusCode, resErr.Message)
	}

	var webhooks []Webhook
	err = json.Unmarshal(bodyBytes, &webhooks)
	if err != nil {
		return nil, "", NewError(err, http.StatusInternalServerError, err.Error())
	}

	return webhooks, nextPage(res.Header, c.NextQueryPageRegexp), nil
}
