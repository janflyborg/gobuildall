package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// Builds all packages in the current working directory and reports any errors to stdout
func buildPackages(done chan struct{}) {
	defer func() { done <- struct{}{}}()
	cmd := exec.Command("go", "build", "./...")
	buff := bytes.Buffer{}
	cmd.Stderr = &buff
	err := cmd.Run()
	if err == nil {
		return
	}
	if _, ok := err.(*exec.ExitError); !ok {
		// Errors that are not an ExitError is a problem
		log.Fatal(err)
	}
	fmt.Print(buff.String())
}

// Builds all tests in the current working directory and reports any errors to stdout
func buildTests(done chan struct{}) {
	defer func() { done <- struct{}{}}()
	cmd := exec.Command("go", "test", "-run", "unlikelypackagename", "./...")
	buff := bytes.Buffer{}
	cmd.Stderr = &buff
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
			fmt.Println(l)
		}
	}
}

func main() {
	// Build all packages and the tests in parallell
	doneChan := make(chan struct{})
	go buildPackages(doneChan)
	go buildTests(doneChan)
	<-doneChan
	<-doneChan
}
