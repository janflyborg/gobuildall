package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Builds all packages in the current working directory and reports any errors to stdout
func buildPackages(directory string, done chan struct{}) {
	defer func() { done <- struct{}{} }()
	cmd := exec.Command("go", "build", "./...")
	buff := bytes.Buffer{}
	cmd.Stderr = &buff
	cmd.Dir = directory
	err := cmd.Run()
	if err == nil {
		return
	}
	if _, ok := err.(*exec.ExitError); !ok {
		// Errors that are not an ExitError is a problem
		log.Fatal(err)
	}
	fmt.Fprint(os.Stderr, buff.String())
}

// Builds all tests in the current working directory and reports any errors to stdout
func buildTests(directory string, done chan struct{}) {
	defer func() { done <- struct{}{} }()
	cmd := exec.Command("go", "test", "-run", "unlikelypackagename", "./...")
	buff := bytes.Buffer{}
	cmd.Stderr = &buff
	cmd.Dir = directory
	err := cmd.Run()
	if err == nil {
		return
	}
	if _, ok := err.(*exec.ExitError); !ok {
		// Errors that are not an ExitError is a problem
		log.Fatal(err)
	}
	// Fix old Mac and Windows line endings
	out := buff.Bytes()
	out = bytes.Replace(out, []byte{13, 10}, []byte{10}, -1)
	out = bytes.Replace(out, []byte{13}, []byte{10}, -1)

	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		// Check if this is an compiler error line
		matched, err := regexp.Match(`(.*):[\d]+:[\d]+: .*`, []byte(l))
		if err != nil {
			panic(err)
		}
		if matched {
			fmt.Fprintln(os.Stderr, l)
		}
	}
}

// Stop with an error message
func fatal(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+s+"\n", args...)
	os.Exit(-1)
}

func main() {
	// Check if the user specified a set of directories to build in - Otherwise we use current working directory
	buildDirs := make([]string, 0)
	if len(os.Args) > 1 {
		buildDirs = os.Args[1:]
	} else {
		buildDirs = []string{"."}
	}

	// Build all packages and the tests in parallel for each directory given
	doneChan := make(chan struct{})
	for _, d := range buildDirs {
		if s, err := os.Stat(d); err != nil {
			if os.IsNotExist(err) {
				fatal("Directory %s does not exist.", d)
			}
			if !s.IsDir() {
				fatal("Path %s does not refer to a directory.", d)
			}
		}
		go buildPackages(d, doneChan)
		go buildTests(d, doneChan)
	}
	// Wait for all builds to finish
	for i := 0; i < 2*len(buildDirs); i++ {
		<-doneChan
	}
}
