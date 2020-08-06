package crud_test

import (
	"log"
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/bejaneps/acroplia/internal/services"
)

func TestWebSendMessageByFullname(t *testing.T) {
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		log.Fatal(err)
	}
	defer wd.Close()

	err = crud.LoginByEmail("bejanhtc@gmail.com", "Saburi123", wd)
	if err != nil {
		log.Fatal(err)
	}

	err = crud.SendMessage("Bezhan Mukhidinov", "", fmt.Sprintf("Test Message %s", uuid.New().String(), wd)
	if err != nil {
		log.Fatal(err)
	}
}

func TestWebSendMessageByUsername(t *testing.T) {
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		log.Fatal(err)
	}
	defer wd.Close()

	err = crud.LoginByEmail("bejanhtc@gmail.com", "Saburi123", wd)
	if err != nil {
		log.Fatal(err)
	}

	err = crud.SendMessage("", "bejaneps", fmt.Sprintf("Test Message %s", uuid.New().String(), wd)
	if err != nil {
		log.Fatal(err)
	}
}

func TestWebSendMessageToInvalidUser(t *testing.T) {
	wd, err := services.NewSelenium("4444", "chrome")
	if err != nil {
		log.Fatal(err)
	}
	defer wd.Close()

	err = crud.LoginByEmail("bejanhtc@gmail.com", "Saburi123", wd)
	if err != nil {
		log.Fatal(err)
	}

	err = crud.SendMessage("Jonn Doe", "", fmt.Sprintf("Test Message %s", uuid.New().String(), wd)
	if err != crud.ErrUserNotFound {
		log.Fatal(err)
	}
}
