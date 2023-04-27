package player

import (
	"time"
)

// NullPlayer advances the current item in the MediaList every 30 minutes. It (poorly) mimics the list of media being watched, as if it was on another channel.
type NullPlayer struct {
	list   *MediaList
	ticker *time.Ticker
	done   chan bool
}

func (n *NullPlayer) Init() error {
	n.ticker = time.NewTicker(time.Minute * 30)
	n.done = make(chan bool)
	return nil
}

func (n *NullPlayer) Play(list *MediaList) error {
	n.list = list
	go func() {
		for {
			select {
			case <-n.done:
				return
			case <-n.ticker.C:
				n.PlayNext()
			}
		}
	}()
	return nil
}

func (n *NullPlayer) PlayNext() error {
	n.list.Advance()
	return nil
}

func (n *NullPlayer) Next() string {
	return n.list.Next()
}

func (n *NullPlayer) Current() string {
	return n.list.Current()
}

func (n *NullPlayer) Shutdown() error {
	n.ticker.Stop()
	n.done <- true
	return nil
}
