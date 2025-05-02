package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caturarp/guoblock/scanner"
)

func main() {
	diffMode := flag.Bool("diff", false, "Scan only Git diff")
	flag.Parse()

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	var findings []scanner.Finding
	var err error

	if *diffMode {
		findings, err = scanner.ScanGitDiff()
	} else {
		findings, err = scanner.ScanDirectory(path)
	}

	if err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(2)
	}

	if len(findings) == 0 {
		fmt.Println("✅ No secrets found")
		os.Exit(0)
	}

	fmt.Println("🚨 Possible secrets found:")
	for _, f := range findings {
		fmt.Printf(" - %s:%d => %s\n", f.File, f.Line, f.Match)
	}
	os.Exit(1)
}
