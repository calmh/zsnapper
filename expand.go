package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var parensExp = regexp.MustCompile(`\$\(([^)]+)\)`)

func expandLines(lines []string) []string {
	var res []string
	for _, line := range lines {
		res = append(res, expand(line)...)
	}
	return res
}

func expand(line string) []string {
	m := parensExp.FindStringSubmatch(line)
	if len(m) < 2 {
		return []string{line}
	}

	bs, err := exec.Command("/bin/sh", "-c", m[1]).Output()
	if err != nil {
		if verbose {
			fmt.Printf("Expanding %q: %v\n", m[1], err)
		}
		return []string{line}
	}

	var res []string
	repl := "$(" + m[1] + ")"
	for _, s := range bytes.Split(bs, []byte("\n")) {
		if len(s) == 0 {
			continue
		}
		res = append(res, expand(strings.Replace(line, repl, string(s), 1))...)
	}

	return res
}

func selectLines(pats, lines []string) []string {
	var res []string
	for _, line := range lines {
		if lineMatches(pats, line) {
			res = append(res, line)
		}
	}
	return res
}

func lineMatches(pats []string, line string) bool {
	for _, pat := range pats {
		if matched, err := path.Match(pat, line); err == nil && matched {
			return true
		}
	}
	return false
}
