package server

import (
	"embed"
	"fmt"
	"math/rand"
	"text/template"
	"time"

	"github.com/clabland/go-homelab-cable/network"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// embedded files
var (
	//go:embed static/*
	staticFS embed.FS
	//go:embed templates/*.html
	templatesFS embed.FS
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
	// need to use MustSubFS since the embedded fs by default includes the
	// subfolder name (in this case "static")
	// if the subfolder name changes, both the //go:embed directive
	// and this will need to be updated
	e.StaticFS("/", echo.MustSubFS(staticFS, "static"))

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseFS(templatesFS, "templates/*.html")),
	}

	e.Renderer = renderer
	e.GET("/api/networks", s.getNetworks)

	e.GET("/api/networks/:callsign/channels", s.getChannels)
	e.GET("/api/networks/:callsign/channels/:channel_id", s.getChannel)
	e.PUT("/api/networks/:callsign/channels/:channel_id/set_live", s.setChannelLive)
	e.PUT("/api/networks/:callsign/channels/:channel_id/play_next", s.playNext)
	e.GET("/api/networks/:callsign/live", s.liveChannel)

	// Routes that always just act upon the current live channel
	e.PUT("/api/networks/:callsign/live/next", s.playLiveNext)

	e.GET("/htmx/meta", s.getHtmxMeta)
	e.GET("/htmx/status", s.getHtmxStatus)
	e.PUT("/htmx/live/next", s.htmxPlayLiveNext)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.port)))
}
