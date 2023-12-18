package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FromFile unmarshals a JSON or YAML configuration file from the given path,
// optionally performing validation if the Validator interface is implemented.
func FromFile(path string, dst any) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(dst)
	if err != nil {
		return fmt.Errorf("yaml: %w", err)
	}

	if v, ok := dst.(Validator); ok {
		err = v.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
