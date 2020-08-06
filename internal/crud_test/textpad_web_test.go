package crud_test

import (
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/bejaneps/acroplia/internal/services"
)

func TestWebTextpadCreate(t *testing.T) {
	// start selenium instance
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Close()

	// perform a login
	err = crud.LoginByEmail("bejanhtc@gmail.com", "Saburi123", wd)
	if err != nil {
		t.Fatal(err)
	}

	// create a textpad
	err = crud.TextpadCreate("Title for Test", "", wd)
	if err != nil {
		t.Fatal(err)
	}
}
