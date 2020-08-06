package crud_test

import (
	"fmt"
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/google/uuid"
)

func TestAPISendMessage(t *testing.T) {
	// perform login through api
	req := crud.NewSignInEmailRequest("bejanhtc@gmail.com", "Saburi123")
	resp, err := req.LoginByEmail()
	if err != nil {
		t.Fatal(err)
	}

	// send message
	msg := crud.NewMessage(fmt.Sprintf("Test Message %s", uuid.New().String()), resp.Data.User)
	respMsg, err := msg.SendMessage("23e7235e-2072-3330-91b7-3b747d5a87e7", resp.Data.AccessToken)
	if err == crud.ErrMessageInternal {
		t.Skip("server sends 500 Internal Server Error, but message is sent")
	} else if err != nil {
		t.Fatal(err)
	}

	_ = respMsg
}
