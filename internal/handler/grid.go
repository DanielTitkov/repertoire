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
	// events
	eventCreateUser = "createUser"
	// params
	paramEmail = "email"
	paramAge   = "age"
)

type (
	GridModel struct {
		Session string
	}
)

func NewGridModel(s *live.Socket) *GridModel {
	m, ok := s.Assigns().(*GridModel)
	if !ok {
		return &GridModel{
			Session: fmt.Sprint(s.Session),
		}
	}

	return m
}

func (h *Handler) Grid() *live.Handler {
	t, err := template.ParseFiles(h.t+"layout.html", h.t+"grid.html")
	if err != nil {
		log.Fatal(err)
	}

	lvh, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}

	// Set the mount function for this handler.
	lvh.Mount = func(ctx context.Context, r *http.Request, s *live.Socket) (interface{}, error) {
		return NewGridModel(s), nil
	}

	lvh.HandleEvent(eventCreateUser, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		email := p.String(paramEmail)
		age := p.Int(paramAge)
		fmt.Println("lvh grid params", email, age)

		return NewGridModel(s), nil
	})

	return lvh
}
