package tgbot

import (
	"fmt"
	"log"
	"time"

	rconcontroller "github.com/MidyMidy-MC/factorio-bot/rcon_controller"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	readTranslation(c)

	id := c.Int64("group-id")
	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint(
		c.String("token"),
		c.String("api-url"),
	)
	if err != nil {
		return err
	}

	u, err := bot.GetMe()
	if err != nil {
		return err
	}
	log.Printf("bot login succeed: %s(%d)", u.String(), u.ID)

	rcon := rconcontroller.New(
		c.String("rcon-address"),
		c.String("rcon-password"),
	)
	log.Print("rcon connected")

	cmd := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "online",
			Description: "players online",
		},
	)

	if _, err := bot.Request(cmd); err != nil {
		return err
	}
	log.Print("command set")

	update := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for {
		select {
		case u := <-update:
			onMessageReceived(bot, u, rcon, id)
		case <-time.After(500 * time.Millisecond):
			onUpdate(bot, rcon, id)
		}
	}
}

func onMessageReceived(bot *tgbotapi.BotAPI, u tgbotapi.Update, rcon *rconcontroller.Controller, id int64) {
	if u.Message == nil {
		return
	}

	if u.Message.Command() == "online" {
		players, err := rcon.GetPlayers(rconcontroller.FlagPlayerOnlne)
		if err != nil {
			log.Print(err)
			return
		}

		if len(players) == 0 {
			reply := tgbotapi.NewMessage(id, MessagePlayersNone)
			if _, err := bot.Send(reply); err != nil {
				log.Print(err)
			}
		} else {
			p := ""
			for _, player := range players {
				p = fmt.Sprintf("%s\n- %s", p, player.Name)
			}
			reply := tgbotapi.NewMessage(id, fmt.Sprintf("%s%s", MessagePlayers, p))
			if _, err := bot.Send(reply); err != nil {
				log.Print(err)
			}
		}
		return
	}

	if err := rcon.PostMessage(u.Message.From.FirstName, u.Message.Text); err != nil {
		log.Print(err)
	}
}

func onUpdate(bot *tgbotapi.BotAPI, rcon *rconcontroller.Controller, id int64) {
	ev, err := rcon.GetUpdate()
	if err != nil {
		log.Print(err)
		return
	}

	var msg tgbotapi.MessageConfig
	switch ev.EventType() {
	case rconcontroller.EventPlayerJoin:
		msg = tgbotapi.NewMessage(id,
			tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2,
				fmt.Sprintf(MessagePlayerJoin, ev.GetString("player_name")),
			),
		)
	case rconcontroller.EventPlayerLeft:
		msg = tgbotapi.NewMessage(id,
			tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2,
				fmt.Sprintf(MessagePlayerLeft, ev.GetString("player_name")),
			),
		)
	case rconcontroller.EventPlayerKilled:
		msg = tgbotapi.NewMessage(id,
			tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2,
				fmt.Sprintf(MessagePlayerKilled, ev.GetString("player_name"), ev.GetString("sarcasm")),
			),
		)
	case rconcontroller.EventResearchStarted:
		msg = tgbotapi.NewMessage(id,
			tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2,
				fmt.Sprintf(MessageResearchStart, getTranslation(
					ev.GetString("research_name"),
				)),
			),
		)
	case rconcontroller.EventResearchFinished:
		msg = tgbotapi.NewMessage(id,
			tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2,
				fmt.Sprintf(MessageResearchFinished, getTranslation(
					ev.GetString("research_name"),
				)),
			),
		)
	case rconcontroller.EventConsolePin:
		fallthrough
	case rconcontroller.EventConsoleChat:
		msg = tgbotapi.NewMessage(id,
			tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2,
				fmt.Sprintf("<%s> %s", ev.GetString("player_name"), ev.GetString("message")),
			),
		)
	case rconcontroller.EventConsoleMe:
		msg = tgbotapi.NewMessage(id,
			tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2,
				fmt.Sprintf("<%s> *%s", ev.GetString("player_name"), ev.GetString("message")),
			),
		)
	case rconcontroller.EventEmpty:
		return
	default:
		log.Print("unhandled event:", ev)
		return
	}

	msg.ParseMode = tgbotapi.ModeMarkdownV2
	reply, err := bot.Send(msg)
	if err != nil {
		log.Print(err)
		return
	}

	if ev.EventType() == rconcontroller.EventConsolePin {
		msg := tgbotapi.PinChatMessageConfig{
			ChatID:              id,
			MessageID:           reply.MessageID,
			DisableNotification: true,
		}
		if _, err := bot.Request(msg); err != nil {
			log.Print(err)
		}
	}
}
