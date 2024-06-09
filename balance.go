package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Balance struct {
	Balance      int64  `json:"balance"`
	TotalBalance int64  `json:"total_balance"`
	Currency     string `json:"currency"`
	SpendToday   int64  `json:"spend_today"`
}

func balance(c *MonzoClient, id string) (*Balance, error) {
	path := "balance"
	requestURL := fmt.Sprintf("%s/%s", c.endpoints["APIURL"], path)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	q := req.URL.Query()
	q.Add("account_id", id)
	req.URL.RawQuery = q.Encode()

	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	if rsp.Body == nil {
		return nil, fmt.Errorf("response body is empty")
	}

	var b Balance

	err = json.NewDecoder(rsp.Body).Decode(&b)
	if err != nil {
		return nil, err
	}

	return &b, nil
}
