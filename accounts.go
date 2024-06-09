package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountsResp struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID                string         `json:"id"`
	Closed            bool           `json:"closed"`
	Created           string         `json:"created"`
	Description       string         `json:"description"`
	Type              string         `json:"type"`
	OwnerType         string         `json:"owner_type"`
	IsFlex            bool           `json:"is_flex"`
	Currency          string         `json:"currency"`
	LegalEntity       string         `json:"legal_entity"`
	CountryCode       string         `json:"country_code"`
	CountryCodeAlpha3 string         `json:"country_code_alpha3"`
	Owners            []Owner        `json:"owners"`
	LinkedAccounts    []string       `json:"linked_accounts"`
	BusinessID        string         `json:"business_id"`
	AccountNumber     string         `json:"account_number"`
	SortCode          string         `json:"sort_code"`
	PaymentDetails    PaymentDetails `json:"payment_details"`
}

type Owner struct {
	UserID             string `json:"user_id"`
	PreferredName      string `json:"preferred_name"`
	PreferredFirstName string `json:"preferred_first_name"`
}

type PaymentDetails struct {
	UK   PaymentDetailsUK   `json:"locale_uk"`
	IBAN PaymentDetailsIBAN `json:"iban"`
}

type PaymentDetailsUK struct {
	AccountNumber string `json:"account_number"`
	SortCode      string `json:"sort_code"`
}

type PaymentDetailsIBAN struct {
	Unformatted string `json:"unformatted"`
	Formatted   string `json:"formatted"`
	BIC         string `json:"bic"`
}

func monzoAccounts(c *MonzoClient) (*[]Account, error) {
	path := "accounts"
	requestURL := fmt.Sprintf("%s/%s", c.endpoints["APIURL"], path)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

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

	var accountsResp AccountsResp

	err = json.NewDecoder(rsp.Body).Decode(&accountsResp)
	if err != nil {
		return nil, err
	}

	// for _, account := range accountsResp.Accounts {
	// 	if !account.Closed {
	// 		fmt.Printf("Account: %s - %s\n", account.Type, account.ID)
	// 	}
	// }

	var accounts []Account

	// If an account is Closed, it will not be returned
	for i, account := range accountsResp.Accounts {
		if !account.Closed {
			accounts = append(accounts, accountsResp.Accounts[i])
		}
	}

	return &accounts, nil
}
