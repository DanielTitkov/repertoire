package domain

type (
	User struct {
		Email string
		Age   int
	}
	// Grid encapsulates repetoire grid data
	Grid struct {
		Config     GridConfig
		Terms      []Term // TODO: maybe use pointer
		Constructs []Construct
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
