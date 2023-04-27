package network

import (
	"github.com/clabland/go-homelab-cable/player"
	"github.com/google/uuid"
)

type Channel struct {
	ID   string
	list *player.MediaList
	p    player.Player
}

func NewChannel(list *player.MediaList) *Channel {
	return &Channel{
		ID:   uuid.New().String(),
		list: list,
	}
}

func (c *Channel) PlayWith(p player.Player) error {
	if c.p != nil {
		if err := c.p.Shutdown(); err != nil {
			return err
		}
	}
	c.p = p

	err := p.Init()
	if err != nil {
		return err
	}
	return p.Play(c.list)
}

func (c *Channel) UpNext() string {
	return c.list.Next()
}

func (c *Channel) Current() string {
	return c.list.Current()
}

func (c *Channel) PlayNext() string {
	return c.list.Advance()
}
