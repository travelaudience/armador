package cluster

import (
	"reflect"
	"testing"

	"github.com/travelaudience/armador/internal/commands"
)

func TestCreateNamespace(t *testing.T) {
	cases := []struct {
		name      string
		cmd       commands.Command
		namespace string
		want      []string
		wantErr   bool
	}{
		{
			name:    "any-name",
			cmd:     commands.CmdMock{CmdExecuted: "kubectl create ns any-name"},
			want:    []string{"Dir: ., Cmd: kubectl create ns any-name"},
			wantErr: false,
		},
		{
			name:    "catch-error",
			cmd:     commands.CmdMock{ReturnError: true},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateNamespace(tt.cmd, tt.name)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected test error for %s: %v", tt.name, err)
				return
			} else if err == nil && tt.wantErr {
				t.Errorf("expected error, no error received for %s", tt.name)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s test failed. \nGot: %v \nExpected: %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_getNamespaces(t *testing.T) {
	cases := []struct {
		name string
		cmd  commands.Command
		want []string
	}{
		{
			name: "check-cmd",
			cmd:  commands.CmdMock{CmdExecuted: "kubectl get ns"},
			want: []string{"Dir: ., Cmd: kubectl get ns"},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNamespaces(tt.cmd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s test failed. \nGot: %v \nExpected: %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestNamespaceExists(t *testing.T) {
	cases := []struct {
		name string
		cmd  commands.Command
		want bool
	}{
		{
			name: "any-name",
			cmd: commands.CmdMock{ReturnValueUnparsed: `
NAME              STATUS    AGE
default           Active    2d
any-name          Active    5d
kube-public       Active    5d
kube-system       Active    5d
`},
			want: true,
		},
		{
			name: "non-exist",
			cmd: commands.CmdMock{ReturnValueUnparsed: `
NAME              STATUS    AGE
default           Active    2d
kube-system       Active    5d
`},
			want: false,
		},
		{
			name: "similar-not",
			cmd: commands.CmdMock{ReturnValueUnparsed: `
NAME              STATUS    AGE
default           Active    2d
similar           Active    2d
similar-name      Active    2d
`},
			want: false,
		},
		{
			name: "catch-error",
			cmd:  commands.CmdMock{ReturnError: true},
			want: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := NamespaceExists(tt.cmd, tt.name); got != tt.want {
				t.Errorf("%s test failed. \nGot: %v \nExpected: %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestListEnvironments(t *testing.T) {
	cases := []struct {
		name string
		cmd  commands.Command
		want []string
	}{
		{
			name: "removed-k8s-defaults",
			cmd: commands.CmdMock{ReturnValueUnparsed: `
NAME              STATUS    AGE
any-name          Active    5d
default           Active    2d
kube-system       Active    5d
test              Active    10d
`},
			want: []string{"any-name", "test"},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListEnvironments(tt.cmd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s test failed. \nGot: %v \nExpected: %v", tt.name, got, tt.want)
			}
		})
	}
}
