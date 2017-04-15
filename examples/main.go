package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgimg"
	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/kr/pretty"
	"github.com/wiless/cellular/deployment"
	"github.com/wiless/channelmodel"
	"github.com/wiless/vlib"
	"github.com/wiless/x11ui"
)

var start time.Time

const dpi = 96

func init() {
	start = time.Now()
	fmt.Println("=========== Started at ", start, "============")
}

func closing() {
	fmt.Println("\n========  End RUNTIME =  ", time.Since(start))
}

var img *image.RGBA

func main() {

	// p.Save(4*vg.Inch, 4*vg.Inch, "output.png")
	app := x11ui.NewApplication("Application ", 400, 400, true, false)
	child := app.NewChildWindow("Child ", 0, 0, 400, 300)
	// widget := x11ui.NewWidget(app.X(), app.AppWin(), "Plot WINDOW", 30, 30, 350, 250)
	app.DefaultKeys(true)

	defer closing()
	fmt.Print("Testing the Channel model")
	var pl CM.RMa
	pl.Init(30, 5, .7)

	// Initialize plotters
	var p *plot.Plot
	img = image.NewRGBA(image.Rect(0, 0, 400, 300))
	options := vgimg.NewWith(vgimg.UseImage(img))
	canvas := draw.New(options)

	/// Acutal Data Manipulations
	points := deployment.RectangularNPoints(vlib.Origin3D.Cmplx(), 1000, 500, 0, 700)
	locations := vlib.FromVectorC(points, 10)
	pretty.Printf("\n%# v ", pl)
	var pls vlib.VectorF
	var dists vlib.VectorF
	for _, ll := range locations {
		v, _ := pl.LossInDb3D(vlib.Origin3D, ll, .7)

		dists.AppendAtEnd(vlib.Origin3D.DistanceFrom(ll))
		pls.AppendAtEnd(-v)
	}
	ds := dataframe.New(series.Floats(dists), series.Floats(pls))
	ds.SetNames([]string{"distance", "PL"})
	sds := ds.Arrange(dataframe.Sort("distance"))
	fmt.Print(sds)

	// fmt.Printf("Path Loss is : %v  %v \n", CM.IndB(v), CM.Inmtrs(vlib.Origin3D.DistanceFrom(ll)))
	fmt.Printf("Path Loss is : %v  \n", pls)

	/// FUNCITON HANDLERS
	plotPL := func() {
		p, _ = plot.New()
		p.Add(plotter.NewGrid())
		var m vlib.MatrixF
		m.AppendColumn(sds.Col("distance").Float()).AppendColumn(sds.Col("PL").Float())
		// plotutil.AddLinePoints(p, m)
		p.X.Label.Text = "distance"
		p.Y.Label.Text = "Pl (dB)"

		// m.AppendColumn(locations.X()).AppendColumn(locations.Y()).AppendColumn(pls)
		// pb, _ := plotter.NewBubbles(m, 0, 10)
		// p.Add(pb)
		// pb.Color = colorful.LinearRgb(1, 0, 0)

		// histogram
		h, _ := plotter.NewHist(pls, 16)
		h.Normalize(1)
		p.Add(h)

		p.Draw(canvas)
		Xtra(child)
	}
	scatter := func() {
		plotutil.AddScatters(p, locations)
		p.Add(plotter.NewGrid())
		p.Draw(canvas)
		Xtra(child)

	}
	newFigure := func() {
		p, _ = plot.New()
		points = deployment.RectangularNPoints(vlib.Origin3D.Cmplx(), 1000, 500, 0, 700)
		locations = vlib.FromVectorC(points, 10)
		scatter()
	}

	savePlot := func() {
		p.Save(500, 500, "output.png")
		log.Print("Saved into File ", "output.png")
	}

	plotPL()
	app.RegisterKey("c", newFigure)
	app.RegisterKey("p", plotPL)
	app.RegisterKey("s", savePlot)

	Xtra(child)
	app.Show()

}

var ximg *xgraphics.Image

func Xtra(w *x11ui.Window) {

	// log.Printf("Bounds ", img.Bounds(), "window rects", w.Rect)
	// ox, oy := w.Rect.X, w.Rect.Y
	rr := w.ImageRect()
	ximg = xgraphics.NewConvert(w.X(), xgraphics.Scale(img, rr.Dx(), rr.Dy()))
	// ximg = xgraphics.NewConvert(w.X(), xgraphics.Scale(img, rr.Dx(), rr.Dy()))

	ximg.XSurfaceSet(w.Id)
	ximg.XPaintRects(w.Id, image.Rect(0, 0, rr.Dx(), rr.Dy()))
	ximg.XDraw()

}
