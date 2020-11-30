package util

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type EnvVariables struct {
	DatabaseType string
	FileDBPath   string
}

func GetConfig() (EnvVariables, error) {
	var e EnvVariables
	err := envconfig.Process("crudstore", &e)

	if err != nil {
		// return log.F
		return EnvVariables{}, errors.Wrap(err, "Failed to read env file")
	}

	return e, nil
}
