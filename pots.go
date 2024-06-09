package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Pots struct {
	Pots []Pot `json:"pots"`
}

type Pot struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Style    string `json:"style"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
	Deleted  bool   `json:"deleted"`
}

func listPots(c *MonzoClient, id string) (*[]Pot, error) {
	path := "pots"
	requestURL := fmt.Sprintf("%s/%s", c.endpoints["APIURL"], path)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	q := req.URL.Query()
	q.Add("current_account_id", id)
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

	var p Pots

	err = json.NewDecoder(rsp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	var pots []Pot

	for _, pot := range p.Pots {
		if !pot.Deleted {
			pots = append(pots, pot)
		}
	}

	return &pots, nil
}

func potIDByName(c *MonzoClient, id, name string) (*string, error) {
	pots, err := listPots(c, id)
	if err != nil {
		return nil, err
	}

	for _, pot := range *pots {
		if pot.Name == name {
			return &pot.ID, nil
		}
	}

	return nil, fmt.Errorf("pot with name %s not found", name)
}

func depositToPot(c *MonzoClient, accountID, potID string, amount int64) error {
	path := "pots/" + potID + "/deposit"
	requestURL := fmt.Sprintf("%s/%s", c.endpoints["APIURL"], path)

	dedupeID := fmt.Sprintf("dedupe_id_%d", time.Now().UnixNano())

	form := url.Values{}
	form.Set("source_account_id", accountID)
	form.Set("amount", fmt.Sprintf("%d", amount))
	form.Set("dedupe_id", dedupeID)
	formData := form.Encode()

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBufferString(formData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	if rsp.Body == nil {
		return fmt.Errorf("response body is empty")
	}

	return nil
}
