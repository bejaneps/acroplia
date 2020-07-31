package login

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/bejaneps/acroplia/internal/services"
	"github.com/jszwec/csvutil"
)

type UserEmail struct {
	Email    string `csv:"email"`
	Password string `csv:"password"`
}

type UserUsername struct {
	Username string `csv:"username"`
	Password string `csv:"password"`
}

type UserPhone struct {
	Phone    string `csv:"phone"`
	Password string `csv:"password"`
}

func TestAPILoginEmail(t *testing.T) {
	f, err := os.Open("./data/email.csv")
	if err != nil {
		t.Fatal(err)
	}

	dec, err := csvutil.NewDecoder(csv.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}

	// read users from csv file
	var users []UserEmail
	for {
		var user UserEmail
		if err := dec.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			t.Errorf("csv decoding: %v", err)
			continue
		}

		users = append(users, user)
	}

	// for each user perform a login
	for _, u := range users {
		s := crud.NewSignInEmailRequest(u.Email, u.Password)

		_, err := s.LoginByEmail()
		if errors.Is(err, crud.ErrInvalidEmail) {
			t.Logf("invalid email: %s", u.Email)
		} else if err != nil {
			t.Errorf("login %s: %v", u.Email, err)
		} else {
			t.Logf("valid email: %s", u.Email)
		}
	}
}

func TestAPILoginUsername(t *testing.T) {
	f, err := os.Open("./data/username.csv")
	if err != nil {
		t.Fatal(err)
	}

	dec, err := csvutil.NewDecoder(csv.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}

	// read users from csv file
	var users []UserUsername
	for {
		var user UserUsername
		if err := dec.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			t.Errorf("csv decoding: %v", err)
			continue
		}

		users = append(users, user)
	}

	// for each user perform a login
	for _, u := range users {
		s := crud.NewSignInUsernameRequest(u.Username, u.Password)

		_, err := s.LoginByUsername()
		if errors.Is(err, crud.ErrInvalidUsername) {
			t.Logf("invalid username: %s", u.Username)
		} else if err != nil {
			t.Errorf("login %s: %v", u.Username, err)
		} else {
			t.Logf("valid username: %s", u.Username)
		}
	}
}

func TestAPILoginPhone(t *testing.T) {
	f, err := os.Open("./data/phone.csv")
	if err != nil {
		t.Fatal(err)
	}

	dec, err := csvutil.NewDecoder(csv.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}

	// read users from csv file
	var users []UserPhone
	for {
		var user UserPhone
		if err := dec.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			t.Errorf("csv decoding: %v", err)
			continue
		}

		users = append(users, user)
	}

	// for each user perform a login
	for _, u := range users {
		s := crud.NewSignInPhoneRequest(u.Phone, u.Password)

		_, err := s.LoginByPhone()
		if errors.Is(err, crud.ErrInvalidPhone) {
			t.Logf("invalid phone: %s", u.Phone)
		} else if err != nil {
			t.Errorf("login %s: %v", u.Phone, err)
		} else {
			t.Logf("valid phone: %s", u.Phone)
		}
	}
}

func TestWebLoginEmail(t *testing.T) {
	f, err := os.Open("./data/email.csv")
	if err != nil {
		t.Fatal(err)
	}

	dec, err := csvutil.NewDecoder(csv.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}

	// read users from csv file
	var users []UserEmail
	for {
		var user UserEmail
		if err := dec.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			t.Errorf("csv decoding: %v", err)
			continue
		}

		users = append(users, user)
	}

	// for each user perform a login
	for _, u := range users {
		wd, err := services.NewSelenium("4444", "firefox")
		if err != nil {
			t.Fatal(err)
		}

		err = crud.LoginByEmail(u.Email, u.Password, wd)
		if errors.Is(err, crud.ErrInvalidEmail) {
			t.Logf("invalid email: %s", u.Email)
		} else if err != nil {
			t.Errorf("login %s: %v", u.Email, err)
		} else {
			t.Logf("valid email: %s", u.Email)
		}

		wd.Close()
	}
}

func TestWebLoginPhone(t *testing.T) {
	f, err := os.Open("./data/phone.csv")
	if err != nil {
		t.Fatal(err)
	}

	dec, err := csvutil.NewDecoder(csv.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}

	// read users from csv file
	var users []UserPhone
	for {
		var user UserPhone
		if err := dec.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			t.Errorf("csv decoding: %v", err)
			continue
		}

		users = append(users, user)
	}

	// for each user perform a login
	for _, u := range users {
		wd, err := services.NewSelenium("4444", "firefox")
		if err != nil {
			t.Fatal(err)
		}

		err = crud.LoginByPhone(u.Phone, u.Password, wd)
		if errors.Is(err, crud.ErrInvalidPhone) {
			t.Logf("invalid phone: %s", u.Phone)
		} else if err != nil {
			t.Errorf("login %s: %v", u.Phone, err)
		} else {
			t.Logf("valid phone: %s", u.Phone)
		}

		wd.Close()
	}
}
