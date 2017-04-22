// CM implements the functions & types for implementation of Channel models in 38-901-e00
package CM

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/wiless/cellular/pathloss"
	"github.com/wiless/vlib"
)

// // dist3D returns the 3D distance between two Nodes considering d_out & d_in as in Eq. 7.4.1
// func dist3D(src, dest deployment.Node) (d float64) {
// 	  d3d:=dest.Location.DistanceFrom(src)
//
// }

var rmaDMax float64 = 15000.0 /// max distance supported in RMA for LOS
var rmaH float64 = 5          // Averge building heits in RuralMacro
var rmaW float64 = 20         // Averge building heits in RuralMacro
var rmaHBS float64 = 35
var rmaHUT float64 = 1.5
var rmaNlosMax float64 = 5000

type RMa struct {
	wsettings  pathloss.ModelSetting
	rmaDMax    float64
	dBP        float64 /// Breaking point distance
	c1, c2, c3 float64 /// internal constants
	freq       float64
	ForceLOS   bool
}

func (w *RMa) Set(pathloss.ModelSetting) {

}
func (w RMa) Get() pathloss.ModelSetting {
	return pathloss.ModelSetting{}
}

func (w *RMa) Init(hBS, hUT, fGHz float64) {
	w.freq = fGHz
	w.dBP = 2 * math.Pi * hBS * hUT * fGHz * 1e9 / C
	w.rmaDMax = rmaDMax
	hh := math.Pow(rmaH, 1.72)
	w.c1 = math.Min(0.03*hh, 10)
	w.c2 = math.Min(0.044*hh, 14.77)
	w.c3 = 0.002 * mlog(rmaH)
	w.ForceLOS = false
}

func (w *RMa) SetDMax(dmax float64) {
	w.rmaDMax = dmax
	rmaDMax = dmax
}

// Exported functions MUST implement

func (r RMa) PLbetween(node1, node2 vlib.Location3D) (plDb float64, isNLOS bool, err error) {

	d3d := node1.DistanceFrom(node2)
	d2d := node1.Distance2DFrom(node2)
	var LOS bool = r.ForceLOS
	if !r.ForceLOS {
		LOS = r.IsLOS(d2d)
	}
	plDb, LOS, err = r.PL(d3d)
	return plDb, LOS, err
}

func (r RMa) IsLOS(d2d float64) bool {
	if d2d <= 10 {
		return true
	} else {
		P_LOS := mexp(-(d2d - 10) / 1000)
		if rand.Float64() < P_LOS {
			return true
		} else {
			return false
		}
	}
}

func (r RMa) PLnlos(dist float64) (plDb float64, e error) {

	pldb, err := r.nlos(dist)
	return pldb, err

}

func (r RMa) PLlos(dist float64) (plDb float64, e error) {

	pldb, err := r.los(dist)
	return pldb, err

}

func (r RMa) PL(dist float64) (plDb float64, isNLOS bool, err error) {

	var LOS bool = r.ForceLOS
	if !r.ForceLOS {
		LOS = r.IsLOS(dist)
	}

	if !LOS {
		pldb, err := r.nlos(dist)
		if dist > rmaNlosMax {
			pldb, err := r.los(dist)
			LOS = true
			return pldb, LOS, err
		}
		return pldb, LOS, err
	} else {
		pldb, err := r.los(dist)

		return pldb, LOS, err
	}
}

// non-exported functions internal / private routines
func (r RMa) nlos(dist float64) (plDb float64, e error) {
	freqGHz := r.freq

	var d3d, d2d float64 = dist, dist
	if d2d < 10 {
		return 0, nil
	}

	if 10 <= d2d && d2d <= 5000 {
		loss1, _ := r.los(d3d)
		loss2 := 161.04 - 7.1*mlog(rmaW) + 7.5*mlog(rmaH) - (24.37-3.7*math.Pow(rmaH/rmaHBS, 2))*mlog(rmaHBS) + (43.42-3.1*mlog(rmaHBS))*(mlog(d3d)-3) + 20*mlog(freqGHz) - mpow(3.2*(mlog(11.75*rmaHUT)), 2) - 4.97

		return max(loss1, loss2), nil
	} else {

		return math.NaN(), fmt.Errorf("Distance %d not supported in this model ", dist)
	}
}

func (r RMa) los(dist float64) (plDb float64, e error) {
	freqGHz := r.freq
	var d3d, d2d float64 = dist, dist

	if d2d < 10 {
		return 0, fmt.Errorf("Should not be issue ")
	}
	if 10 <= d2d && d2d <= r.dBP {
		loss, _ := r.p1(d3d, freqGHz)

		return loss, nil
	} else if d2d > r.dBP && d2d <= r.rmaDMax {
		loss, _ := r.p2(d3d, freqGHz)
		return loss, nil
	} else {
		// return math.NaN(), fmt.Errorf("Unsupported distance %d for LOS ", dist)
		return 0, fmt.Errorf("Unsupported distance %d for LOS ", dist)
	}
}

func (r *RMa) losNodes(src, dest vlib.Location3D) (plDb float64, valid bool) {
	freqGHz := r.freq
	d3d := src.DistanceFrom(dest)
	d2d := src.Distance2DFrom(dest)

	if 10 <= d2d && d2d <= r.dBP {
		loss, _ := r.los(d3d)
		return loss, true
	} else if d2d > r.dBP && d2d <= r.rmaDMax {
		loss, _ := r.p2(d3d, freqGHz)
		plDb, _ = r.los(d3d)
		plDb += 40 * mlog(d3d/r.dBP)
		return loss, valid
	} else {
		log.Printf("\nDistance not supported in this model")
		return 0, false
	}

}

func (r *RMa) p1(d3d, freqGHz float64) (plDb float64, valid bool) {

	plDb = 20*mlog(40*pi*d3d*freqGHz/3) + r.c1*mlog(d3d) - r.c2 + r.c3*d3d
	return plDb, true
}

func (r *RMa) p2(d3d, freqGHz float64) (plDb float64, valid bool) {

	plDb, valid = r.p1(r.dBP, freqGHz)
	plDb += 40 * mlog(d3d/r.dBP)
	return plDb, true
}
