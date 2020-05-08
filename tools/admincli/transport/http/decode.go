package http

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func decodeResponse(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return errors.Wrap(http.ErrBodyReadAfterClose, string(body))
	}

	return nil
}
