# Jesting Jaguar - Golang Template Escaper

## Overview
A simple tool to escape Golang template brackets to prevent interpolation.

## Usage Cheatsheet
```sh
# Escape template brackets in a single file
jestingjaguar escape file.txt

# Recursively escape template brackets in a directory
jestingjaguar escape ./templates/

# Increase verbosity level
jestingjaguar escape file.txt -v
jestingjaguar escape ./templates/ -vv
```

## Installation
```sh
go install github.com/gkwa/jestingjaguar@latest
```

## Example
Input:
```
artifacts/{{ workflow.name }}
```

Output:
```
artifacts/{{"{{"}} workflow.name {{"}}"}}
```

## Claude.ai generated this app with the prompt output from this

```
rm -rf /tmp/t && boilerplate --missing-key-action zero --non-interactive --output-folder=/tmp/t --template-url git::https://github.com/gkwa/manyeast.git//use-golang-to-escape-golang-template --var GoModuleName=jestingjaguar && find /tmp/t && cat /tmp/t/use-golang-to-escape-golang-template.tmpl | less
```