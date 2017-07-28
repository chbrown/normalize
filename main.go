package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

func main() {
	dos := flag.Bool("dos", false, "convert to dos line-endings (\\r\\n)")
	trim := flag.Bool("trim", false, "trim trailing whitespace")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [file]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Without any flags, convert to unix line-endings (\\n)\n")
		fmt.Fprintf(os.Stderr, "Without specifying file, reads from /dev/stdin\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	names := flag.Args()
	if len(names) > 1 {
		fmt.Fprintln(os.Stderr, "You must not supply more than one positional argument")
		os.Exit(64)
	}

	var input io.Reader
	if len(names) == 1 {
		// read from file
		var err error
		input, err = os.Open(names[0])
		// this input is an os.File, which implements io.Reader
		if err != nil {
			panic(err)
		}
	} else {
		// read from /dev/stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintln(os.Stderr, "/dev/stdin must be piped in")
			os.Exit(66)
		}
		input = os.Stdin
	}

	line_ending := []byte("\n")
	if *dos {
		line_ending = []byte("\r\n")
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Bytes()
		if *trim {
			// trim line in-place (this is effectively half of bytes.TrimSpace)
			line = bytes.TrimRightFunc(line, unicode.IsSpace)
		}
		// output line content
		_, err := os.Stdout.Write(line)
		if err != nil {
			panic(err)
		}
		// output line ending
		_, err2 := os.Stdout.Write(line_ending)
		if err2 != nil {
			panic(err2)
		}
	}
}
