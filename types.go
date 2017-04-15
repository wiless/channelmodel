// CM implements the functions & types for implementation of Channel models in 38-901-e00
package CM

import (
	"fmt"
	"log"
	"math"

	clr "github.com/fatih/color"
	"github.com/wiless/cellular/deployment"
	"github.com/wiless/cellular/pathloss"
	"github.com/wiless/vlib"
)

// // dist3D returns the 3D distance between two Nodes considering d_out & d_in as in Eq. 7.4.1
// func dist3D(src, dest deployment.Node) (d float64) {
// 	  d3d:=dest.Location.DistanceFrom(src)
//
// }
const C float64 = 3.0e8         // Speed of light
const rmaDMax float64 = 10000.0 /// max distance supported in RMA for LOS
const rmaH float64 = 20         // Averge building heits in RuralMacro
var mlog = math.Log10

const pi = math.Pi

func IndB(v float64) string {
	return fmt.Sprintf("%.3f%s", v, clr.CyanString("dB"))
}
func InGHz(v float64) string {
	return fmt.Sprintf("%.3f%s", v, clr.GreenString("GHz"))
}
func Inmtrs(v float64) string {
	return fmt.Sprintf("%.3f%s", v, clr.GreenString("m "))
}

type Frequency float64

func (f *Frequency) fromGHz(fghz Frequency) Frequency {
	*f = 1.0e9 * fghz
	return *f
}

func (f *Frequency) GHz(fghz Frequency) {
	*f = 1.0e-9 * fghz
}

type RMa struct {
	wsettings  pathloss.ModelSetting
	rmaDMax    float64
	dBP        float64 /// Breaking point distance
	c1, c2, c3 float64 /// internal constants
}

func (w *RMa) Set(pathloss.ModelSetting) {

}
func (w RMa) Get() pathloss.ModelSetting {
	return pathloss.ModelSetting{}
}

func (w *RMa) Init(hBS, hUT, fGHz float64) {
	w.dBP = 2 * math.Pi * hBS * hUT * fGHz * 1e9 / C
	w.rmaDMax = rmaDMax
	hh := math.Pow(rmaH, 1.72)
	w.c1 = math.Min(0.03*hh, 10)
	w.c2 = math.Min(0.044*hh, 14.77)
	w.c3 = 0.002 * mlog(rmaH)
}

// type Model interface {
// 	Set(ModelSetting)
// 	Get() ModelSetting
// 	LossInDbNodes(txnode, rxnode deployment.Node, freqGHz float64) (plDb float64, valid bool)
// 	LossInDb3D(txnode, rxnode vlib.Location3D, freqGHz float64) (plDb float64, valid bool)
// }

func (r *RMa) LossInDb(dist float64, freqGHz float64) (plDb float64, valid bool) {
	// FreqMHz := freqGHz * 1.0e3                 // Frequency is in MHz
	// distance := src.DistanceFrom(dest) / 1.0e3 // Convert to km (most equations have d in km)
	// var result float64
	// result = -1
	// result = 46.3 + 33.9*math.Log10(FreqMHz) - 13.82*math.Log10(src.Z) - a + (44.9-6.55*math.Log10(src.Z))*math.Log10(distance) + 3
	var d3d, d2d float64 = dist, dist
	if 10 <= d2d && d2d <= r.dBP {
		loss, valid := r.LOSp1(d3d, freqGHz)
		return loss, valid
	} else if d2d > r.dBP && d2d <= r.rmaDMax {
		loss, valid := r.LOSp2(d3d, freqGHz)
		return loss, valid
	} else {
		log.Printf("\nDistance not supported in this model")
		return 0, false
	}

}

func (r *RMa) LossInDb3D(src, dest vlib.Location3D, freqGHz float64) (plDb float64, valid bool) {
	// FreqMHz := freqGHz * 1.0e3                 // Frequency is in MHz
	// distance := src.DistanceFrom(dest) / 1.0e3 // Convert to km (most equations have d in km)
	// var result float64
	// result = -1
	// result = 46.3 + 33.9*math.Log10(FreqMHz) - 13.82*math.Log10(src.Z) - a + (44.9-6.55*math.Log10(src.Z))*math.Log10(distance) + 3
	d3d := src.DistanceFrom(dest)
	d2d := src.Distance2DFrom(dest)
	if 10 <= d2d && d2d <= r.dBP {
		loss, valid := r.LOSp1(d3d, freqGHz)
		return loss, valid
	} else if d2d > r.dBP && d2d <= r.rmaDMax {
		loss, valid := r.LOSp2(d3d, freqGHz)
		return loss, valid
	} else {
		log.Printf("\nDistance not supported in this model")
		return 0, false
	}

}

func (r *RMa) LOSp1(d3d, freqGHz float64) (plDb float64, valid bool) {
	log.Println("\nWithin Break Point distance ", d3d)
	plDb = 20*mlog(40*pi*d3d*freqGHz/3) + r.c1*mlog(d3d) - r.c2 + r.c3*d3d
	return plDb, true
}

func (r *RMa) LOSp2(d3d, freqGHz float64) (plDb float64, valid bool) {
	log.Println("\n>>> BP distance")
	plDb, valid = r.LOSp1(r.dBP, freqGHz)
	plDb += 40 * mlog(d3d/r.dBP)
	return plDb, true
}
func (r *RMa) LossInDbNodes(txnode, rxnode deployment.Node, freqGHz float64) (plDb float64, valid bool) {

	return 0, true
}
