package crud_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
)

func TestAPITextpadCreate(t *testing.T) {
	// perform login to get user's info
	s := crud.NewSignInEmailRequest("bejanhtc@gmail.com", "Saburi123")
	authResp, err := s.LoginByEmail()
	if err != nil {
		t.Fatal(err)
	}

	// create a textpad using user's info
	textpad := crud.NewTextpad("Test Title", "Test Subtitle", authResp.Data.User)
	textpadResp, err := textpad.Create(authResp.Data.AccessToken)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response: \n\n")
	b, _ := json.MarshalIndent(textpadResp, "", "  ")
	_, _ = os.Stdout.Write(b)
}
