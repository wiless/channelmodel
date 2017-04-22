package main

import (
	"fmt"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/wiless/channelmodel"
	"github.com/wiless/plotutils"
	"github.com/wiless/vlib"
)

var start time.Time

const dpi = 96

var p *plot.Plot

func init() {
	start = time.Now()
	fmt.Println("=========== Started at ", start, "============")
}

func closing() {
	fmt.Println("\n========  End RUNTIME =  ", time.Since(start))
}

var pl CM.RMa
var pl2 CM.UMa
var ds dataframe.DataFrame
var recentlyRESIZED bool
var locations vlib.VectorPos3D

func main() {
	recentlyRESIZED = true
	defer closing()

	fmt.Print("Testing the Channel model\n ")

	pl.Init(30, 1.5, .7)
	pl.ForceLOS = false

	/// Acutal Data Manipulations
	// points := deployment.RectangularNPoints(vlib.Origin3D.Cmplx(), 1000, 500, 30, 700)
	// locations = vlib.FromVectorC(points, 10)
	// updatePathLoss(pl, locations)
	p, _ = plot.New()
	var N = 15300
	var dist, vpl vlib.VectorF

	vlos := make([]bool, N)
	dist.Resize(N)
	vpl.Resize(N)
	cnt := 0
	for ii := 10; ii < N; ii++ {
		d := float64(ii)
		loss, islos, err := pl.PL(d)
		if err == nil {
			dist[cnt] = d
			vpl[cnt] = -loss
			vlos[cnt] = islos
		} else {
			fmt.Println(d, loss, err)
		}
		cnt++
	}
	N = cnt
	dist.Resize(N)
	vpl.Resize(N)

	ds = dataframe.New(series.Ints(vlib.NewSegmentI(1, N)), series.Floats(dist), series.Floats(vpl))
	ds.SetNames([]string{"Index", "distance", "PL"})

	fmt.Print(ds)
	// fmt.Println("Selected ", ds.Subset([]int(vlib.ToVectorI("0:10:100"))))
	// fmt.Println("Random ", ds.Subset(vlib.RandI(30, N)))
	var m vlib.MatrixF
	m.AppendColumn(ds.Col("distance").Float()).AppendColumn(ds.Col("PL").Float()).AppendColumn(vpl.Add(-30))

	pf.Plot(&m)
	// pf.HoldOn()
	// pf.Plot(&m, 0, 2)
	// ds.WriteCSV(os.Stdout)
	// fmt.Printf("Path Loss is : %v  \n", ds)
	// log.Print("Hello")
	// xinfo := []int(vlib.NewSegmentI(650, 10))
	// xxf := ds.Filter(dataframe.F{"distance", series.GreaterEq, 100})

	// log.Print(xxf)

	// plotPL()

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
	// if recentlyRESIZED {
	// 	resize()
	//
	// }
	p.Title.Text = fmt.Sprintf("Time %v", time.Now())

	// canvas.Rotate(math.Pi * rand.Float64() / 6)

}

func newFigure(fname string) {
	// canvas.Scale(3, 4)

	plotutil.AddScatters(p, locations)
	p.Add(plotter.NewGrid())
	p.Title.Text = fmt.Sprintf("Time %v", time.Now().Format(time.Stamp))

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Save(500, 500, fname)

}
