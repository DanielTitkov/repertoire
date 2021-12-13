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
	eventAddTerm            = "addTerm"
	eventRemoveTerm         = "removeTerm"
	eventGenerateTriads     = "generateTriads"
	eventMoveTerm           = "moveTerm"
	eventUpdateTriad        = "updateTriad"
	eventNextTriad          = "nextTriad"
	eventGenerateConstructs = "generateConstructs"
	eventUpdateLinking      = "updateLinking"
	eventGridResult         = "gridResult"
	// params
	paramEmail        = "email"
	paramAge          = "age"
	paramTermID       = "termid"
	paramTriadID      = "triadid"
	paramConstructID  = "constructid"
	paramLinkingValue = "linkingvalue"
	paramMoveTermFrom = "from"
	paramLeftPole     = "leftPole"
	paramRightPole    = "rightPole"
	// params values
	valueMoveTermFromLeft  = "left"
	valueMoveTermFromRight = "right"
)

var funcMap = template.FuncMap{
	"N":     iter.N,
	"Title": strings.Title,
	"sub": func(x, y int) int {
		return x - y
	},
}

type (
	GridModel struct {
		Grid              *domain.Grid
		Session           string
		AddTermError      string
		FormFieldDebounce int // ms
		CurrentTriadID    int
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
					MinTerms:       4,
					MaxTerms:       12,
					TriadMethod:    domain.TriadMethodChoice,
					MinConstructs:  7,
					ConstructSteps: 5,
				},
			),
			Session:           fmt.Sprint(s.Session),
			FormFieldDebounce: 400,
			CurrentTriadID:    0,
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
		h.t+"grid_linking.html",
		h.t+"grid_result.html",
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

		m.CurrentTriadID = 0 // set first triad

		// TODO handle if for some reason there is no triads
		return m, nil
	})

	lvh.HandleEvent(eventMoveTerm, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()

		termID := p.Int(paramTermID)
		triad := m.Grid.GetTriadByIndex(p.Int(paramTriadID))

		switch from := p.String(paramMoveTermFrom); from {
		case valueMoveTermFromLeft:
			err := triad.MoveFromLeft(termID)
			if err != nil {
				return m, err
			}
		case valueMoveTermFromRight:
			err := triad.MoveFromRight(termID)
			if err != nil {
				return m, err
			}
		default:
			return m, fmt.Errorf("unknown term move direction: %s", from)
		}

		return m, nil
	})

	lvh.HandleEvent(eventUpdateTriad, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()

		m.Grid.GetTriadByIndex(m.CurrentTriadID).SetPoles(
			p.String(paramLeftPole),
			p.String(paramRightPole),
		)

		return m, nil
	})

	lvh.HandleEvent(eventNextTriad, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()

		triadID := p.Int(paramTriadID)
		if !(triadID < len(m.Grid.Triads)) {
			return m, nil
		}

		m.CurrentTriadID += 1

		return m, nil
	})

	lvh.HandleEvent(eventGenerateConstructs, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()

		err := m.Grid.GenerateConstructs()
		if err != nil {
			return m, err
		}

		return m, nil
	})

	lvh.HandleEvent(eventUpdateLinking, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()

		m.Grid.Matrix.Set(
			p.Int(paramConstructID),
			p.Int(paramTermID),
			float64(p.Int(paramLinkingValue)),
		) // TODO: probaby move to method

		return m, nil
	})

	lvh.HandleEvent(eventGridResult, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
		m := AssignGridModel(s)
		m.clearErrors()

		err := m.Grid.CalculateResult()
		if err != nil {
			return m, err
		}

		return m, nil
	})

	// lvh.HandleSelf(eventAppendTerm, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
	// 	m := AssignGridModel(s)

	// 	return m, nil
	// })

	return lvh
}
