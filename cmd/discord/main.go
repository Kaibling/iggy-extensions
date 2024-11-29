package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaibling/iggy-extensions/config"
	"github.com/kaibling/iggy-extensions/pkg/log"
	"github.com/kaibling/iggy-extensions/service"
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
	// ctx, cancel := context.WithCancel(context.Background())

	// start broker
	client, err := service.NewNATSClient(cfg, l)
	if err != nil {
		l.Error(err)

		return
	}

	dc, err := service.NewDiscordClient(cfg.ChannelPrefix, cfg.DiscordToken, l, client)
	if err != nil {
		l.Error(err)

		return
	}

	go func() {
		if err := client.Subscribe(cfg.ChannelPrefix+".write", dc); err != nil {
			l.Error(err)
			return
		}
	}()

	l.Info("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dc.Close()
}
