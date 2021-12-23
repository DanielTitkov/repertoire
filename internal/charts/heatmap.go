package charts

const (
	typeHeatmap = "heatmap"
	colorScale  = "Portland"
)

type Heatmap struct {
	Z           [][]float64 `json:"z"`
	X           []string    `json:"x"`
	Y           []string    `json:"y"`
	Type        string      `json:"type"`
	HoverOnGaps bool        `json:"hoverongaps"`
	ColorScale  string      `json:"colorscale"`
}

func NewHeatmap(
	z [][]float64,
	x []string,
	y []string,
) Heatmap {
	return Heatmap{
		Z:           z,
		X:           x,
		Y:           y,
		Type:        typeHeatmap,
		HoverOnGaps: false,
		ColorScale:  colorScale,
	}
}

func NewHeatmapData(
	z [][]float64,
	x []string,
	y []string,
) []Heatmap {
	return []Heatmap{NewHeatmap(z, x, y)}
}
