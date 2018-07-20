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

func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// like bufio.ScanLines but with non-mandatory "\n"
	// i.e., a line break can be \r, \r\n, or just \n
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\n\r"); i >= 0 {
		if data[i] == '\r' {
			// if the current buffer happens to end with \r, and there might be more data...
			if len(data)-1 == i {
				// for instance, when we're not at EOF, and that data could start with \n
				// (which should not be treated as a separate newline)
				if !atEOF {
					// insist on seeing more data first
					return 0, nil, nil
				}
			} else if data[i+1] == '\n' {
				// we're not at the end, and there is a newline following
				return i + 2, data[0:i], nil
			}
		}
		// unless we've bailed out by now, we've got the basic case:
		// either a \n, or a \r followed by something that is not a \n
		// where that something could be EOF
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

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

	lineEnding := []byte("\n")
	if *dos {
		lineEnding = []byte("\r\n")
	}

	scanner := bufio.NewScanner(input)
	scanner.Split(scanLines)
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
		_, err2 := os.Stdout.Write(lineEnding)
		if err2 != nil {
			panic(err2)
		}
	}
}
