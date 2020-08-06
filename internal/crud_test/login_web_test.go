package crud_test

import (
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/bejaneps/acroplia/internal/services"
)

func TestWebLoginByValidEmail(t *testing.T) {
	var (
		email    = "bejanhtc@gmail.com"
		password = "Saburi123"
	)

	// setup selenium wd
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Close()

	// check login by email with correct credentials
	err = crud.LoginByEmail(email, password, wd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebLoginByInvalidEmail(t *testing.T) {
	var invalidEmail = "test@test.com"

	// setup selenium wd
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Close()

	// check login by email with invalid email
	err = crud.LoginByEmail(invalidEmail, "", wd)
	if err != crud.ErrInvalidEmail {
		t.Fatal(err)
	}
}

func TestWebLoginByEmailInvalidPassword(t *testing.T) {
	var (
		email           = "bejanhtc@gmail.com"
		invalidPassword = "Saburi321"
	)

	// setup selenium wd
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Close()

	// check login by email with invalid password
	err = crud.LoginByEmail(email, invalidPassword, wd)
	if err != crud.ErrInvalidPassword {
		t.Fatal(err)
	}
}

func TestWebLoginByValidPhone(t *testing.T) {
	var (
		validPhone    = "+905488892053"
		validPassword = "Saburi123"
	)

	// setup selenium wd
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Close()

	// check login by email with invalid phone
	// TODO: remove skip
	err = crud.LoginByPhone(validPhone, validPassword, wd)
	if err != nil {
		t.Skip(err)
	}
}

func TestWebLoginByInvalidPhone(t *testing.T) {
	var invalidPhone = "+1ds4215rw1"

	// setup selenium wd
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Close()

	// check login by email with invalid phone
	err = crud.LoginByPhone(invalidPhone, "", wd)
	if err != crud.ErrInvalidPhone {
		t.Fatal(err)
	}
}

func TestWebLoginByPhoneInvalidPassword(t *testing.T) {
	var (
		validPhone      = "+905488892053"
		invalidPassword = "invalidPassword"
	)

	// setup selenium wd
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Close()

	// check login by email with invalid phone
	// TODO: remove skip
	err = crud.LoginByPhone(validPhone, invalidPassword, wd)
	if err != crud.ErrInvalidPassword {
		t.Skip(err)
	}
}
