# YARA Rule Cleaner
This CLI application written in Golang can **strip Metadata** and **Tags** from YARA rules. It **removes comments** and **formats** the rules as well.

## Example
```sh
yara_cleaner -output CLEANED_RULES_DIR -stripMeta -stripTags -recursive RAW_RULES_DIR
```
## Usage
```
yara_cleaner -h

Usage:  yara_cleaner [OPTION]... FILE | DIR
Provide program flags and at least one directory or file to scan.
  -debug
        Enable debug mode
  -output string
        Output directory
  -recursive
        Recursively scan yara rules (default true)
  -strict
        Strict mode: exit with error if a rule is invalid
  -stripMeta
        Strip metadata from yara rules
  -stripTags
        Strip Tags from yara rules
```
## Build
```sh
make build
```

## Things to note
* This is a work in progress