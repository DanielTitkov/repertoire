package domain

type (
	User struct {
		Email string
		Age   int
	}
	// Grid encapsulates repetoire grid data
	Grid struct {
		Session    string
		Config     GridConfig
		Terms      []Term // TODO: maybe use pointer
		Constructs []Construct
		TermsN     int
	}
	Term struct {
		Title string
	}
	Construct struct {
		Title string
	}
	GridConfig struct {
		MinTerms int
		MaxTerms int
	}
)
