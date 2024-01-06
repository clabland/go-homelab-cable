package server

import (
	"io"
	"net/http"
	"text/template"

	"github.com/clabland/go-homelab-cable/domain"
	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Meta struct {
	Name  string
	Owner string
}

func (s *Server) getHtmxMeta(e echo.Context) error {
	return e.Render(http.StatusOK, "meta.html", Meta{Name: s.Network.Name, Owner: s.Network.Owner})
}

func (s *Server) getHtmxStatus(e echo.Context) error {
	c, err := s.Network.CurrentChannel()
	if err != nil {

		return err
	}
	return e.Render(http.StatusOK, "status.html", domain.ToChannelModel(s.Network, c))
}

func (s *Server) htmxPlayLiveNext(e echo.Context) error {
	return s.playLiveNext(e)
}
