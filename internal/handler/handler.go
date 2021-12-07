package handler

import (
	"github.com/DanielTitkov/repertoire/internal/app"
)

type (
	Handler struct {
		app *app.App
		t   string // template path
	}
)

func NewHandler(
	app *app.App,
	t string,
) *Handler {
	return &Handler{
		app: app,
		t:   t,
	}
}
