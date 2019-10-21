package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	// "github.com/kniren/gota/dataframe"
	// "github.com/kniren/gota/series"
	"github.com/wiless/cellular/deployment"
	"github.com/wiless/channelmodel"
	"github.com/wiless/vlib"
	"github.com/wiless/x11ui"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
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

// Initialize plotters
// var p *plot.Plot
var canvas draw.Canvas
var plotwin *x11ui.Window

var pl CM.RMa
var ds dataframe.DataFrame

var recentlyRESIZED bool

var resize = func() {
	pwin, _ := plotwin.Parent()
	appgeo, _ := pwin.Geometry()

	// appgeo, _ := app.AppWin().Geometry()
	log.Print(appgeo.Width(), appgeo.Height())

	plotwin.Window.Resize(appgeo.Width(), appgeo.Height())
	img = plotwin.CreateRawImage(0, 0, appgeo.Width(), appgeo.Height())
	recentlyRESIZED = false
	// canvas = draw.New(vgimg.NewWith(vgimg.UseImage(img)))
	// plotPL()
}

func main() {
	recentlyRESIZED = true
	defer closing()

	// p.Save(4*vg.Inch, 4*vg.Inch, "output.png")
	app := x11ui.NewApplication("Application ", 400, 300, true, false)
	app.Debug = false
	app.DefaultKeys(true)
	plotwin = app.NewChildWindow("child", 0, 0, 400, 300)

	img = plotwin.CreateRawImage(0, 0, 400, 300)

	// xevent.ResizeRequestFun(
	// 	func(X *xgbutil.XUtil, e xevent.ResizeRequestEvent) {
	// 		log.Print("1.========Outside plot win Received ", e)
	// 		// plotwin.ReSize(int(e.Width), int(e.Height))
	// 		plotwin.Window.Resize(int(e.Width), int(e.Height))
	//
	// 	}).Connect(app.X(), app.AppWin().Id)
	//

	xevent.ConfigureNotifyFun(
		func(p *xgbutil.XUtil, e xevent.ConfigureNotifyEvent) {
			log.Printf("2. PLTWN ceived CONFIGNOTIFICATION ", e)
			// plotwin.Window.Resize(int(e.Width), int(e.Height))
			// plotwin.ReDrawImage()
			recentlyRESIZED = true

		}).Connect(app.X(), app.AppWin().Id)

	app.RegisterKey("r", resize)

	fmt.Print("Testing the Channel model")
	pl.Init(5)
	// pl.Init(30, 1.5, .7)
	pl.ForceLOS = false

	/// Acutal Data Manipulations
	points := deployment.RectangularNPoints(vlib.Origin3D.Cmplx(), 1000, 500, 30, 700)
	locations := vlib.FromVectorC(points, 10)
	updatePathLoss(pl, locations)
	// fmt.Printf("Path Loss is : %v  \n", ds)
	// log.Print("Hello")
	// xinfo := []int(vlib.NewSegmentI(650, 10))
	// xxf := ds.Filter(dataframe.F{"distance", series.GreaterEq, 100})

	// log.Print(xxf)

	plotPL()

	// p.Save(500, 500, "output.png")
	app.RegisterKey("c", newFigure)
	app.RegisterKey("p", plotPL)
	app.RegisterKey("s", savePlot)
	plotwin.ReDrawImage()
	app.Show()

}

func updatePathLoss(pl CM.RMa, locations vlib.VectorPos3D) {

	var pls vlib.VectorF
	var dists vlib.VectorF
	// var los series.Series
	los := series.New([]bool{}, series.Bool, "ISLOS")

	for _, ll := range locations {
		v, islos, _ := pl.PLbetween(vlib.Origin3D, ll)
		dists.AppendAtEnd(vlib.Origin3D.DistanceFrom(ll))
		pls.AppendAtEnd(-v)
		los.Append(islos)
	}

	ds = dataframe.New(series.New(dists, series.Float, "distance"), series.New(pls, series.Float, "PL"), los, series.Floats(locations.X()), series.Floats(locations.Y())).Arrange(dataframe.Sort("distance"))

}

/// FUNCITON HANDLERS
func plotPL() {
	p, _ := plot.New()
	p.Add(plotter.NewGrid())

	// var m vlib.MatrixF
	// m.AppendColumn(ds.Col("distance").Float()).AppendColumn(ds.Col("PL").Float())
	// plotutil.AddLines(p, m)
	// p.X.Label.Text = "distance"
	// p.Y.Label.Text = "Pl (dB)"

	// m.AppendColumn(locations.X()).AppendColumn(locations.Y()).AppendColumn(pls)
	// pb, _ := plotter.NewBubbles(m, 0, 10)
	// p.Add(pb)
	// pb.Color = colorful.LinearRgb(1, 0, 0)

	// histogram
	var v vlib.VectorF = ds.Col("PL").Float()
	h, _ := plotter.NewHist(v, 16)
	// h.Normalize(1)
	p.Add(h)
	// p.Y.Max = 1
	if recentlyRESIZED {
		resize()

	}
	p.Title.Text = fmt.Sprintf("Time %v", time.Now())
	canvas = draw.New(vgimg.NewWith(vgimg.UseImage(img)))
	// canvas.Rotate(math.Pi * rand.Float64() / 6)

	p.Save(500, 500, "output.png")
	p.Draw(canvas)
	plotwin.ReDrawImage()

}

func newFigure() {
	// canvas.Scale(3, 4)
	p, _ := plot.New()
	points := deployment.RectangularNPoints(vlib.Origin3D.Cmplx(), 1000, 500, 30, 700)
	locations := vlib.FromVectorC(points, 10)
	updatePathLoss(pl, locations)
	plotutil.AddScatters(p, locations)
	p.Add(plotter.NewGrid())
	p.Title.Text = fmt.Sprintf("Time %v", time.Now().Format(time.Stamp))
	if recentlyRESIZED {
		resize()

	}
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	canvas = draw.New(vgimg.NewWith(vgimg.UseImage(img)))
	canvas.Translate(vg.Point{10, 10})
	canvas.Scale(.9, .9)
	p.Save(500, 500, "output.png")
	p.Draw(canvas)
	plotwin.ReDrawImage()
}

func savePlot() {
	var p *plot.Plot
	p.Save(500, 500, "output.png")
	log.Print("Saved into File ", "output.png")
}
