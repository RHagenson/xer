package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"unicode"
)

// General usage flags
var (
	read  = flag.String("f", "", "Input file (default: stdin)")
	write = flag.String("o", "", "Output file (default: stdout)")
	help  = flag.Bool("h", false, "Print help and exit")
)

// Unmasking flags
var (
	exp         = flag.String("x", "", "Regex to unmask")
	symbols     = flag.Bool("s", false, "Unmask symbols")
	punctuation = flag.Bool("p", false, "Unmask punctuation")
	digit       = flag.Bool("d", false, "Unmask digits")
	number      = flag.Bool("n", false, "Unmask numbers")
)

var (
	exitCode = 0
	masker   = byte('X')
	regex    = new(regexp.Regexp)
)

func usage() {
	flag.PrintDefaults()
}

// xer takes in a file and projects an X in place of non-whitespace
// This is intended to be used to look at the foramtting of code inline with
// the idea that the logical structure of code should be suggested by
// its formatting.
func main() {
	// call xerMain in a separate function
	// so that it can use defer and have them
	// run before the exit.
	xerMain()
	os.Exit(exitCode)
}

func xerMain() {
	// Setup
	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	regex = regexp.MustCompile(*exp)

	// Open input
	var in *os.File
	var err error
	switch *read {
	case "":
		in = os.Stdin
	default:
		if in, err = os.Open(*read); err != nil {
			log.Fatal(err)
		}
	}
	defer in.Close()

	// Read in content
	cont, err := ioutil.ReadAll(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: could not read file %s", *read)
		exitCode = 2
		return
	}

	// Mask content by rune
	xedcont := maskRunes(cont)

	// Unmask content by regex
	nonemsked := regex.FindAllIndex(cont, -1)
	for _, pair := range nonemsked {
		for i := pair[0]; i < pair[1]; i++ {
			xedcont[i] = cont[i]
		}
	}

	// Opne output
	var out *os.File
	switch *write {
	case "":
		out = os.Stdout
	default:
		if out, err = os.Open(*write); err != nil {
			log.Fatal(err)
		}
	}
	defer out.Close()

	out.Write(xedcont)
}

func maskRunes(b []byte) []byte {
	mskd := make([]byte, len(b))

	for i, r := range b {
		switch {
		case *symbols && unicode.IsSymbol(rune(r)):
			mskd[i] = r
		case *punctuation && unicode.IsPunct(rune(r)):
			mskd[i] = r
		case *digit && unicode.IsDigit(rune(r)):
			mskd[i] = r
		case *number && unicode.IsNumber(rune(r)):
			mskd[i] = r
		case unicode.IsSpace(rune(r)):
			mskd[i] = r
		default:
			mskd[i] = masker
		}
	}
	return mskd
}
