package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"unicode"
)

// General usage flags
var (
	read  = flag.String("f", "", "Input file")
	write = flag.String("o", "", "Output file")
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
	mask = flag.String("m", "X", "What character to mask with")
)

var (
	exitCode = 0
	regex    = new(regexp.Regexp)
)

func setup() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}
	regex = regexp.MustCompile(*exp)
}

// xer takes in a file and projects an X in place of non-whitespace
// This is intended to be used to look at the foramtting of code inline with
// the idea that the logical structure of code should be suggested by
// its formatting.
func main() {
	// call xerMain in a separate function
	// so that it can use defer and have them
	// run before the exit.
	setup()
	xerMain()
	os.Exit(exitCode)
}

func xerMain() {
	in := openInput(*read)
	defer in.Close()

	out := openOutput(*write)
	defer out.Close()

	cont := readContent(in)

	xedcont := maskRunes(cont)

	rgis := findRegexIndexes(cont)

	xedcont = unmaskByRegex(xedcont, cont, rgis)

	out.Write(xedcont)
	return
}

func openInput(s string) *os.File {
	return openFileOrUse(s, os.Stdin)
}

func openOutput(s string) *os.File {
	return openFileOrUse(s, os.Stdout)
}

func openFileOrUse(s string, f *os.File) *os.File {
	if s != "" {
		var (
			out *os.File
			err error
		)
		out, err = os.Open(*read)
		if err != nil {
			log.Fatal(err)
		}
		return out
	}
	return f
}

func readContent(f *os.File) []byte {
	cont, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("error: could not read file %s", *read)
	}
	return cont
}

func findRegexIndexes(b []byte) [][]int {
	return regex.FindAllIndex(b, -1)
}

func unmaskByRegex(masked, orig []byte, is [][]int) []byte {
	for _, pair := range is {
		for i := pair[0]; i < pair[1]; i++ {
			masked[i] = orig[i]
		}
	}
	return masked
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
			mskd[i] = byte((*mask)[0])
		}
	}
	return mskd
}
