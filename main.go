package main

import (
	"log"
	"os"

	"github.com/MidyMidy-MC/factorio-bot/tgbot"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "config yaml file",
			Aliases: []string{"c"},
			Value:   "config.yaml",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "token",
			Usage:    "telegram bot token",
			Aliases:  []string{"t"},
			EnvVars:  []string{"TGBOT_TOKEN", "TOKEN"},
			Required: true,
		}),
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:     "group-id",
			Usage:    "telegram group ID, must less than 0",
			Aliases:  []string{"g"},
			EnvVars:  []string{"TG_GID"},
			Required: true,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "rcon-address",
			Usage:    "address for factorio rcon",
			Aliases:  []string{"addr", "address", "r"},
			EnvVars:  []string{"RCON_ADDR"},
			Required: true,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "rcon-password",
			Usage:    "password for factorio rcon",
			Aliases:  []string{"password", "p"},
			EnvVars:  []string{"RCON_PASSWORD", "PASSWORD"},
			Required: true,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "tech-translation-file",
			Usage: "key-value file for tech translation",
			Value: "tech-translation.yaml",
		}),
	}

	app := &cli.App{
		Name:   "factorio-bot",
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Flags:  flags,
		Action: tgbot.Run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
