package scanner

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Finding struct {
	File  string
	Line  int
	Match string
}

func ScanDirectory(root string) ([]Finding, error) {
	var results []Finding
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			// ignore dir
			return nil
		}
		if strings.Contains(path, ".git") || strings.Contains(path, "vendor") {
			return nil
		}

		lines, err := scanFile(path)
		if err == nil {
			results = append(results, lines...)
		}
		return nil
	})
	return results, err
}

// scanGitDiff scans the current git diff for secrets.
func ScanGitDiff() ([]Finding, error) {
	cmd := exec.Command("git", "diff", "--cached", "--unified=0")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var findings []Finding
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var currentFile string
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "+++ b/") {
			currentFile = strings.TrimPrefix(line, "+++ b/")
			continue
		}

		if strings.HasPrefix(line, "@@ ") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				lineInfo := strings.TrimPrefix(parts[2], "+")
				lineNum, _ = strconv.Atoi(strings.Split(lineInfo, ",")[0])
			}
			continue
		}

		// scan added lines
		if strings.HasPrefix(line, "+") {
			line = strings.TrimPrefix(line, "+")
			for _, rule := range Rules {
				if rule.Pattern.MatchString(line) {
					findings = append(findings, Finding{
						File:  currentFile,
						Line:  lineNum,
						Match: rule.Name,
					})
					break
				}
			}
			lineNum++
		}
	}

	return findings, nil
}

func scanFile(path string) ([]Finding, error) {
	var findings []Finding

	file, err := os.Open(path)
	if err != nil {
		return findings, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		for _, rule := range Rules {
			if rule.Pattern.MatchString(line) {
				findings = append(findings, Finding{
					File:  path,
					Line:  lineNum,
					Match: rule.Name,
				})
				break // Stop after first match
			}
		}
		lineNum++
	}
	return findings, nil
}
