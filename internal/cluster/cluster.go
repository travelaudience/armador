package cluster

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/travelaudience/armador/internal/commands"
)

type ClusterConfig struct {
	Type    string
	Name    string
	Zone    string
	Project string
}

func ClusterConnect(cmd commands.Command, clusterConfig ClusterConfig) error {
	switch clusterType := clusterConfig.Type; clusterType {
	case "google":
		return gcloudConnect(cmd, clusterConfig.Zone, clusterConfig.Project, clusterConfig.Name)
	case "minikube":
		return minikubeConnect(cmd, clusterConfig.Name)
	default:
		return fmt.Errorf("cluster type '%s' is not supported", clusterType)
	}
}

func gcloudConnect(cmd commands.Command, zone, project, cluster string) error {
	_, err := cmd.Exec("cluster-connect",
		"gcloud", "container", "clusters", "get-credentials",
		"--zone", zone,
		"--project", project, cluster)
	if err != nil {
		return fmt.Errorf("could not complete GcloudConnectCmd: %s", err)
	}
	return nil
}

func minikubeConnect(cmd commands.Command, contextName string) error {
	_, err := cmd.Exec("minikube-connect", "kubectl", "config", "use-context", contextName)
	if err != nil {
		return fmt.Errorf("could not not connect to %s: %s", contextName, err)
	}
	return nil
}

func GetClusterConfig() (retValues ClusterConfig, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to parse cluster info in config file")
			retValues = ClusterConfig{}
		}
	}()

	cluster := viper.Get("cluster").(map[string]interface{})

	if cluster["google"] != nil {
		retValues.Type = "google"

		vals := cluster["google"].(map[string]interface{})
		invalid := ""
		n, ok := vals["name"].(string)
		if !ok {
			invalid += "name "
		}
		z, ok := vals["zone"].(string)
		if !ok {
			invalid += "zone "
		}
		p, ok := vals["project"].(string)
		if !ok {
			invalid += "project "
		}
		if len(invalid) > 0 {
			return retValues, fmt.Errorf("google cluster configuration is missing: %s", invalid)
		}
		retValues.Name = n
		retValues.Zone = z
		retValues.Project = p
		return retValues, err
	}

	if cluster["minikube"] != nil {
		retValues.Type = "minikube"
		vals := cluster["minikube"].(map[string]interface{})
		invalid := ""
		// FIXME: contextname should be case sensitive: contextName
		n, ok := vals["contextname"].(string)
		if !ok {
			invalid += "contextName "
		}
		if len(invalid) > 0 {
			return retValues, fmt.Errorf("minikube configuration is missing: %s", invalid)
		}
		retValues.Name = n
		return retValues, err
	}

	// If it's here, then the retValues haven't been set and there is a configuration problem
	return retValues, fmt.Errorf("cluster configuration is not available: %v", cluster)
}
