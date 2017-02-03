package genutil

import (
	"os"
	"os/exec"
	"strings"
)

// GoFmt formats go code file
func GoFmt(filename string) error {
	if info, err := os.Stat(filename); err != nil || info == nil {
		// skip wrong file
		return nil
	}
	if !strings.HasSuffix(filename, ".go") {
		// ignore non-golang file
		return nil
	}
	const gofmt = "gofmt"
	if _, err := exec.LookPath(gofmt); err != nil {
		// do nothing if failed to lookup `gofmt`
		return nil
	}
	cmd := exec.Command(gofmt, "-w", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CppFormat formats cpp code file
func CppFormat(filename string) error {
	//TODO
	return nil
}
