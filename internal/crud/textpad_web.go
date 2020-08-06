package crud

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

const acropliaHomeURL = "https://demo.acroplia.com/connect/"

const (
	plusButton            = `div[data-onboarding="mainArea"] > div:nth-child(1) > div:nth-child(3)`
	createInLibraryButton = `div[data-onboarding="mainArea"] > div:nth-child(2) > div:nth-child(2) > div:nth-child(1) > div:nth-child(3) > div:nth-child(1)`
	noteButton            = `div.context_item:nth-child(1)`
	noteTitle             = `.clear-input`
	welcomeContinueButton = `.btn-main`
	communityButton       = `div > a:nth-child(1)`
	tipsSkipButton        = `div[data-xpath="portal"] > div:nth-child(1) > div:nth-child(1) > div:nth-child(4) > div:nth-child(1)`
)

// acropliaPre performs pre checkups before doing some action in Acroplia
func acropliaPre(wd selenium.WebDriver) error {
	// navigate to home page
	err := wd.Get(acropliaHomeURL)
	if err != nil {
		return errors.Wrapf(err, "navigating to %s", acropliaHomeURL)
	}

	// check if welcome message is present
	elem, err := wd.FindElement(selenium.ByCSSSelector, welcomeContinueButton)
	if err == nil {
		// click on it to proceed
		err = elem.Click()
		if err != nil {
			return errors.Wrap(err, "clicking on welcome continue button")
		}
	}

	// wait till page loads
	time.Sleep(1 * time.Second)

	return nil
}

// TextpadCreate creates a textpad in Acroplia using web
func TextpadCreate(title string, subtitle string, wd selenium.WebDriver) error {
	if err := acropliaPre(wd); err != nil {
		return err
	}

	// find community button
	elem, err := wd.FindElement(selenium.ByCSSSelector, communityButton)
	if err != nil {
		return errors.Wrap(err, "finding community button")
	}

	// click on community button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking community button")
	}

	// wait till community loads
	time.Sleep(4 * time.Second)

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

	// find plus button
	elem, err = wd.FindElement(selenium.ByCSSSelector, plusButton)
	if err != nil {
		return errors.Wrap(err, "finding plus button")
	}

	// click on plus button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking plus button")
	}

	// find create in library button
	elem, err = wd.FindElement(selenium.ByCSSSelector, createInLibraryButton)
	if err != nil {
		return errors.Wrap(err, "finding create in library button")
	}

	// click on create in library button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking create in library button")
	}

	// find note button
	elem, err = wd.FindElement(selenium.ByCSSSelector, noteButton)
	if err != nil {
		return errors.Wrap(err, "finding note button")
	}

	// click on note button
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking note button")
	}

	// wait until note loads
	time.Sleep(3 * time.Second)

	// find note title
	elem, err = wd.FindElement(selenium.ByCSSSelector, noteTitle)
	if err != nil {
		return errors.Wrap(err, "finding note title")
	}

	// clear note title
	err = elem.Clear()
	if err != nil {
		return errors.Wrap(err, "clearing note title")
	}

	// fill note title
	err = elem.SendKeys(title)
	if err != nil {
		return errors.Wrap(err, "filling note title")
	}

	time.Sleep(1 * time.Second)

	// find plus button
	elem, err = wd.FindElement(selenium.ByCSSSelector, plusButton)
	if err != nil {
		return errors.Wrap(err, "finding plus button")
	}

	// click on plus button so title saves
	err = elem.Click()
	if err != nil {
		return errors.Wrap(err, "clicking plus button")
	}

	time.Sleep(1 * time.Second)

	return nil
}
