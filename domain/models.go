package domain

import (
	"fmt"
	"path"

	"github.com/clabland/go-homelab-cable/network"
)

type Channel struct {
	ID      string `json:"id"`
	Playing string `json:"playing"`
	UpNext  string `json:"up_next"`
	Live    bool   `json:"live"`
}

type Network struct {
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	CallSign string `json:"call_sign"`
}

func ToChannelModel(n *network.Network, c *network.Channel) Channel {
	return Channel{
		ID:      c.ID,
		Playing: path.Base(c.Current()),
		UpNext:  path.Base(c.UpNext()),
		Live:    n.Live() == c.ID,
	}
}

func (c Channel) String() string {
	return fmt.Sprintf("Channel '%s'\nPlaying: %s\nUp Next: %s\nLive: %v\n", c.ID, c.Playing, c.UpNext, c.Live)
}
