package network

import (
	"context"
	"errors"

	"github.com/clabland/go-homelab-cable/player"
)

var ErrNetworkNoChannelPlaying = errors.New("network no channel playing")
var ErrNetworkChannelNotFound = errors.New("network channel not found")

var DefaultChannelName = "default"

type Network struct {
	Ctx   context.Context
	Name  string
	Owner string

	channels     map[string]*Channel
	tunedChannel string // The channel which is currently displaying video on the host
}

func NewNetwork(name string, owner string) *Network {
	if name == "" {
		name = "Homelab Cable"
	}
	if owner == "" {
		owner = "clabretro"
	}

	return &Network{
		Name:     name,
		Owner:    owner,
		channels: make(map[string]*Channel, 0),
	}
}

func (n *Network) AddChannel(list *player.MediaList) *Channel {
	c := NewChannel(list)
	if len(n.channels) == 0 {
		c.ID = DefaultChannelName
	}
	n.channels[c.ID] = c
	return c
}

func (n *Network) Channel(ID string) (*Channel, error) {
	if c, ok := n.channels[ID]; ok {
		return c, nil
	} else {
		return nil, ErrNetworkChannelNotFound
	}
}

func (n *Network) Channels() []*Channel {
	channels := make([]*Channel, 0, len(n.channels))
	for _, c := range n.channels {
		channels = append(channels, c)
	}
	return channels
}

func (n *Network) CurrentChannel() (*Channel, error) {
	if len(n.channels) == 0 || n.tunedChannel == "" {
		return nil, ErrNetworkNoChannelPlaying
	}
	return n.Channel(n.tunedChannel)
}

func (n *Network) SetChannelLive(ID string) error {
	c, err := n.Channel(ID)
	if err != nil {
		return err
	}
	current, err := n.CurrentChannel()
	if err != nil && !errors.Is(err, ErrNetworkNoChannelPlaying) {
		return err
	}
	if current != nil {
		err := current.PlayWith(&player.NullPlayer{})
		if err != nil {
			return err
		}
	}
	err = c.PlayWith(&player.VLCPlayer{})
	if err != nil {
		return err
	}
	n.tunedChannel = c.ID
	return nil
}

func (n *Network) Live() string {
	return n.tunedChannel
}
