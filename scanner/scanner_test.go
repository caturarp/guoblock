package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
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

func TestScanGitDiff(t *testing.T) {
	// Set up temporary git repo
	tmpDir, err := os.MkdirTemp("", "guoblock-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Create and stage a file with secrets
	secretFile := filepath.Join(tmpDir, "secret.txt")
	content := `aws_secret_key='wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'`
	if err := os.WriteFile(secretFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "add", "secret.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	findings, err := ScanGitDiff()
	if err != nil {
		t.Fatal(err)
	}

	// Verify findings
	if len(findings) == 0 {
		t.Error("Expected to find secrets in git diff")
	}

	found := false
	for _, f := range findings {
		if f.Match == "AWS Secret Key" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find AWS Secret Key in diff")
	}
}

func TestScanGitDiffWithUnstagedSecrets(t *testing.T) {
	// Set up temporary git repo
	tmpDir, err := os.MkdirTemp("", "guoblock-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Create a file with secrets but don't stage it
	secretFile := filepath.Join(tmpDir, "secret.txt")
	content := `aws_secret_key='wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'`
	if err := os.WriteFile(secretFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	findings, err := ScanGitDiff()

	// Verify results
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}
	if len(findings) != 0 {
		t.Errorf("Expected no findings for unstaged secrets, got %d findings", len(findings))
	}
}
