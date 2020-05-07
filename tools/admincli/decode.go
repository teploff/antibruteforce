package main

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func decodeResponse(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return errors.New(string(body))
	}

	return nil
}
