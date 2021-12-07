package handler

import (
	"html/template"
	"log"

	"github.com/jfyne/live"
)

func (h *Handler) Home() *live.Handler {
	t, err := template.ParseFiles(h.t+"layout.html", h.t+"home.html")
	if err != nil {
		log.Fatal(err)
	}

	liveHandler, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}

	return liveHandler

}
