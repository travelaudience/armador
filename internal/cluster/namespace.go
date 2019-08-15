package cluster

import (
	"fmt"
	"regexp"
)

func CheckName(namespace string) error {
	if len(namespace) < 2 || len(namespace) > 25 {
		return fmt.Errorf("%s is invalid. It must be greater than 2 chars, and less than 25", namespace)
	}

	validNameRegex := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`).MatchString
	if !validNameRegex(namespace) {
		errMsg := "Name must consist of lower case alphanumeric characters or '-' (e.g. 'tom',  or 'tom-new-feature')"
		return fmt.Errorf("%s \n\t %s is invalid", errMsg, namespace)
	}
	return nil
}
