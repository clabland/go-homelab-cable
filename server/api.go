package server

import (
	"net/http"

	"github.com/clabland/go-homelab-cable/domain"
	"github.com/labstack/echo/v4"
)

func (s *Server) getNetworks(e echo.Context) error {
	// There's only one network for now.
	return e.JSON(http.StatusOK, []domain.Network{
		{
			Name:     s.Network.Name,
			Owner:    s.Network.Owner,
			CallSign: "KHLC", // "KHomeLabCable"
		},
	})
}

func (s *Server) getChannels(e echo.Context) error {
	channels := make([]any, 0)
	for _, c := range s.Network.Channels() {
		channels = append(channels, domain.ToChannelModel(s.Network, c))
	}
	return e.JSON(http.StatusOK, channels)
}

func (s *Server) getChannel(e echo.Context) error {
	c, err := s.Network.Channel(e.Param("channel_id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, domain.ToChannelModel(s.Network, c))
}

func (s *Server) setChannelLive(e echo.Context) error {
	c, err := s.Network.Channel(e.Param("channel_id"))
	if err != nil {
		return err
	}
	err = s.Network.SetChannelLive(c.ID)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, domain.ToChannelModel(s.Network, c))
}

func (s *Server) playNext(e echo.Context) error {
	c, err := s.Network.Channel(e.Param("channel_id"))
	if err != nil {
		return err
	}
	_ = c.PlayNext()
	return e.JSON(http.StatusOK, domain.ToChannelModel(s.Network, c))
}

func (s *Server) playLiveNext(e echo.Context) error {
	c, err := s.Network.CurrentChannel()
	if err != nil {
		return err
	}
	_ = c.PlayNext()
	return e.JSON(http.StatusOK, domain.ToChannelModel(s.Network, c))
}

func (s *Server) liveChannel(e echo.Context) error {
	c, err := s.Network.CurrentChannel()
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, domain.ToChannelModel(s.Network, c))
}
