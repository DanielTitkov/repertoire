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
		Triads     []*Triad
		Step       string
	}
	Term struct {
		Title string
	}
	Triad struct {
		LeftTerms  []*Term
		RightTerms []*Term
		LeftPole   string
		RightPole  string
		Step       string
	}
	Construct struct {
		Title string
	}
	GridConfig struct {
		TriadMethod   string
		MinTerms      int
		MaxTerms      int
		MinConstructs int
	}
)
