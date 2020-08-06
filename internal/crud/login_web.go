package crud

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

const acropliaLoginURL = "https://demo.acroplia.com/login"

const (
	emailInput    = `input[name*="email"]`
	passwordInput = `input[name*="password"]`
	phoneInput    = `input[name*="phone"]`

	submitButton   = `div[class*="btn-main"]`
	continueButton = `btn-nofiled-blue`
)

var (
	// ErrInvalidEmail is used when user submitted email is wrong when performing login
	ErrInvalidEmail = errors.New("invalid email")

	// ErrInvalidPassword is used when user submitted password is wrong when performing login
	ErrInvalidPassword = errors.New("invalid password")

	// ErrInvalidPhone is used when user submitted phone is wrong when performing login
	ErrInvalidPhone = errors.New("invalid phone")

	// ErrInvalidUsername is used when user submitted username is wrong when performing login
	ErrInvalidUsername = errors.New("invalid username")
)

// LoginByEmail performs a login to Acroplia using email
func LoginByEmail(email, password string, wd selenium.WebDriver) error {
	var err error

	// navigate to login page
	err = wd.Get(acropliaLoginURL)
	if err != nil {
		return errors.Wrapf(err, "navigating to %s", acropliaLoginURL)
	}

	// check if old version of browser is shown
	elem, err := wd.FindElement(selenium.ByClassName, continueButton)
	if err == nil {
		err = elem.Click() // click on continue button to proceed to login page
		if err != nil {
			return errors.Wrap(err, "clicking on continue anyway button")
		}
	}

	// wait until login page loads
	time.Sleep(1 * time.Second)

	// find email input
	elem, err = wd.FindElement(selenium.ByCSSSelector, emailInput)
	if err != nil {
		return errors.Wrap(err, "finding email input")
	}

	// fill email input
	err = elem.SendKeys(email)
	if err != nil {
		return errors.Wrap(err, "filling email input")
	}

	// find next button
	elem, err = wd.FindElement(selenium.ByCSSSelector, submitButton)
	if err != nil {
		return errors.Wrap(err, "finding submit button")
	}

	// click on next button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking submit button")
	}

	// wait until password page loads
	time.Sleep(2 * time.Second)

	// find password input
	elem, err = wd.FindElement(selenium.ByCSSSelector, passwordInput)
	if err != nil {
		return errors.Wrap(err, "finding password input")
	}

	// check if password input is visible, if no then email was wrong
	attr, _ := elem.GetAttribute("aria-hidden")
	if attr == "true" {
		return ErrInvalidEmail
	}

	// fill password input
	err = elem.SendKeys(password)
	if err != nil {
		return errors.Wrap(err, "filling password input")
	}

	// find sign in button
	elem, err = wd.FindElement(selenium.ByCSSSelector, submitButton)
	if err != nil {
		return errors.Wrap(err, "finding sign in button")
	}

	// click on sign in button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking sign in button")
	}

	time.Sleep(2 * time.Second)

	// check if we can find password field, if yes then password was wrong
	_, err = wd.FindElement(selenium.ByCSSSelector, passwordInput)
	if err == nil {
		return ErrInvalidPassword
	}

	return nil
}

// LoginByPhone performs a login to Acroplia using phone
func LoginByPhone(phone, password string, wd selenium.WebDriver) error {
	var err error

	// navigate to login page
	err = wd.Get(acropliaLoginURL)
	if err != nil {
		return errors.Wrapf(err, "navigating to %s", acropliaLoginURL)
	}

	// check if old version of browser is shown
	elem, err := wd.FindElement(selenium.ByClassName, continueButton)
	if err == nil {
		err = elem.Click() // click on continue button to proceed to login page
		if err != nil {
			return errors.Wrap(err, "clicking on continue anyway button")
		}
	}

	// wait until login page loads
	time.Sleep(1 * time.Second)

	// find by phone number link
	elem, err = wd.FindElement(selenium.ByCSSSelector, "a")
	if err != nil {
		return errors.Wrap(err, "finding by phone number link")
	}

	// click on link
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking on by phone number link")
	}

	// find phone input
	elem, err = wd.FindElement(selenium.ByCSSSelector, phoneInput)
	if err != nil {
		return errors.Wrap(err, "finding phone input")
	}

	// clear initial value of phone input
	err = elem.Clear()
	if err != nil {
		return errors.Wrap(err, "clearing phone input")
	}

	// remove + sign from phone number if exists
	if i := strings.Index(phone, "+"); i > -1 {
		phone = phone[i+1:]
	}

	// input new value to phone input
	err = elem.SendKeys(phone)
	if err != nil {
		return errors.Wrap(err, "filling phone input")
	}

	// find next button
	elem, err = wd.FindElement(selenium.ByCSSSelector, submitButton)
	if err != nil {
		return errors.Wrap(err, "finding next button")
	}

	// click on next button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking next button")
	}

	// wait until password page loads up
	time.Sleep(2 * time.Second)

	// find password input
	elem, err = wd.FindElement(selenium.ByCSSSelector, passwordInput)
	if err != nil {
		return errors.Wrap(err, "finding password input")
	}

	// check if password input is visible, if no then email was wrong
	attr, _ := elem.GetAttribute("aria-hidden")
	if attr == "true" {
		return ErrInvalidPhone
	}

	// fill password input
	err = elem.SendKeys(password)
	if err != nil {
		return errors.Wrap(err, "filling password input")
	}

	// find sign in button
	elem, err = wd.FindElement(selenium.ByCSSSelector, submitButton)
	if err != nil {
		return errors.Wrap(err, "finding sign in button")
	}

	// click on sign in button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking sign in button")
	}

	time.Sleep(2 * time.Second)

	// check if we can find password field, if yes then password was wrong
	_, err = wd.FindElement(selenium.ByCSSSelector, passwordInput)
	if err == nil {
		return ErrInvalidPassword
	}

	return nil
}
