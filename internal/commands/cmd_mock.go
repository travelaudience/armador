package commands

import (
	"errors"
	"fmt"
	"strings"
)

type CmdMock struct {
	ReturnValue         []string
	ReturnValueUnparsed string
	ReturnError         bool
	ReturnErrorMsg      string
	CmdExecuted         string
}

func (c CmdMock) Exec(n string, cmds ...string) ([]string, error) {
	return c.ExecInDir(n, ".", cmds...)
}
func (c CmdMock) ExecInDir(n string, d string, cmds ...string) ([]string, error) {
	outStr, err := c.ExecInDirUnparsed(n, d, cmds...)
	return parseStdout(outStr), err
}
func (c CmdMock) ExecUnparsed(n string, cmds ...string) (string, error) {
	return c.ExecInDirUnparsed(n, ".", cmds...)
}
func (c CmdMock) ExecInDirUnparsed(n string, d string, cmds ...string) (string, error) {
	if c.ReturnError {
		return "", errors.New(c.ReturnErrorMsg)
	}
	if c.ReturnValueUnparsed == "" {
		return fmt.Sprintf("Dir: %s, Cmd: %s", d, strings.Join(cmds, " ")), nil
	}
	return c.ReturnValueUnparsed, nil
}
