package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/99designs/keyring"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	client := NewClient()

	if client.clientID == "" || client.clientSecret == "" {
		return fmt.Errorf("the Client ID and Client secret were not found in env vars")
	}

	ring, err := keyring.Open(keyring.Config{
		ServiceName: "monzo-access-token",
	})
	if err != nil {
		return err
	}

	i, err := ring.Get("tokens")
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			fmt.Println("tokens not found in keychain, starting auth flow")
			if err := oauth(client); err != nil {
				return fmt.Errorf("failed to authenticate: %w", err)
			}
			// join the tokens and save them to the keychain
			tokens := client.accessToken + "::" + client.refreshToken
			if err := ring.Set(keyring.Item{
				Key:         "tokens",
				Data:        []byte(tokens),
				Label:       "Monzo Access Token",
				Description: "Access and refresh tokens for the Monzo API",
			}); err != nil {
				return fmt.Errorf("failed to set tokens in keychain: %w", err)
			}
		} else {
			return err
		}
	} else {
		// split the tokens from i
		tokens := string(i.Data)
		tokenSlice := strings.Split(tokens, "::")
		if len(tokenSlice) != 2 {
			return fmt.Errorf("unexpected token format: %s", tokens)
		}
		client.accessToken = tokenSlice[0]
		client.refreshToken = tokenSlice[1]
	}

	err = pingTest(client)
	if err != nil {
		return err
	}

	accounts, err := monzoAccounts(client)
	if err != nil {
		return err
	}

	var ukRetailAccountID string
	for _, account := range *accounts {
		if account.Type == "uk_retail" {
			ukRetailAccountID = account.ID
			break
		}
	}

	b, err := balance(client, ukRetailAccountID)
	if err != nil {
		return err
	}

	if b.Balance%100 == 0 {
		return nil
	}

	rounded := (b.Balance / 100) * 100
	pennies := b.Balance - rounded
	if pennies == 0 {
		return nil
	}

	SavingsPotID, err := potIDByName(client, ukRetailAccountID, "Savings")
	if err != nil {
		return err
	}

	err = depositToPot(client, ukRetailAccountID, *SavingsPotID, pennies)
	if err != nil {
		return err
	}

	return nil
}
