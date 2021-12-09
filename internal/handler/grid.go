package handler

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/DanielTitkov/repertoire/internal/domain"
	"github.com/bradfitz/iter"

	"github.com/jfyne/live"
)

const (
	// events
	eventAddTerm        = "addTerm"
	eventRemoveTerm     = "removeTerm"
	eventGenerateTriads = "generateTriads"
	// params
	paramEmail  = "email"
	paramAge    = "age"
	paramTermID = "termid"
)

var funcMap = template.FuncMap{
	"N":     iter.N,
	"Title": strings.Title,
}

type (
	GridModel struct {
		Grid         *domain.Grid
		Session      string
		AddTermError string
	}
)

func (gm *GridModel) clearErrors() {
	gm.AddTermError = ""
}

func AssignGridModel(s *live.Socket) *GridModel {
	m, ok := s.Assigns().(*GridModel)
	if !ok {
		return &GridModel{
			Grid: domain.NewGrid(
				domain.GridConfig{
					MinTerms:    5,
					MaxTerms:    12,
					TriadMethod: domain.TriadMethodForced,
				},
			),
			Session: fmt.Sprint(s.Session),
		}

	}

	return m
}

func (h *Handler) Grid() *live.Handler {
	t := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles(
		h.t+"layout.html",
		h.t+"grid.html",
		h.t+"grid_terms.html",
		h.t+"grid_triads.html",
		h.t+"alerts.html",
	))

	lvh, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}

	// Set the mount function for this handler.
	lvh.Mount = func(ctx context.Context, r *http.Request, s *live.Socket) (interface{}, error) {
		return AssignGridModel(s), nil
	}

	lvh.HandleEvent(eventAddTerm, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()
		termValue := p.String("term")
		err := m.Grid.AddTerm(domain.Term{Title: termValue})
		if err != nil {
			m.AddTermError = err.Error()
		}

		return m, nil
	})

	lvh.HandleEvent(eventRemoveTerm, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()
		termID := p.Int(paramTermID)

		m.Grid.RemoveTermByIndex(termID)

		return m, nil
	})

	lvh.HandleEvent(eventGenerateTriads, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()

		err := m.Grid.GenerateTriads()
		if err != nil {
			fmt.Println(eventGenerateTriads, err)
		}

		m.Grid.Step = domain.GridStepTriads // TODO: maybe move inside method

		return m, nil
	})

	// lvh.HandleEvent(eventUpdateTerm, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
	// 	// FIXME this doesn't work correctly
	// 	m := AssignGridModel(s)
	// 	// fmt.Println(eventUpdateTerm, p)
	// 	return m, nil
	// })

	// lvh.HandleEvent(eventUpdateTerms, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
	// 	m := AssignGridModel(s)
	// 	var newTerms []domain.Term
	// 	for i := range iter.N(m.TermsN) {
	// 		termValue := p.String("term-" + strconv.Itoa(i))
	// 		if termValue != "" {
	// 			newTerms = append(newTerms, domain.Term{Title: termValue})
	// 		}
	// 	}

	// 	m.UpdateTerms(newTerms)

	// 	fmt.Println(eventUpdateTerms, p, m)
	// 	return m, nil
	// })

	// lvh.HandleSelf(eventAppendTerm, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
	// 	m := AssignGridModel(s)

	// 	return m, nil
	// })

	return lvh
}
