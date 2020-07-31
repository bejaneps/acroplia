package services

import (
	"runtime"
	"testing"

	"github.com/bejaneps/acroplia/internal/services"
)

func TestNewSelenium(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOARCH != "amd64" {
		t.Fatal("test available only on linux amd64 arch")
	}

	// test supported browsers,
	// if browser opens and it can navigate to google.com
	// then it passes test
	for _, b := range services.SupportedBrowsers {
		t.Logf("testing %s browser...", b)
		driver, err := services.NewSelenium("4444", b, "--headless") // run browser in headless mode without ui
		if err != nil {                                              // if it fails here, continue to next browser
			t.Errorf("couldn't start %s browser: %v", b, err)
			continue
		}

		err = driver.Get("https://www.google.com")
		if err != nil {
			t.Errorf("couldn't navigate to google.com using %s browser: %v", b, err)
		}

		driver.Close()

		t.Logf("test for %s browser passed", b)
	}
}
