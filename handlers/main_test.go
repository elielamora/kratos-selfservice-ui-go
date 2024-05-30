package handlers

import (
	"flag"
	"fmt"
	"io"
	"os"
	"testing"
)

var (
	update = flag.Bool("update", false, "update the golden files for tests in the package")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func getAsset(t *testing.T, filename string) []byte {
	t.Helper()
	return goldenValue(t, filename, nil, false)
}

// goldenValue gets the golden value writing the file if the update flag was set
func goldenValue(t *testing.T, filename string, actual []byte, update bool) []byte {
	t.Helper()
	goldenPath := fmt.Sprintf("testdata/%s", filename)

	flags := os.O_RDWR
	if update {
		flags |= os.O_CREATE
	}

	f, err := os.OpenFile(goldenPath, flags, 0644)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if update {
		_, err := f.Write(actual)
		if err != nil {
			t.Fatalf("Error writing to file %s: %s", goldenPath, err)
		}

		return actual
	}

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", goldenPath, err)
	}
	return content
}
