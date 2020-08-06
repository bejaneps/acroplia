package crud_test

import (
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
)

func TestCheckInternetConnection(t *testing.T) {
	if err := crud.CheckInternetConnection(); err != nil {
		t.Fatal(err)
	}
}
