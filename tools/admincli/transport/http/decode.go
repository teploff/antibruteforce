package http

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/teploff/antibruteforce/internal/shared"
)

func decodeResponse(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return errors.Wrap(shared.ErrEmpty, string(body))
	}

	return nil
}
