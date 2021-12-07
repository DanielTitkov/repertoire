package main

import (
	"net/http"

	"github.com/DanielTitkov/repertoire/internal/app"

	"github.com/DanielTitkov/repertoire/internal/handler"

	"github.com/DanielTitkov/repertoire/internal/chat"
	"github.com/jfyne/live"
)

func main() {
	h := handler.NewHandler(
		&app.App{},
		"templates/",
	)

	// Run the server.
	http.Handle("/", h.Home())
	http.Handle("/chat", chat.NewHandler())
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
