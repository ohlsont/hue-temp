package huetemp

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func indoorPlot(pts []*TemperaturePoint) (*plotter.Line, error) {
	indoorPts := plotter.XYs{}
	for _, p := range pts {
		indoorPts = append(indoorPts, plotter.XY{
			X: float64(p.Time.Unix()),
			Y: p.TempCelsius,
		})
	}
	line, err := plotter.NewLine(indoorPts)
	if err != nil {
		return nil, fmt.Errorf("draw: %w", err)
	}
	line.Color = color.RGBA{G: 255, A: 255}
	return line, nil
}

func outDoorPlot(dataPoints []*TemperaturePoint) (*plotter.Line, error) {
	pts := plotter.XYs{}
	for _, p := range dataPoints {
		pts = append(pts, plotter.XY{
			X: float64(p.Time.Unix()),
			Y: p.OutdoorTemperatureCelsius,
		})
	}
	line, err := plotter.NewLine(pts)
	if err != nil {
		return nil, fmt.Errorf("draw: %w", err)
	}
	line.Color = color.RGBA{R: 255, A: 255}
	return line, nil
}

func Draw(pts []*TemperaturePoint) error {
	p, err := plot.New()
	if err != nil {
		return fmt.Errorf("plot: %w", err)
	}
	p.Title.Text = "Temperature in home"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}
	p.Y.Label.Text = "Temperature"
	p.Y.Tick.Marker = plot.DefaultTicks{}
	indoorLine, err := indoorPlot(pts)
	if err != nil {
		return fmt.Errorf("plot: %w", err)
	}
	outDoorLine, err := outDoorPlot(pts)
	if err != nil {
		return fmt.Errorf("plot: %w", err)
	}
	p.Add(plotter.NewGrid(), indoorLine, outDoorLine)
	p.Legend.Add("indoorLine", indoorLine)
	p.Legend.XOffs = -1 * vg.Centimeter
	p.Legend.Add("outDoorLine", outDoorLine)
	err = p.Save(30*vg.Centimeter, 5*vg.Centimeter, "timeseries.png")
	if err != nil {
		return fmt.Errorf("plot: %w", err)
	}
	return nil
}
