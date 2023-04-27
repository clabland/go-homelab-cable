package player

import (
	vlc "github.com/adrg/libvlc-go/v3"
)

type VLCPlayer struct {
	list *MediaList

	player    *vlc.Player
	currMedia *vlc.Media

	next chan struct{}
}

func (p *VLCPlayer) Init() error {
	err := vlc.Init("--quiet")
	if err != nil {
		return err
	}

	p.next = make(chan struct{})

	p.player, err = vlc.NewPlayer()
	if err != nil {
		return err
	}

	p.player.SetFullScreen(true)

	manager, err := p.player.EventManager()
	if err != nil {
		return err
	}

	eventCallback := func(event vlc.Event, userData interface{}) {
		switch event {
		case vlc.MediaPlayerEndReached:
			p.next <- struct{}{}
		}
	}

	_, err = manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
	if err != nil {
		return err
	}

	go func(p *VLCPlayer) {
		for range p.next {
			err := p.PlayNext()
			if err != nil {
				panic(err)
			}
		}
	}(p)

	return nil
}

func (p *VLCPlayer) Shutdown() error {
	if p.player != nil {
		p.player.Stop()
		p.player.Release()
	}
	if p.currMedia != nil {
		p.currMedia.Release()
	}
	return vlc.Release()
}

func (p *VLCPlayer) Play(list *MediaList) error {
	if p.player == nil {
		return ErrPlayerNotInitialized
	}

	p.list = list

	var err error
	p.currMedia, err = p.player.LoadMediaFromPath(p.list.Current())
	if err != nil {
		return err
	}
	return p.player.Play()
}

func (p *VLCPlayer) PlayNext() error {
	if p.player == nil {
		return ErrPlayerNotInitialized
	}

	var err error
	err = p.player.Stop()
	if err != nil {
		return err
	}
	if p.currMedia != nil {
		p.currMedia.Release()
	}
	p.currMedia, err = p.player.LoadMediaFromPath(p.list.Advance())
	if err != nil {
		return err
	}
	return p.player.Play()
}

func (p *VLCPlayer) Next() string {
	return p.list.Next()
}

func (p *VLCPlayer) Current() string {
	return p.list.Current()
}
