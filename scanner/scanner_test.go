package scanner

import (
	"os"
	"testing"
)

func TestScanFileWithSecrets(t *testing.T) {
	content := `AWS_SECRET_KEY = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'`
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
func TestScanFileWithLowercaseSecrets(t *testing.T) {
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

func TestScanFileWithInvalidSecrets(t *testing.T) {
	// Invalid AWS keys
	content := `
		aws_access_key = "AKIA123"
		aws_secret = "short"
		aws_key = "invalid"
		api_key = "too_short"
		token = "123"
	`
	_ = os.WriteFile("test_invalid.txt", []byte(content), 0644)
	defer os.Remove("test_invalid.txt")

	results, err := scanFile("test_invalid.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("Expected no findings for invalid secrets, got %d that is: \n %v", len(results), results)
	}
}
