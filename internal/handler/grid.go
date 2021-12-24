package handler

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/DanielTitkov/repertoire/internal/charts"

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
	paramTermValue    = "term"
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
	"Reverse": domain.ReverseSliceTerms,
}

type (
	GridModel struct {
		Grid              *domain.Grid
		UpdateValue       int
		Session           string
		AddTermError      string
		FormFieldDebounce int // ms
		CurrentTriadID    int
		Charts            struct {
			TermsCorr      []charts.Heatmap
			ConstructsCorr []charts.Heatmap
		}
		TermsCorrs      []Corr
		ConstructsCorrs []Corr
	}
	Corr struct {
		Left  string
		Right string
		Value float64
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
					MinTerms:       7,
					MaxTerms:       12,
					TriadMethod:    domain.TriadMethodChoice,
					MinConstructs:  7,
					ConstructSteps: 4,
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
		err := m.Grid.AddTerm(domain.Term{Title: p.String(paramTermValue)})
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

		// correlation heatmaps
		m.Charts.TermsCorr = makeTermsHeatmapData(m.Grid)
		m.Charts.ConstructsCorr = makeConstructsHeatmapData(m.Grid)

		// correlations
		m.ConstructsCorrs = makeConstructsCorrs(m.Grid)
		m.TermsCorrs = makeTermsCorrs(m.Grid)

		m.UpdateValue += 1

		return m, nil
	})

	// lvh.HandleSelf(eventAppendTerm, func(ctx context.Context, s *live.Socket, p live.Params) (interface{}, error) {
	// 	m := AssignGridModel(s)

	// 	return m, nil
	// })

	return lvh
}

func makeTermsHeatmapData(g *domain.Grid) []charts.Heatmap {
	var titles []string
	for _, t := range g.Terms {
		titles = append(titles, t.Title)
	}

	var data [][]float64
	for i := range g.Terms {
		var row []float64
		for j := range g.Terms {
			row = append(row, g.Analysis.TermsCorrMatrix.At(i, j))
		}
		data = append(data, row)
	}

	return charts.NewHeatmapData(
		data,
		titles,
		titles,
	)
}

func makeConstructsHeatmapData(g *domain.Grid) []charts.Heatmap {
	var titles []string
	for _, c := range g.Constructs {
		titles = append(titles, c.Title())
	}

	var data [][]float64
	for i := range g.Constructs {
		var row []float64
		for j := range g.Constructs {
			row = append(row, g.Analysis.ConstructsCorrMatrix.At(i, j))
		}
		data = append(data, row)
	}

	return charts.NewHeatmapData(
		data,
		titles,
		titles,
	)
}

func makeConstructsCorrs(g *domain.Grid) []Corr {
	var corrs []Corr

	for i := 0; i < len(g.Constructs); i++ {
		for j := 0; j < (i + 1); j++ {
			if g.Constructs[i].Title() != g.Constructs[j].Title() {
				corrs = append(corrs, Corr{
					Left:  g.Constructs[i].Title(),
					Right: g.Constructs[j].Title(),
					Value: math.Round(g.Analysis.ConstructsCorrMatrix.At(i, j)*1000) / 1000,
				})
			}
		}
	}

	sort.Slice(corrs, func(i, j int) bool {
		return corrs[i].Value > corrs[j].Value
	})

	return corrs
}

func makeTermsCorrs(g *domain.Grid) []Corr {
	var corrs []Corr

	for i := 0; i < len(g.Terms); i++ {
		for j := 0; j < (i + 1); j++ {
			if g.Terms[i].Title != g.Terms[j].Title {
				corrs = append(corrs, Corr{
					Left:  g.Terms[i].Title,
					Right: g.Terms[j].Title,
					Value: math.Round(g.Analysis.TermsCorrMatrix.At(i, j)*1000) / 1000,
				})
			}
		}
	}

	sort.Slice(corrs, func(i, j int) bool {
		return corrs[i].Value > corrs[j].Value
	})

	return corrs
}
