package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/travelaudience/armador/internal/logger"

	"go.uber.org/zap"
)

type Command interface {
	Exec(string, ...string) ([]string, error)
	ExecInDir(string, string, ...string) ([]string, error)
	ExecUnparsed(string, ...string) (string, error)
	ExecInDirUnparsed(string, string, ...string) (string, error)
}

type Cmd struct {
	Ctx context.Context
}

// Exec generates a command from a string, runs it, then returns that
// commands stdout.
func (c Cmd) Exec(name string, commands ...string) ([]string, error) {
	return c.ExecInDir(name, ".", commands...)
}

func (c Cmd) ExecInDir(name string, dir string, commands ...string) ([]string, error) {
	outStr, err := c.ExecInDirUnparsed(name, dir, commands...)
	if err != nil {
		return nil, err
	}
	if len(outStr) > 0 {
		logger.GetLogger().Debugf("\n%s", outStr)
	}
	return parseStdout(outStr), nil
}

func (c Cmd) ExecUnparsed(name string, commands ...string) (string, error) {
	return c.ExecInDirUnparsed(name, ".", commands...)
}

// ExecInDirUnparsed generates a command from a string, runs it within the specified directory
func (c Cmd) ExecInDirUnparsed(name, dir string, commands ...string) (string, error) {
	// Set up a logger with the command and name fields. By doing this, any time we
	// log something, the command and name will be printed by default.
	logger := logger.GetLogger(zap.Fields(
		zap.String("name", name),
		zap.String("dir", dir),
		zap.String("command", strings.Join(commands, " "))))

	// Empty command is not allowed
	if len(commands) == 0 {
		return "", errors.New("Cmd: command is mandatory")
	}

	logger.Debug("running command")

	// Create the channel to signal when command finishes
	done := make(chan error)

	// Construct command
	cmd := exec.Command(commands[0], commands[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = dir

	err := cmd.Start()
	if err != nil {
		logger.Errorf("Command execution not possible %s", err)
		return "", err
	}

	go func() { done <- cmd.Wait() }() // When command finishes, it will signal done

	// Context is expected to have a timeout. Either the command finishes, or the ctx times out
	select {
	case err := <-done:
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if err != nil {
			logger.Errorf("An error occured in the command: %v with string %v", err, errStr)
			return "", fmt.Errorf("cmd: Run failed with %s", errStr)
		}
		return outStr, nil
	case <-c.Ctx.Done():
		logger.Error("Command timed out")
		return "", fmt.Errorf("command could not complete: %v", c.Ctx.Err())
	}
}

func parseStdout(cmdOut string) []string {
	output := []string{}
	lines := strings.Split(cmdOut, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		output = append(output, strings.Join(fields, " "))
	}
	return output
}
