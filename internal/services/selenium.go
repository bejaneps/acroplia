package services

import (
	"fmt"

	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/firefox"

	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

var (
	// SupportedBrowsers is a list of browsers, that are tested.
	SupportedBrowsers = []string{
		"firefox",
		"chrome",
	}
)

// NewSelenium creates new Selenium WebDriver using specified port and browser.
// Don't forget to quit the driver after the job is done.
func NewSelenium(port, browser string, args ...string) (selenium.WebDriver, error) {
	seleniumURL := fmt.Sprintf("http://localhost:%s/wd/hub", port)
	caps := selenium.Capabilities{"browserName": browser}
	if browser == "firefox" {
		firefoxCaps := firefox.Capabilities{}
		firefoxCaps.Args = append(firefoxCaps.Args, args...)

		caps.AddFirefox(firefoxCaps)
	} else if browser == "chrome" {
		chromeCaps := chrome.Capabilities{}
		chromeCaps.Args = append(chromeCaps.Args, args...)

		caps.AddChrome(chromeCaps)
	} else {
		return nil, errors.Errorf("%s browser is not supported", browser)
	}

	driver, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		return nil, errors.Wrap(err, "creating new Selenium Webdriver")
	}

	return driver, nil
}
