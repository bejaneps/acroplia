package crud_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
)

func TestAPILoginByValidEmail(t *testing.T) {
	// create bulk request body
	s := crud.NewSignInEmailRequest("bejanhtc@gmail.com", "Saburi123")

	// make request to server
	resp, err := s.LoginByEmail()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response: \n\n")
	b, _ := json.MarshalIndent(resp, "", "  ")
	_, _ = os.Stdout.Write(b)

	// compare results from response to expected results
	expFirstName := "Bezhan"
	expLastName := "Mukhidinov"
	if resp.Data.User.FirstName != expFirstName || resp.Data.User.LastName != expLastName {
		t.Fatalf("expected firstName: %s, lastName: %s; got firstName %s, lastName %s", expFirstName, expLastName, resp.Data.User.FirstName, resp.Data.User.LastName)
	}
}

func TestAPILoginByInvalidEmail(t *testing.T) {
	// create bulk request body
	s := crud.NewSignInEmailRequest("invalidEmail@invalid.com", "invalidPassword")

	// make request to server
	_, err := s.LoginByEmail()
	if err != crud.ErrInvalidEmail {
		t.Fatal(err)
	}
}

func TestAPILoginByValidUsername(t *testing.T) {
	// create bulk request body
	s := crud.NewSignInUsernameRequest("bejanhtc", "Saburi123")

	// make request to server
	resp, err := s.LoginByUsername()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response: \n\n")
	b, _ := json.MarshalIndent(resp, "", "  ")
	_, _ = os.Stdout.Write(b)

	// compare results from response to expected results
	expFirstName := "Bezhan"
	expLastName := "Mukhidinov"
	if resp.Data.User.FirstName != expFirstName || resp.Data.User.LastName != expLastName {
		t.Fatalf("expected firstName: %s, lastName: %s; got firstName %s, lastName %s", expFirstName, expLastName, resp.Data.User.FirstName, resp.Data.User.LastName)
	}
}

func TestAPILoginByInvalidUsername(t *testing.T) {
	// create bulk request body
	s := crud.NewSignInUsernameRequest("invalidUsername", "invalidPassword")

	// make request to server
	_, err := s.LoginByUsername()
	if err != crud.ErrInvalidUsername {
		t.Fatal(err)
	}
}

func TestAPILoginByValidPhone(t *testing.T) {
	// create bulk request body
	s := crud.NewSignInPhoneRequest("+905488892053", "Saburi123")

	// make request to server
	// TODO: remove skip
	resp, err := s.LoginByPhone()
	if err != nil {
		t.Skip(err)
	}

	t.Logf("Response: \n\n")
	b, _ := json.MarshalIndent(resp, "", "  ")
	_, _ = os.Stdout.Write(b)

	// compare results from response to expected results
	expFirstName := "Manal"
	expLastName := "Khalid"
	if resp.Data.User.FirstName != expFirstName || resp.Data.User.LastName != expLastName {
		t.Fatalf("expected firstName: %s, lastName: %s; got firstName %s, lastName %s", expFirstName, expLastName, resp.Data.User.FirstName, resp.Data.User.LastName)
	}
}

func TestAPILoginByInvalidPhone(t *testing.T) {
	// create bulk request body
	s := crud.NewSignInPhoneRequest("+123456", "invalidPassword")

	// make request to server
	_, err := s.LoginByPhone()
	if err != crud.ErrInvalidPhone {
		t.Fatal(err)
	}
}
