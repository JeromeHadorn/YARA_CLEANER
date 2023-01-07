package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Northern-Lights/yara-parser/data"
	"github.com/Northern-Lights/yara-parser/grammar"
)

// TODO: short and long flags
// TODO: Nice Logging
// TODO: Option for merging all rules into one YARA File
var (
	stripMeta bool
	stripTags bool

	recursive bool
	strict    bool
	debug     bool

	output string
)

func main() {
	flag.BoolVar(&stripMeta, "stripMeta", false, "Strip metadata from yara rules")
	flag.BoolVar(&stripTags, "stripTags", false, "Strip Tags from yara rules")
	flag.BoolVar(&recursive, "recursive", true, "Recursively scan yara rules")
	flag.BoolVar(&strict, "strict", false, "Strict mode: exit with error if a rule is invalid")
	flag.StringVar(&output, "output", "", "Output directory")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")

	flag.Usage = usage
	flag.Parse()
	checkUsage(flag.NArg())

	scanDirs := flag.Args()

	for dir := range scanDirs {
		if _, err := os.Stat(scanDirs[dir]); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR - File/Directory '%s' does not exist.\n", scanDirs[dir])
			os.Exit(1)
		}

		if recursive {
			ScanFiles(scanDirs[dir], output, strict, stripMeta, stripTags)
		} else {
			baseDir := filepath.Dir(scanDirs[dir])
			ScanFile(baseDir, scanDirs[dir], output, strict, stripMeta, stripTags)
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:  yara_cleaner [OPTION]... FILE | DIR`)
	fmt.Fprintln(os.Stderr, "Provide program flags and at least one directory or file to scan.")
	flag.PrintDefaults()
	os.Exit(1)
}

func checkUsage(nargs int) {
	if nargs < 1 {
		fmt.Fprintln(os.Stderr, `ERROR - Unexpected number of arguments. Please see below all accepted arguments and their default values.`)
		usage()
	}
}

func ScanFile(baseDir string, file string, outputDir string, strict bool, stripMeta bool, stripTags bool) error {
	relativePath := strings.Replace(file, baseDir, "", 1)
	newPath := filepath.Join(outputDir, relativePath)
	newDir := filepath.Dir(newPath)

	if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
		panic(err)
	}

	// Copy Non-YARA files
	// TODO: YARA rules saved in files without .yar | .yara extension
	if !strings.Contains(filepath.Ext(file), ".yar") {
		return CopyFile(file, newPath)
	}

	// Open YARA file
	yaraFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("can not open file: %s, err: %v", file, err)
	}
	defer yaraFile.Close()

	// Parse YARA file
	ruleset, err := grammar.Parse(yaraFile, os.Stdout)
	if err != nil {
		if strict {
			return fmt.Errorf("file couldn't be parsed, err: %s", err)
		}
		fmt.Printf("Skipped unparsable file: %s, err: %s\n", file, err)
		return nil
	}

	// Optional - strip metadata/tags
	if stripMeta || stripTags {
		Strip(&ruleset, stripMeta, stripTags)
	}

	// Serialize Yara Files
	serialized, err := ruleset.Serialize()
	if err != nil {
		return fmt.Errorf("YARA file '%s' could not seralized, err: %s", relativePath, err)
	}

	// Drop polished file
	if err := DropFile(newPath, serialized); err != nil {
		return fmt.Errorf("YARA file '%s' could not be created, err: %s", newPath, err)
	}

	return nil
}

func ScanFiles(baseDir string, outputDir string, strict bool, stripMeta bool, stripTags bool) {
	filepath.Walk(baseDir, func(pathEntry string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return ScanFile(baseDir, pathEntry, outputDir, strict, stripMeta, stripTags)
	})
}

// Strip data from a YARA Rule Set
func Strip(ruleset *data.RuleSet, stripMeta bool, stripTags bool) {
	stripped_rules := make([]data.Rule, 0, len(ruleset.Rules))
	for _, rule := range ruleset.Rules {
		if stripMeta {
			rule.Meta = nil
		}

		if stripTags {
			rule.Tags = []string{}
		}
		stripped_rules = append(stripped_rules, rule)
	}
	ruleset.Rules = stripped_rules
}
