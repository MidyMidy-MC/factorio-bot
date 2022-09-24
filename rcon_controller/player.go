package rconcontroller

import "fmt"

type Players struct {
	Players []*Player `json:"players"`
}

type Player struct {
	Name              string `json:"name"`
	Online            bool   `json:"connected"`
	AFKTime           int    `json:"afk_time"`
	LastOnline        int    `json:"last_online"`
	DisplayResolution struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"display_resolution"`
	Spectator bool `json:"spectator"`
}

func (p *Player) String() string {
	return fmt.Sprintf("player: %s", p.Name)
}
