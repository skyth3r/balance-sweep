package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func pingTest(c *MonzoClient) error {
	path := "ping/whoami"
	requestURL := fmt.Sprintf("%s/%s", c.endpoints["APIURL"], path)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	// if http status code is is 401, refresh the token
	if rsp.StatusCode == http.StatusUnauthorized {
		err = refreshToken(c)
		if err != nil {
			return err
		}
	}

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	if rsp.Body == nil {
		return fmt.Errorf("response body is empty")
	}

	rspJson := map[string]interface{}{}
	err = json.NewDecoder(rsp.Body).Decode(&rspJson)
	if err != nil {
		return err
	}

	_, ok := rspJson["user_id"].(string)
	if !ok {
		return fmt.Errorf("cannot find user ID in response")
	}

	//fmt.Println("Test API call successful 🎉")

	return nil
}
