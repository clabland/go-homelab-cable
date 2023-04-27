package server

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/clabland/go-homelab-cable/network"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	port    string
	Network *network.Network
}

func NewServer(port string, n *network.Network) *Server {
	rand.Seed(time.Now().UnixNano())
	return &Server{
		port:    port,
		Network: n,
	}
}

func (s *Server) Serve() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	pwd, err := os.Getwd()
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Static("/", pwd+"/server/static/")

	e.GET("/api/networks", s.getNetworks)

	e.GET("/api/networks/:callsign/channels", s.getChannels)
	e.GET("/api/networks/:callsign/channels/:channel_id", s.getChannel)
	e.PUT("/api/networks/:callsign/channels/:channel_id/set_live", s.setChannelLive)
	e.PUT("/api/networks/:callsign/channels/:channel_id/play_next", s.playNext)
	e.GET("/api/networks/:callsign/live", s.liveChannel)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.port)))
}
