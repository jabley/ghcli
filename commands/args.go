package commands

import (
	"fmt"
	"strings"
)

type Args struct {
	GlobalFlags []string
	Command     string
	ProgramPath string
	Params      []string
	Noop        bool
	Terminator  bool
}

func NewArgs(args []string) *Args {
	var (
		command     string
		params      []string
		noop        bool
		globalFlags []string
	)

	slurpGlobalFlags(&args, &globalFlags)
	noop = removeValue(&globalFlags, noopFlag)

	if len(args) == 0 {
		params = []string{}
	} else {
		command = args[0]
		params = args[1:]
	}

	return &Args{
		GlobalFlags: globalFlags,
		Command:     command,
		Params:      params,
		Noop:        noop,
	}
}

const (
	noopFlag    = "--noop"
	versionFlag = "--version"
	helpFlag    = "--help"
	flagPrefix  = "-"
)

func looksLikeFlag(value string) bool {
	return strings.HasPrefix(value, flagPrefix)
}

func slurpGlobalFlags(args *[]string, globalFlags *[]string) {
	slurpNextValue := false
	commandIndex := 0

	for i, arg := range *args {
		if slurpNextValue {
			commandIndex = i + 1
			slurpNextValue = false
		} else if arg == versionFlag || arg == helpFlag || !looksLikeFlag(arg) {
			break
		} else {
			commandIndex = i + 1
		}
	}

	if commandIndex > 0 {
		aa := *args
		*globalFlags = aa[0:commandIndex]
		*args = aa[commandIndex:]
	}
}

func (a *Args) FirstParam() string {
	if a.ParamsSize() == 0 {
		panic(fmt.Sprintf("Index 0 is out of bound"))
	}

	return a.Params[0]
}

func (a *Args) HasSubcommand() bool {
	return !a.IsParamsEmpty() && a.Params[0][0] != '-'
}

func (a *Args) IndexOfParam(param string) int {
	for i, p := range a.Params {
		if p == param {
			return i
		}
	}

	return -1
}

func (a *Args) ParamsSize() int {
	return len(a.Params)
}

func (a *Args) IsParamsEmpty() bool {
	return a.ParamsSize() == 0
}

func (a *Args) HasFlags(flags ...string) bool {
	for _, f := range flags {
		if i := a.IndexOfParam(f); i != -1 {
			return true
		}
	}

	return false
}

func removeItem(slice []string, index int) (newSlice []string, item string) {
	if index < 0 || index > len(slice)-1 {
		panic(fmt.Sprintf("Index %d is out of bound", index))
	}

	item = slice[index]
	newSlice = append(slice[:index], slice[index+1:]...)

	return newSlice, item
}

func removeValue(slice *[]string, value string) (found bool) {
	aa := *slice
	for i := len(aa) - 1; i >= 0; i-- {
		arg := aa[i]
		if arg == value {
			found = true
			*slice, _ = removeItem(*slice, i)
		}
	}
	return found
}
