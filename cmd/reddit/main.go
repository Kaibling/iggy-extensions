package main

import (
	"fmt"

	"github.com/kaibling/iggy-extensions/config"
	"github.com/kaibling/iggy-extensions/pkg/log"
)

var (
	buildTime string //nolint: gochecknoglobals
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err.Error()) //nolint: forbidigo

		return
	}

	l := log.New(cfg.LogLevel, cfg.LogJSON)
	l.Info("build time: " + buildTime)
}
