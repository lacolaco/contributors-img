package goutils

import (
	"fmt"
	"regexp"
)

func ValidateRepositoryName(s string) error {
	if s == "" {
		return fmt.Errorf("repository name cannot be empty")
	}
	if match, err := regexp.MatchString(`^[\w\-._]+\/[\w\-._]+$`, s); !match || err != nil {
		return fmt.Errorf("invalid repository name: %s", s)
	}
	return nil
}
