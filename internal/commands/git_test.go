package commands

import (
	"testing"
)

func TestGitGet(t *testing.T) {
	tests := []struct {
		name        string
		cmd         Command
		repo        string
		destination string
		wantErr     bool
	}{
		{
			name:        "valid-clone",
			cmd:         CmdMock{ReturnError: false},
			repo:        "github.com/armador",
			destination: "new-destination",
			wantErr:     false,
		},
		{
			name:        "no-destination",
			cmd:         CmdMock{ReturnError: false},
			repo:        "github.com/armador",
			destination: "",
			wantErr:     true,
		},
		{
			name:        "valid-pull",
			cmd:         CmdMock{ReturnError: false},
			repo:        "github.com/armador",
			destination: "../",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GitGet(tt.cmd, tt.repo, tt.destination)
			if err != nil && !tt.wantErr {
				t.Errorf("%s: unexpected test error: %v", tt.name, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("%s: expected error, no error received", tt.name)
			}
		})
	}
}
