package handler

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jfyne/live"
)

const (
	// template vars
	varSession = "session"
)

type (
	HomeInstance struct {
		Session string
	}
)

func NewHomeInstance(s *live.Socket) *HomeInstance {
	m, ok := s.Assigns().(*HomeInstance)
	if !ok {
		return &HomeInstance{
			Session: fmt.Sprint(s.Session),
		}
	}

	return m
}

func (h *Handler) Home() *live.Handler {
	t, err := template.ParseFiles(h.t+"layout.html", h.t+"home.html")
	if err != nil {
		log.Fatal(err)
	}

	lvh, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}

	// Set the mount function for this handler.
	lvh.Mount = func(ctx context.Context, r *http.Request, s *live.Socket) (interface{}, error) {
		return NewHomeInstance(s), nil
	}

	return lvh
}
