package crud

import (
	"net/http"

	"github.com/pkg/errors"
)

// CheckInternetConnection is a utility function that checks if user has internet connection, if there is no connection then tool won't run.
func CheckInternetConnection() error {
	resp, err := http.Get(acropliaLoginURL)
	if err != nil {
		return errors.Wrapf(err, "couldn't navigate to %s", acropliaLoginURL)
	}

	if resp.StatusCode != 200 {
		return errors.Errorf("no internet connection")
	}

	return nil
}
