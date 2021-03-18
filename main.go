package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	fileArg := flag.String(`f`, ``, `nftables config file (for example: /etc/nftables.conf)`)

	flag.Parse()

	if (len(os.Args) - 1) == 0 {
		// No CLI arguments given
		_, _ = fmt.Fprintf(os.Stdin, `See --help for usage`)
		os.Exit(0)
	}

	f, err := os.Open(*fileArg)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error: %v`, err)
		os.Exit(1)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `error: %v`, err)
		os.Exit(1)
	}

	if !fi.Mode().IsRegular() {
		_, _ = fmt.Fprintf(os.Stderr, `error: not a regular file: %q`, *fileArg)
		os.Exit(1)
	}

	// Write output to memory buffer
	var w bytes.Buffer
	w.Reset()
	w.Grow(1024 * 1024)

	r := bufio.NewScanner(f)

	lvl := 0

	for r.Scan() {
		// Remove all spaces
		line := strings.TrimSpace(r.Text())

		// Skip indent level change
		skip := false

		if strings.ContainsRune(line, '{') && strings.ContainsRune(line, '}') {
			// '{' and '}' at the same line, skip
			skip = true
		}

		if !skip && strings.HasSuffix(line, `}`) {
			lvl--
		}

		newline := strings.Repeat("\t", lvl) + line + "\n"
		if strings.TrimSpace(newline) == `` {
			// Empty line
			newline = "\n"
		}

		// Write to memory
		w.WriteString(newline)

		if !skip && strings.HasSuffix(line, `{`) {
			lvl++
		}
	}

	// Output memory buffer
	fmt.Print(w.String())

}
