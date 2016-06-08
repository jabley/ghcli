package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type ExecError struct {
	Err      error
	ExitCode int
}

func (execError *ExecError) Error() string {
	return execError.Err.Error()
}

func newExecError(err error) ExecError {
	exitCode := 0
	if err != nil {
		exitCode = 1
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}

	return ExecError{Err: err, ExitCode: exitCode}
}

type Runner struct {
	commands map[string]*Command
}

func NewRunner() *Runner {
	return &Runner{
		commands: make(map[string]*Command),
	}
}

func (r *Runner) Use(command *Command, aliases ...string) {
	r.commands[command.Name()] = command
	if len(aliases) > 0 {
		r.commands[aliases[0]] = command
	}
}

func (r *Runner) Lookup(name string) *Command {
	return r.commands[name]
}

func (r *Runner) Execute() ExecError {
	args := NewArgs(os.Args[1:])
	args.ProgramPath = os.Args[0]

	if args.Command == "" {
		printUsage()
		return newExecError(fmt.Errorf(""))
	}

	cmd := r.Lookup(args.Command)

	if cmd != nil && cmd.Runnable() {
		return r.Call(cmd, args)
	}

	return newExecError(nil)
}

func (r *Runner) Call(cmd *Command, args *Args) ExecError {
	err := cmd.Call(args)
	if err != nil {
		if err == flag.ErrHelp {
			err = nil
		}
	}
	return newExecError(err)
}
