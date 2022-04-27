package configuration

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

type Configuration struct {
	Grafana struct {
		Address string `yaml:"address"`
		Key     string `yaml:"key"`
	} `yaml:"grafana"`
	MetricsAPI struct {
		Address string `yaml:"address"`
	} `yaml:"metricsAPI"`
	Schedule struct {
		Repeat  bool `yaml:"repeat"`
		Days    int  `yaml:"days"`
		Hours   int  `yaml:"hours"`
		Minutes int  `yaml:"minutes"`
	} `yaml:"schedule"`
}

func ConfigurationInit() (*Configuration, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("couldn't load .env file")
	}
	var config Configuration
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "grafana",
				Aliases:     []string{"gr"},
				Usage:       "URL for Grafana connection",
				EnvVars:     []string{"GRAFANA_URL"},
				Destination: &config.Grafana.Address,
			},
			&cli.StringFlag{
				Name:        "key",
				Aliases:     []string{"k"},
				Usage:       "Key for Grafana connection",
				EnvVars:     []string{"GRAFANA_KEY"},
				Destination: &config.Grafana.Key,
			},
			&cli.StringFlag{
				Name:        "metrics",
				Aliases:     []string{"m"},
				Usage:       "URL for MetricsAPI connection",
				EnvVars:     []string{"METRICSAPI_URL"},
				Destination: &config.MetricsAPI.Address,
			},
			&cli.BoolFlag{
				Name:        "repeat",
				Aliases:     []string{"r"},
				Usage:       "Repeat star by schedule",
				EnvVars:     []string{"SCHEDULE"},
				Destination: &config.Schedule.Repeat,
			},
			&cli.IntFlag{
				Name:        "days",
				Aliases:     []string{"d"},
				Usage:       "Set days for schedule",
				EnvVars:     []string{"DAYS"},
				Destination: &config.Schedule.Days,
			},
			&cli.IntFlag{
				Name:        "minutes",
				Aliases:     []string{"min"},
				Usage:       "Set minutes for schedule",
				EnvVars:     []string{"MINUTES"},
				Destination: &config.Schedule.Minutes,
			},
			&cli.IntFlag{
				Name:        "hours",
				Aliases:     []string{"hr"},
				Usage:       "Set hours for schedule",
				EnvVars:     []string{"HOURS"},
				Destination: &config.Schedule.Hours,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return nil, errors.New("problem with app.Run")
	}

	return &config, nil
}
