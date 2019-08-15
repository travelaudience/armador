package commands

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

var cmdCases = []struct {
	name        string
	cmds        []string
	isErr       bool
	strExpect   string
	sliceExpect []string
}{
	{
		name:        "output-to-stdout",
		cmds:        []string{"echo", "hello-world"},
		isErr:       false,
		strExpect:   "hello-world",
		sliceExpect: []string{"hello-world"},
	},
	{
		name:        "output-to-stderr",
		cmds:        []string{"awk", `BEGIN {print "hello-world" > "/dev/stderr";}`},
		isErr:       false,
		strExpect:   "",
		sliceExpect: []string{},
	},
	{
		name:        "output-to-stderr-and-stdout",
		cmds:        []string{"awk", `BEGIN {print "hello-world" > "/dev/stderr"; print ":)"}`},
		isErr:       false,
		strExpect:   ":)",
		sliceExpect: []string{":)"},
	},
	{
		name:        "program-error",
		cmds:        []string{"awk", `BEGIN {syntax error + -}`},
		isErr:       true,
		strExpect:   "",
		sliceExpect: nil,
	},
	{
		name:        "program-times-out",
		cmds:        []string{"sleep", "10"}, // test's are configured with 3 second timer
		isErr:       true,
		strExpect:   "",
		sliceExpect: nil,
	},
	{
		name:        "no-commands-supplied",
		cmds:        []string{}, // test's are configured with 3 second timer
		isErr:       true,
		strExpect:   "",
		sliceExpect: nil,
	},
	{
		name:  "multi-line-output",
		cmds:  []string{"awk", `BEGIN {for (x=0;x<5;x++){if (x % 2 == 0) {print x;} else {print ""}}}`},
		isErr: false,
		strExpect: `0

2

4`,
		sliceExpect: []string{"0", "2", "4"},
	},
}

// Parsed output in slice format
func TestCmd_Exec(t *testing.T) {
	for _, test := range cmdCases {
		// All tests will be run with a 3 second timer.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		c := Cmd{
			Ctx: ctx,
		}
		got, err := c.Exec(test.name, test.cmds...)
		if err != nil && !test.isErr {
			t.Errorf("%s: unexpected parsed output test error: %v", test.name, err)
			continue
		} else if err == nil && test.isErr {
			t.Errorf("%s: expected error, no error received", test.name)
			continue
		}

		if !reflect.DeepEqual(got, test.sliceExpect) {
			t.Errorf("parsed output test failure:\n  TEST: %s\n  EXPECTED: %s (len: %d)\n  GOT: %s (len: %d)", test.name, test.sliceExpect, len(test.sliceExpect), got, len(got))
		}
	}
}

// Unparsed output in string format
func TestCmd_Unparsed(t *testing.T) {
	for i, test := range cmdCases {
		// All tests will be run with a 3 second timer.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		c := Cmd{
			Ctx: ctx,
		}

		got, err := c.ExecUnparsed(test.name, test.cmds...)
		if err != nil && !test.isErr {
			t.Errorf("%s: unexpected test error: %v", test.name, err)
			continue
		} else if err == nil && test.isErr {
			t.Errorf("%s: expected error, no error received", test.name)
			continue
		}

		// trim the leading/trailing whitespaces
		got = strings.TrimSpace(got)
		if got != test.strExpect {
			t.Errorf("unexpected test failure:\n  TEST: %s (%d)\n  EXPECTED: %s\n  GOT: %s", test.name, i, test.strExpect, got)
		}
	}
}
