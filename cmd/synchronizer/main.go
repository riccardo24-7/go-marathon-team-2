package main

import (
	"go-marathon-team-2/pkg/configuration"
	"go-marathon-team-2/pkg/logger"
	"go-marathon-team-2/pkg/syncer"
)

func main() {
	config, err := configuration.ConfigurationInit()
	if err != nil {
		logger.LogMessageFatal("Cannot load configuration", err)
	}
	syncer.StartApp(config)
}
