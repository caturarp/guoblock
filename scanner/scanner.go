package scanner

import (
	"bufio"
	"os"
	"path/filepath"
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
