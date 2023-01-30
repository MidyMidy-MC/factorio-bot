package tgbot

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var translation map[string]string

func init() {
	translation = make(map[string]string)
}

func readTranslation(c *cli.Context) {
	f, err := os.Open(c.String("tech-translation-file"))
	if err != nil {
		log.Print("skip tech translation file reading")
		return
	}
	defer f.Close()

	yaml.NewDecoder(f).Decode(translation)
}

func getTranslation(k string) string {
	if v, ok := translation[k]; ok {
		return v
	}
	return k
}
