package scanner

import (
	"os"
	"testing"
)

func TestScanFileWithSecrets(t *testing.T) {
	content := `aws_secret_key='wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'`
	_ = os.WriteFile("test.txt", []byte(content), 0644)
	defer os.Remove("test.txt")

	results, err := scanFile("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Error("Expected at least one finding")
	}
}
