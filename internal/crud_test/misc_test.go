package crud_test

import (
	"runtime"
	"testing"

	"github.com/bejaneps/acroplia/internal/crud"
)

func TestCheckInternetConnection(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOARCH != "amd64" {
		t.Fatal("test available only on linux amd64 arch")
	}

	if err := crud.CheckInternetConnection(); err != nil {
		t.Fatal(err)
	}
}
