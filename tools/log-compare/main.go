package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var lineRE = regexp.MustCompile(`SP: [\dA-F]{4}`)
var whitespaceRE = regexp.MustCompile(`\s+`)
var continueRE = regexp.MustCompile(`^\([\dA-F]{2}$`)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("./log-compare [go-log] [c-log]")
	}

	gfh, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("failed to open go log")
	}

	cfh, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal("failed to open c log")
	}

	gs := bufio.NewScanner(gfh)
	cs := bufio.NewScanner(cfh)

	var prevLine string
	var i int
	for {
		i++
		gl := nextLine(gs)
		cl := nextLine(cs)

		if gl == cl {
			println("Files are identical")
			return
		}

		if gl == "" || cl == "" {
			fmt.Printf("Files are identical until end of one: line(%d)\ngo(%s)\nc(%s)\n", i, gl, cl)
			return
		}

		gl = normalizeGo(gl)
		cl = normalizeC(cl)

		if gl != cl {
			fmt.Printf("Files diverge:\nGO: %s\nC:  %s\n\nPREV: %s\n", gl, cl, prevLine)
			return
		}
		prevLine = gl
	}
}

func nextLine(s *bufio.Scanner) string {
	for s.Scan() {
		line := s.Bytes()
		if lineRE.Match(line) {
			return string(line)
		}
	}

	return ""
}

func normalizeGo(line string) string {
	parts := whitespaceRE.Split(line[20:], -1)
	return strings.Join(parts, " ")
}

func normalizeC(line string) string {
	parts := whitespaceRE.Split(line, -1)

	var n int
	for i, p := range parts[4:] {
		i++
		if continueRE.Match([]byte(p)) {
			n = i + 4
			break
		}
	}

	return strings.Join(append(parts[:4], parts[n-1:]...), " ")
}
