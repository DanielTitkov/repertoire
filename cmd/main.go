package main

import (
	"log"
	"net/http"

	"github.com/DanielTitkov/repertoire/internal/app"

	"github.com/DanielTitkov/repertoire/internal/handler"

	"github.com/DanielTitkov/repertoire/internal/chat"
	"github.com/jfyne/live"
)

func main() {
	log.Println("Starting server...")

	h := handler.NewHandler(
		&app.App{},
		"templates/",
	)

	// Run the server.
	http.Handle("/chat", chat.NewHandler())
	http.Handle("/grid", h.Grid())
	http.Handle("/", h.Home())
	// live scripts
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	// favicon
	http.HandleFunc("/favicon.ico", faviconHandler)
	// serve
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/favicon.ico")
}
