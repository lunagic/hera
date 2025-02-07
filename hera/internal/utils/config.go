package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func MustReadYamlConfig(
	paths []string,
	target any,
) error {
	for _, path := range paths {
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}

		if err := yaml.Unmarshal(fileBytes, target); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("no config file found at any of the following locations: %s", strings.Join(paths, ", "))
}
