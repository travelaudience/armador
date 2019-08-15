package commands

import (
	"fmt"
	"os"

	"github.com/travelaudience/armador/internal/logger"
)

func GitGet(cmd Command, repo, destination string) error {
	// no desitinataion, so just try to clone and hope for the best
	if destination == "" {
		return fmt.Errorf("The path to clone to is not clearly defined, and may already exist")
	}

	// check if dir exists, is so just git pull, else git clone
	_, err := os.Stat(destination)
	if err != nil && !os.IsNotExist(err) {
		logger.GetLogger().Warnf("Problem checking if the repo already exists in %s\n %s", destination, err)
	}
	if err == nil {
		cmds := []string{"git", "pull", "--ff-only"}
		_, err = cmd.ExecInDir("git-update", destination, cmds...)
	} else {
		cmds := []string{"git", "clone", repo, destination}
		_, err = cmd.Exec("git-clone", cmds...)
	}
	if err != nil {
		return fmt.Errorf("Could not get %s\n %s", repo, err)
	}
	return nil
}
