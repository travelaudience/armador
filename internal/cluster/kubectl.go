package cluster

import (
	"strings"

	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/logger"
)

func CreateNamespace(cmd commands.Command, namespace string) ([]string, error) {
	return cmd.Exec("create-ns", "kubectl", "create", "ns", namespace)
}

func getNamespaces(cmd commands.Command) []string {
	logger := logger.GetLogger()

	res, err := cmd.Exec("get-ns", "kubectl", "get", "ns")
	if err != nil {
		logger.Error("Could not get namespaces")
	}
	return res
}

func NamespaceExists(cmd commands.Command, search string) bool {
	namespaces := getNamespaces(cmd)
	for _, ns := range namespaces {
		if strings.HasPrefix(ns, search+" ") {
			return true
		}
	}
	return false
}

func DeleteNamespace(cmd commands.Command, namespace string) error {
	_, err := cmd.Exec("delete-ns", "kubectl", "delete", "ns", namespace)
	if err != nil {
		return err
	}
	return nil
}

func ListEnvironments(cmd commands.Command) []string {
	namespaces := []string{}
	allNamespaces := getNamespaces(cmd)
	for _, res := range allNamespaces {
		ns := strings.Fields(res)[0]
		// TODO: CONFIG: create global settings for default namespaces
		if ns != "kube-system" && ns != "kube-public" && ns != "default" && ns != "NAME" && ns != "core" {
			namespaces = append(namespaces, ns)
		}
	}
	return namespaces
}
