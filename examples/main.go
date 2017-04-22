package main

import (
	"fmt"
	"time"

	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/wiless/channelmodel"
	"github.com/wiless/plotutils"
	"github.com/wiless/vlib"
)

var start time.Time

const dpi = 96

func init() {
	start = time.Now()
	fmt.Println("=========== Started at ", start, "============")
}

func closing() {
	pf.CloseAll()
	fmt.Println("\n========  End RUNTIME =  ", time.Since(start))
}

var pl CM.RMa
var pl2 CM.UMa
var ds dataframe.DataFrame

var locations vlib.VectorPos3D

func main() {

	defer closing()
	pf.StartX()

	fmt.Print("Testing the Channel model\n ")

	pl.Init(30, 1.5, .7)
	pl.ForceLOS = false
	pl.SetDMax(15000)
	/// Acutal Data Manipulations
	// points := deployment.RectangularNPoints(vlib.Origin3D.Cmplx(), 1000, 500, 30, 700)
	// locations = vlib.FromVectorC(points, 10)
	// updatePathLoss(pl, locations)

	var N = 10000
	var dist, vpl vlib.VectorF

	vlos := make([]bool, N)
	dist.Resize(N)
	vpl.Resize(N)
	cnt := 0
	for ii := 10; ii < N; ii += 100 {
		d := float64(ii)
		loss, islos, err := pl.PL(d)
		_ = err
		// if err == nil {
		dist[cnt] = d
		vpl[cnt] = -loss
		vlos[cnt] = islos
		// } else {
		cnt++
		if err != nil {
			fmt.Print(err)
		}
		// 	fmt.Println(d, loss, err)
		// }

	}
	N = cnt
	dist.Resize(N)
	vpl.Resize(N)

	ds = dataframe.New(series.Floats(dist), series.Floats(vpl))
	ds.SetNames([]string{"distance", "PL"})
	fmt.Println(ds.Dims())

	// ds.Arrange(dataframe.Order{Colname: "PL"})
	// ds = ds.Filter(dataframe.F{"PL", series.Neq, 0})

	var m vlib.MatrixF
	m.AppendColumn(ds.Col("distance").Float()).AppendColumn(ds.Col("PL").Float()).AppendColumn(vpl.Add(-10))

	pf.Fig("Kavish Plot")
	pf.Plot(&m)
	//
	// pf.SetXlabel("Distance (m)")
	// pf.SetYlabel(`PL ($dB_m$) `)
	pf.HoldOff()
	//
	pf.Plot(&m, 0, 2)
	// pf.ShowX11()
	// ds.WriteCSV(os.Stdout)
	pf.Wait()
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
