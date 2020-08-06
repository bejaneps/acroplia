package crud

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

const (
	peopleButton = `nav > div > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > a:nth-child(4)`
	peopleList   = `nav > div > div:nth-child(2) > div:nth-child(2) > div:nth-child(3) > div:nth-child(1) > div`
	messageArea  = `textarea[placeholder="Message"]`
	sendButton   = `button[title*="Send"]`
	topUser      = `nav > div > div:nth-child(2) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > div:nth-child(2)`
)

var (
	// ErrUserNotFound used when searched user isn't found
	ErrUserNotFound = errors.New("user not found")
)

// SendMessage to a user in Acroplia using web
func SendMessage(fullname string, username string, message string, wd selenium.WebDriver) error {
	if err := acropliaPre(wd); err != nil {
		return err
	}

	// find people button
	elem, err := wd.FindElement(selenium.ByCSSSelector, peopleButton)
	if err != nil {
		return errors.Wrap(err, "finding people button")
	}

	// click on people button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking on people button")
	}

	// wait until people page loads
	time.Sleep(5 * time.Second)

	// check if tips button is present
	elem, err = wd.FindElement(selenium.ByCSSSelector, tipsSkipButton)
	if err == nil {
		// click on it to proceed
		err = elem.Click()
		if err != nil {
			return errors.Wrap(err, "clicking on skip all button")
		}
	}

	// wait until skip all button removed
	time.Sleep(1 * time.Second)

	// check if specified user is in top chat
	topUserElem, err := wd.FindElement(selenium.ByCSSSelector, topUser)
	if err == nil {
		// check by fullname or username
		if fullname != "" {
			fullnameElem, err := topUserElem.FindElement(selenium.ByCSSSelector, `div:nth-child(1) > div:nth-child(1)`)
			if err != nil {
				return errors.Wrap(err, "finding user's fullname")
			}

			// get it's value
			text, err := fullnameElem.Text()
			if err != nil {
				return errors.Wrap(err, "getting user's fullname")
			}

			// compare it
			if text == fullname {
				if err := sendMessage(message, elem, wd); err != nil {
					return err
				}
				return nil
			}
		} else {
			usernameElem, err := topUserElem.FindElement(selenium.ByCSSSelector, `div:nth-child(2)`)
			if err != nil {
				return errors.Wrap(err, "finding user's username")
			}

			// get it's value
			text, err := usernameElem.Text()
			if err != nil {
				return errors.Wrap(err, "getting user's username")
			}

			// compare it
			if text[1:] == username {
				if err := sendMessage(message, elem, wd); err != nil {
					return err
				}
				return nil
			}
		}
	}

	// find people list
	elems, err := wd.FindElements(selenium.ByCSSSelector, peopleList)
	if err != nil {
		return errors.Wrap(err, "finding people list")
	}

	// iterate by each user
	var flag bool
	for _, elem := range elems {
		if fullname != "" {
			// find user's fullname
			fullnameElem, err := elem.FindElement(selenium.ByCSSSelector, `div:nth-child(2) > div:nth-child(1) > div:nth-child(1)`)
			if err != nil {
				return errors.Wrap(err, "finding user's fullname")
			}

			// get it's value
			text, err := fullnameElem.Text()
			if err != nil {
				return errors.Wrap(err, "getting user's fullname")
			}

			// compare it
			if text == fullname {
				flag = true
				if err := sendMessage(message, elem, wd); err != nil {
					return err
				}
				break
			}
		} else {
			// find user's username
			usernameElem, err := elem.FindElement(selenium.ByCSSSelector, `div:nth-child(2) > div:nth-child(2)`)
			if err != nil {
				return errors.Wrap(err, "finding user's username")
			}

			// get it's value
			text, err := usernameElem.Text()
			if err != nil {
				return errors.Wrap(err, "getting user's username")
			}

			// compare it (remove @)
			if text[1:] == username {
				flag = true
				if err := sendMessage(message, elem, wd); err != nil {
					return err
				}
				break
			}
		}
	}

	// check if user was found
	if !flag {
		return ErrUserNotFound
	}

	return nil
}

func sendMessage(message string, elem selenium.WebElement, wd selenium.WebDriver) error {
	// click on user
	err := elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking on user")
	}

	// find message area
	msgElem, err := wd.FindElement(selenium.ByCSSSelector, messageArea)
	if err != nil {
		return errors.Wrap(err, "finding message area")
	}

	// fill message
	err = msgElem.SendKeys(message)
	if err != nil {
		return errors.Wrap(err, "sending message")
	}

	// find send button
	msgButton, err := wd.FindElement(selenium.ByCSSSelector, sendButton)
	if err != nil {
		return errors.Wrap(err, "finding send message button")
	}

	// click on send button
	err = msgButton.Click()
	if err != nil {
		return errors.Wrap(err, "clicking on send message button")
	}

	// wait until message is sent
	time.Sleep(time.Second * 2)

	return nil
}
