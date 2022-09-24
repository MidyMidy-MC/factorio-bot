package rconcontroller

import (
	"encoding/json"
	"fmt"

	"github.com/james4k/rcon"
)

type Controller struct {
	addr, password string
}

func New(addr, password string) *Controller {
	return &Controller{addr: addr, password: password}
}

func (c *Controller) PostMessage(username, message string) error {
	conn, err := rcon.Dial(c.addr, c.password)
	if err != nil {
		return err
	}
	defer conn.Close()

	if message == "" {
		return nil
	}

	msg := &userChatMessage{
		Messages: []msg{
			{
				Name:    username,
				Message: message,
			},
		},
	}

	p, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("/_midymidyws post_messages %s", p)
	if _, err := conn.Write(cmd); err != nil {
		return err
	}

	if _, _, err := conn.Read(); err != nil {
		return err
	}

	return nil
}

const FlagPlayerOnlne = true

func (c *Controller) GetPlayers(isOnline bool) ([]*Player, error) {
	conn, err := rcon.Dial(c.addr, c.password)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if _, err := conn.Write("/_midymidyws get_players"); err != nil {
		return nil, err
	}

	resp, _, err := conn.Read()
	if err != nil {
		return nil, err
	}

	players := &Players{
		Players: make([]*Player, 0),
	}
	if err := json.Unmarshal([]byte(resp), &players); err != nil {
		return nil, err
	}

	if isOnline {
		playersOnline := make([]*Player, 0)
		for _, player := range players.Players {
			if player.Online {
				playersOnline = append(playersOnline, player)
			}
		}
		return playersOnline, nil
	}
	return players.Players, nil
}

func (c *Controller) GetUpdate() (*Event, error) {
	conn, err := rcon.Dial(c.addr, c.password)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if _, err := conn.Write("/_midymidyws get_update"); err != nil {
		return nil, err
	}

	resp, _, err := conn.Read()
	if err != nil {
		return nil, err
	}

	update := &Event{}
	if err := json.Unmarshal([]byte(resp), &update.Munch); err != nil {
		return nil, err
	}

	return update, nil
}
