package main

import (
	"flag"
	"fmt"
	"bytes"
	"io/ioutil"
	"io"
	"os"
	"path"
	// "strings"
	"unicode"
)

func verboseRemove(name string) error {
	fmt.Printf("Removing file, %s\n", name)
	return os.Remove(name)
}

var CR = []byte("\r")
var CRLF = []byte("\r\n")
var LF = []byte("\n")

func main() {
	unix := flag.Bool("unix", false, "convert to unix line-endings (\\n)")
	dos := flag.Bool("dos", false, "convert to dos line-endings (\\r\\n)")
	trim := flag.Bool("trim", false, "trim trailing whitespace")

	clean := flag.Bool("clean", false, "remove temporary files")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file1 file2\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Without any flags, %s performs no actions\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *unix && *dos {
		fmt.Fprintln(os.Stderr, "you cannot use -unix and -dos at the same time")
		os.Exit(64)
	}

	names := flag.Args()
	fmt.Printf("Processing %d file(s)\n", len(names))

	tmp_dir := os.TempDir()

	var open_perms os.FileMode // = 0

	for _, name := range names {
		file, err := os.OpenFile(name, os.O_RDWR, open_perms)
		if err != nil {
			panic(err)
		}
		// first, copy the file contents to a temporary file
		// tmp_name := path.Join(tmp_dir, path.Base(name))
		tmp_file, err := ioutil.TempFile(tmp_dir, path.Base(name))
		if err != nil {
			panic(err)
		}
		if *clean {
			defer verboseRemove(tmp_file.Name())
		}
		// we can't use os.Link since we need to keep the shell of the old file for its metadata (tags, etc.)
		_, err = io.Copy(tmp_file, file)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Backed up %s to %s\n", name, tmp_file.Name())
		// rewind the temporary file we just wrote to the beginning
		_, err = tmp_file.Seek(0, io.SeekStart)
		if err != nil {
			panic(err)
		}
		// TODO: implement this streamingly
		// contents is a []byte
		contents, err := ioutil.ReadAll(tmp_file)
		if err != nil {
			panic(err)
		}
		// fix line-endings
		if *unix {
			contents = bytes.Replace(contents, CRLF, LF, -1) // \r\n -> \n
			contents = bytes.Replace(contents, CR, LF, -1) // \r -> \n
		}
		if *dos {
			// so perverse :(
			contents = bytes.Replace(contents, LF, CRLF, -1) // \n -> \r\n
			// TODO: do not modify pre-existing CRLF line-endings
		}
		// trim whitespace
		if *trim {
			// TODO: handle CRLF somehow?
			lines := bytes.Split(contents, LF)
			for i := range lines {
				lines[i] = bytes.TrimRightFunc(lines[i], unicode.IsSpace)
			}
			contents = bytes.Join(lines, LF)
		}
		// write (potentially) changed contents
		// TODO: skip writing if contents has not changed
		// _, err := file.Seek(0, io.SeekStart)
		// original_file, err := os.Create(name)
		// original_file.Write(contents)
		err = ioutil.WriteFile(name, contents, open_perms)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Wrote changes (?) to %s\n", name)
	}
}
