// CM implements the functions & types for implementation of Channel models in 38-901-e00
package CM

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/wiless/cellular/deployment"
	"github.com/wiless/cellular/pathloss"
	"github.com/wiless/vlib"
)

//  dist3D returns the 3D distance between two Nodes considering d_out & d_in as in Eq. 7.4.1

var umaDMax float64 = 21000.0 /// max distance supported in uma for LOS
var umaH float64 = 5          // Averge building heits in RuralMacro
var umaW float64 = 20         // Averge building heits in RuralMacro
var umaHBS float64 = 35
var umaHUT float64 = 5
var umaNlosMax float64 = 5000

type UMa struct {
	wsettings  pathloss.ModelSetting
	umaDMax    float64
	dBP        float64 /// Breaking point distance
	c1, c2, c3 float64 /// internal constants
	freq       float64
	ForceLOS   bool
}

func (w *UMa) Set(pathloss.ModelSetting) {

}
func (w UMa) Get() pathloss.ModelSetting {
	return pathloss.ModelSetting{}
}

func (w *UMa) Init(hBS, hUT, fGHz float64) {

	// TODO to be fixed from the Note 1 of Ref. document Table 7.4.1-1
	/*
		Breakpoint distance d'BP = 4 h'BS h'UT fc/c, where fc is the centre frequency in Hz, c = 3.0×108 m/s is the propagation velocity in free space, and h'BS and h'UT are the effective antenna heights at the BS and the UT, respectively. The effective antenna heights h'BS and h'UT are computed as follows: h'BS = hBS – hE, h'UT = hUT – hE, where hBS and hUT are the actual antenna heights, and hE is the effective environment height. For UMi hE = 1.0m. For UMa hE=1m with a probability equal to 1/(1+C(d2D, hUT)) and chosen from a discrete uniform distribution uniform(12,15,...,(hUT-1.5)) otherwise. With C(d2D, hUT) given by
		⎧0 ,hUT <13m C(d2D,hUT)=⎪⎨⎛hUT −13⎞1.5
		⎪⎩⎜⎝ 10 ⎟⎠g(d2D),13m≤hUT≤23m  ,
		 where
	*/
	w.freq = fGHz
	w.dBP = 4 * hBS * hUT * fGHz * 1e9 / C
	w.umaDMax = umaDMax
	hh := math.Pow(umaH, 1.72)
	w.c1 = math.Min(0.03*hh, 10)
	w.c2 = math.Min(0.044*hh, 14.77)
	w.c3 = 0.002 * mlog(umaH)
	w.ForceLOS = false
}

func (r UMa) PLbetween(node1, node2 vlib.Location3D) (plDb float64, isNLOS bool, err error) {

	d3d := node1.DistanceFrom(node2)
	d2d := node1.Distance2DFrom(node2)
	var LOS bool = r.ForceLOS
	if !r.ForceLOS {
		LOS = r.IsLOS(d2d)
	}
	plDb, LOS, err = r.PL(d3d)
	return plDb, LOS, err
}
func (r UMa) IsLOS(d2d float64) bool {
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

func (r UMa) PLnlos(dist float64) (plDb float64, e error) {

	pldb, err := r.nlos(dist)
	return pldb, err

}

func (r UMa) PLlos(dist float64) (plDb float64, e error) {

	pldb, err := r.los(dist)
	return pldb, err

}

func (r UMa) PL(dist float64) (plDb float64, isNLOS bool, err error) {

	var LOS bool = r.ForceLOS
	if !r.ForceLOS {
		LOS = r.IsLOS(dist)
	}

	if !LOS {
		pldb, err := r.nlos(dist)
		if dist > umaNlosMax {
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

func (r UMa) nlos(dist float64) (plDb float64, e error) {
	freqGHz := r.freq

	var d3d, d2d float64 = dist, dist
	if d2d < 10 {
		return 0, nil
	}

	if 10 <= d2d && d2d <= 5000 {
		loss1, _ := r.los(d3d)
		loss2 := 161.04 - 7.1*mlog(umaW) + 7.5*mlog(umaH) - (24.37-3.7*math.Pow(umaH/umaHBS, 2))*mlog(umaHBS) + (43.42-3.1*mlog(umaHBS))*(mlog(d3d)-3) + 20*mlog(freqGHz) - mpow(3.2*(mlog(11.75*umaHUT)), 2) - 4.97

		return max(loss1, loss2), nil
	} else {

		return math.NaN(), fmt.Errorf("Distance %d not supported in this model ", dist)
	}
}

func (r UMa) los(dist float64) (plDb float64, e error) {
	freqGHz := r.freq
	var d3d, d2d float64 = dist, dist

	if d2d < 10 {
		return 0, nil
	}
	if 10 <= d2d && d2d <= r.dBP {
		loss, _ := r.p1(d3d, freqGHz)

		return loss, nil
	} else if d2d > r.dBP && d2d <= r.umaDMax {
		loss, _ := r.p2(d3d, freqGHz)
		return loss, nil
	} else {

		return math.NaN(), fmt.Errorf("Unsupported distance %d for LOS ", dist)
	}
}

func (r *UMa) losNodes(src, dest vlib.Location3D) (plDb float64, valid bool) {
	freqGHz := r.freq
	d3d := src.DistanceFrom(dest)
	d2d := src.Distance2DFrom(dest)

	if 10 <= d2d && d2d <= r.dBP {
		loss, _ := r.los(d3d)
		return loss, true
	} else if d2d > r.dBP && d2d <= r.umaDMax {
		loss, _ := r.p2(d3d, freqGHz)
		plDb, _ = r.los(d3d)
		plDb += 40 * mlog(d3d/r.dBP)
		return loss, valid
	} else {
		log.Printf("\nDistance not supported in this model")
		return 0, false
	}

}

func (r *UMa) p1(d3d, freqGHz float64) (plDb float64, valid bool) {
	plDb = 28.0*22*mlog(d3d) + 20*mlog(freqGHz) + r.c1*mlog(d3d) - r.c2 + r.c3*d3d
	return plDb, true
}

func (r *UMa) p2(d3d, freqGHz float64) (plDb float64, valid bool) {

	plDb, valid = r.p1(r.dBP, freqGHz)
	plDb += 40 * mlog(d3d/r.dBP)
	return plDb, true
}
func (r *UMa) LossInDbNodes(txnode, rxnode deployment.Node, freqGHz float64) (plDb float64, valid bool) {

	return 0, true
}
