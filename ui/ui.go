package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

// Interface defining how we want to print output to the world
type UI interface {
	Printf(format string, a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
	Errorf(format string, a ...interface{}) (n int, err error)
	Errorln(a ...interface{}) (n int, err error)
}

var (
	Stdout     = colorable.NewColorableStdout()
	Stderr     = colorable.NewColorableStderr()
	Default UI = Console{Stdout: Stdout, Stderr: Stderr}
)

// Printf prints the message to the UI
func Printf(format string, a ...interface{}) (n int, err error) {
	return Default.Printf(format, a...)
}

// Println prints the message to the UI
func Println(a ...interface{}) (n int, err error) {
	return Default.Println(a...)
}

// Errorf prints the error to the UI
func Errorf(format string, a ...interface{}) (n int, err error) {
	return Default.Errorf(format, a...)
}

// Errorln prints the message to the UI
func Errorln(a ...interface{}) (n int, err error) {
	return Default.Errorln(a...)
}

// IsTerminal returns true if the file is a tty
func IsTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd())
}

// Console is an abstraction for printing output
type Console struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (c Console) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stdout, format, a...)
}

func (c Console) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.Stdout, a...)
}

func (c Console) Errorf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stderr, format, a...)
}

func (c Console) Errorln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.Stderr, a...)
}
