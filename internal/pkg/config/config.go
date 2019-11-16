package config

import (
	"fmt"

	"github.com/GwangGwang/ganeungbot/pkg/util"
)

const secretDir string = "/secrets"

// Get returns the config with the given name
func Get(name string) (string, error) {

	configDir := fmt.Sprintf("%s/%s", secretDir, name)
	value, err := util.FileReadString(configDir)
	if err != nil {
		return "", fmt.Errorf("config '%s' not found at '%s'", name, configDir)
	}

	return value, nil
}
