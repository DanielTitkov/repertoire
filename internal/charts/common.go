package charts

type (
	Title struct {
		Text string  `json:"text"`
		Font Font    `json:"font"`
		Xref string  `json:"xref"`
		X    float64 `json:"x"`
	}

	Font struct {
		Family string `json:"family"`
		Size   int    `json:"size"`
	}
)

func NewTitle(text string) Title {
	return Title{
		Text: text,
		Xref: "paper",
		X:    0.05,
		Font: Font{
			Family: "Ubuntu, monospace",
			Size:   24,
		},
	}
}
