package util

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type EnvVariables struct {
	DatabaseType string `default:"filestore"`
	FileDBPath   string `default:"books.json"`
	BoltDBName   string `default:"main.db"`
	PGUsername   string
	PGPassword   string
	PGHost       string
	PGName       string
}

type postgresConf struct {
	PGUsername string
	PGPassword string
	PGHost     string
	PGName     string
}

func GetConfig() (EnvVariables, error) {
	var e EnvVariables
	err := envconfig.Process("crudstore", &e)

	if err != nil {
		// return log.F
		return EnvVariables{}, errors.Wrap(err, "Failed to read env variables")
	}

	return e, nil
}
