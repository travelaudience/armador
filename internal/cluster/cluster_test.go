package cluster

import (
	"testing"

	"github.com/travelaudience/armador/internal/commands"
)

func TestClusterConnect(t *testing.T) {
	tests := []struct {
		name          string
		cmd           commands.Command
		clusterConfig ClusterConfig
		wantErr       bool
	}{
		{
			name:          "minikube-connect",
			cmd:           commands.CmdMock{},
			clusterConfig: ClusterConfig{Type: "minikube", Name: "test"},
			wantErr:       false,
		},
		{
			name:          "google-connect",
			cmd:           commands.CmdMock{},
			clusterConfig: ClusterConfig{Type: "google", Name: "test", Zone: "eu", Project: "test"},
			wantErr:       false,
		},
		{
			name:          "google-failed-connect",
			cmd:           commands.CmdMock{ReturnError: true},
			clusterConfig: ClusterConfig{Type: "google", Name: "test", Zone: "eu", Project: "test"},
			wantErr:       true,
		},
		{
			name:          "unsupported-type",
			cmd:           commands.CmdMock{ReturnError: false},
			clusterConfig: ClusterConfig{Type: "invalid"},
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ClusterConnect(tt.cmd, tt.clusterConfig)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected test error for %s: %v", tt.name, err)
				return
			} else if err == nil && tt.wantErr {
				t.Errorf("expected error, no error received for %s", tt.name)
				return
			}
		})
	}
}
