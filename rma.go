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

const rmaDMax = 10000.0 /// max distance supported in RMA for LOS
const rmaH float64 = 5  // Averge building heits in RuralMacro
const rmaW float64 = 20 // Averge road width in RuralMacro
const rmaHBS float64 = 35
const rmaHUT = 1.5

// var rmaNlosMax float64 = 5000

//wraps the interface for supporting deployment link

type RMa struct {
	*pathloss.ModelSetting
	dBP            float64 /// Breaking point distance
	c1, c2, c3     float64 /// internal constants for LOS
	c4, c5, c6, c7 float64 /// internal constants for NLOS
	ForceLOS       bool
	ForceNLOS      bool
	isOK           bool
	Extended       bool
	rmaNlosMax     float64
	streetW        float64
	nlossOffset    float64 // Offset for LMLC
}

// ForcesLOS for all
func (r *RMa) ForceAllLOS(f bool) {
	r.ForceLOS = f
	if f {
		r.ForceNLOS = false
	}
}

// ForcesLOS for all
func (r *RMa) ForceAllNLOS(f bool) {
	r.ForceNLOS = f
	if f {
		r.ForceLOS = false
	}
}

// Returns a default RMA Model settings
func RMADefault() *pathloss.ModelSetting {
	ms := pathloss.NewModelSetting()
	ms.SetFGHz(.7)
	ms.CutOffDistance = 21000 /// updated with recent ITU doc
	ms.Name = "RMa"
	ms.AddParam("HBS", rmaHBS).AddParam("HUT", rmaHUT)
	return ms
}

func (r *RMa) Check() {
	if !r.isOK {
		log.Panic("RMa Model not initialized, Call .Init() ")
	}
}

func (w *RMa) Set(ms *pathloss.ModelSetting) {
	if ms.CutOffDistance == 0 {
		ms.CutOffDistance = rmaDMax
	}
	w.rmaNlosMax = 5000

	// copy(w.ModelSetting, *ms)
	w.ModelSetting = ms
	hBS, hUT := w.Value("hBS"), w.Value("hUT")
	fGHz := w.FGHz()
	if fGHz == 0 {
		w.isOK = false
		log.Print("RMA ModelSettings Frequency not set !!")
	}
	if hBS == 0 || hUT == 0 {
		log.Print("Could not find paramters hBS / hUT in the setting, Setting to Default 35.0m & 1.5m")
		hBS = rmaHBS
		hUT = rmaHUT
	} else {
		log.Print("Good found values !! ")
	}
	hh := math.Pow(rmaH, 1.72)
	w.c1 = math.Min(0.03*hh, 10)
	w.c2 = math.Min(0.044*hh, 14.77)
	w.c3 = 0.002 * mlog(rmaH)

	w.streetW = rmaW
	// nlos constants
	w.c4 = 161.04 - 7.1*mlog(w.streetW) + 7.5*mlog(rmaH)
	w.c5 = -(24.37 - 3.7*math.Pow(rmaH/hBS, 2)) * mlog(hBS)
	w.c6 = (43.42 - 3.1*mlog(hBS))
	w.c7 = 20*mlog(fGHz) - (3.2*(math.Pow(mlog(11.75*hUT), 2)) - 4.97)

	w.ForceLOS = false
	w.ForceNLOS = false
	w.dBP = w.BPDistance()

	w.nlossOffset = 12.0 // 12dB offset if its LMLC

	w.isOK = true
}

//DMax returns the maximum supported distance
func (r RMa) DMax() float64 {
	return r.CutOffDistance
}

func (r *RMa) SetExtended(isext bool) {
	r.Extended = isext
	if isext {
		r.SetDMax(21000)     // based on update for LMLC
		r.SetNLosDMax(21000) // based on update for LMLC
	} else {
		r.SetDMax(10000)     // based on update for LMLC
		r.SetNLosDMax(10000) // based on update for LMLC
	}
}

func (r *RMa) SetStreetW(ww float64) {
	r.streetW = ww
}
func (w RMa) Get() *pathloss.ModelSetting {
	return w.ModelSetting
}

// Initializes with the default RMA modelsettings at fGHz frequency
func (w *RMa) Init(fGHz float64) {
	w.Set(RMADefault().SetFGHz(fGHz))
}

func (w *RMa) SetDMax(dmax float64) {
	w.Check()
	w.CutOffDistance = dmax
}

func (w *RMa) SetNLosDMax(dmax float64) {
	w.Check()

	w.rmaNlosMax = dmax
}

func (r RMa) BPDistance() float64 {
	hBS, hUT := r.Value("hBS"), r.Value("hUT")
	log.Print(hBS, hUT)
	r.dBP = 2 * math.Pi * hBS * hUT * r.FGHz() * 1e9 / C
	// r.dBP=7667

	return r.dBP
}

// Exported functions MUST implement
func (r RMa) IsSupported(fghz float64) bool {
	r.Check()
	// 30GHz According to Note 2 of Table 7.4.1-1 Path Loss (fgHz in GHz)
	if r.FGHz() < 30 && r.FGHz() > 0.5 {
		return true
	}
	return false
}

func (r RMa) PLbetween(node1, node2 vlib.Location3D) (plDb float64, isNLOS bool, err error) {
	if !r.isOK {
		log.Panicln("RMA: Model not Initialized ...")
	}
	d3d := node1.DistanceFrom(node2)
	d2d := node1.Distance2DFrom(node2)

	var LOS bool = r.ForceLOS
	if !r.ForceLOS && !r.ForceNLOS {
		LOS = r.IsLOS(d2d)
	}
	plDb, LOS, err = r.PL(d3d)
	if LOS && err != nil {
		log.Printf("PL  ISLOS=%v , ERROR=%v , %v to %v  ", LOS, err, node1, node2)
	}
	return plDb, LOS, err
}

func (r RMa) IsLOS(d2d float64) bool {
	if !r.isOK {
		log.Panicln("RMA: Model not Initialized ...")
	}

	if d2d <= 10 {
		return true
	} else {

		// Model does not suppor NLOS > rmaNlosMax (=5000)
		if d2d > r.rmaNlosMax {
			return true
		}

		P_LOS := mexp(-(d2d - 10) / 1000)
		if rand.Float64() < P_LOS {
			return true
		} else {
			return false
		}
	}
}

func (r RMa) PLosPDF(d2d float64) float64 {

	if d2d <= 10 {
		return 1
	} else {

		// Model does not suppor NLOS > rmaNlosMax (=5000)
		if d2d > r.rmaNlosMax {
			return 1
		}

		P_LOS := mexp(-(d2d - 10) / 1000)
		return P_LOS
	}
}
func (r RMa) PLnlos(dist float64) (plDb float64, e error) {
	r.Check()

	pldb, err := r.nlos(dist)

	return pldb, err

}

func (r RMa) PLlos(dist float64) (plDb float64, e error) {
	r.Check()

	pldb, err := r.los(dist)

	return pldb, err

}

func (r RMa) PL(dist float64) (plDb float64, isNLOS bool, err error) {
	r.Check()

	var LOS bool = r.ForceLOS
	if !r.ForceLOS && !r.ForceNLOS {
		if dist < r.rmaNlosMax {
			LOS = r.IsLOS(dist)
		}

	}

	if !LOS {
		pldb, err := r.nlos(dist)

		return pldb, LOS, err

	} else {
		pldb, err := r.los(dist)

		return pldb, LOS, err
	}
}

// non-exported functions internal / private routines
func (r RMa) nlos(dist float64) (plDb float64, e error) {
	freqGHz := r.FGHz()
	// log.Print("NLOS Freq is ",freqGHz)
	var d3d, d2d float64 = dist, dist
	if d2d < 10 {
		return FreeSpace(d2d, freqGHz), nil
	}

	if 10 <= d2d && d2d <= r.rmaNlosMax {

		loss1, err := r.los(d3d)

		// P3(indx)=C4+C5+C6*(mlog(d3d)-3)+C7;
		r.c4 = 161.04 - 7.1*mlog(r.streetW) + 7.5*mlog(rmaH)
		loss2 := r.c4 + r.c5 + r.c6*(mlog(d3d)-3) + r.c7
		// loss2 := 161.04 - 7.1*mlog(r.streetW) + 7.5*mlog(rmaH) - (24.37-3.7*math.Pow(rmaH/hBS, 2))*mlog(hBS) + (43.42-3.1*mlog(hBS))*(mlog(d3d)-3) + 20*mlog(freqGHz) - mpow(3.2*(mlog(11.75*hUT)), 2) - 4.97

		if err != nil {
			log.Println("NLOS ERR ", loss1, err)
		}

		nlosdB := max(loss1, loss2)

		if r.Extended { // if its LMLC mode
			nlosdB = max(loss1, loss2-r.nlossOffset)
		}

		//  }
		// }
		return nlosdB, nil
	} else {
		return 99999, fmt.Errorf("NLOS : Distance %f not supported in this model ", dist)
	}
}

func (r *RMa) p1(d3d, freqGHz float64) (plDb float64, valid bool) {

	plDb = 20*mlog(40*pi*d3d*freqGHz/3) + r.c1*mlog(d3d) - r.c2 + r.c3*d3d
	return plDb, true
}

func (r RMa) los(dist float64) (plDb float64, e error) {
	freqGHz := r.FGHz()

	var d3d, d2d float64 = dist, dist
	e = nil
	if d2d < 10 {
		flpl := FreeSpace(d2d, freqGHz)
		return flpl, e
	}
	if 10 <= d2d && d2d <= r.dBP {
		loss, ok := r.p1(d3d, freqGHz)
		if !ok {
			e = fmt.Errorf("LOS : PL1(%v,%v) Error ", dist, freqGHz)
		}
		return loss, e
	} else if d2d > r.dBP && d2d <= r.CutOffDistance {

		p1BP, ok := r.p1(r.dBP, freqGHz)
		// p1BP,ok:=r.p1(760, freqGHz)
		if !ok {
			e = fmt.Errorf("LOS : PL1(%v,%v) Error ", dist, freqGHz)
			return 99999, e
		}
		loss := p1BP + 40.0*mlog(d3d/r.dBP)
		//  loss:=p1BP+	40.0*mlog(d3d/769.0)
		return loss, nil
	} else {
		// return math.NaN(), fmt.Errorf("Unsupported distance %d for LOS ", dist)
		return 0, fmt.Errorf("LOS :Unsupported distance %.2f Max:[%f] ", dist, r.CutOffDistance)
	}
}

func (r *RMa) losNodes(src, dest vlib.Location3D) (plDb float64, valid bool) {

	d3d := src.DistanceFrom(dest)
	d2d := src.Distance2DFrom(dest)

	if 10 <= d2d && d2d <= r.dBP {
		loss, _ := r.los(d3d)
		return loss, true
	} else if d2d > r.dBP && d2d <= r.CutOffDistance {
		// loss, _ := r.p1(r.dBP, freqGHz)
		plDb, err := r.los(d3d)
		// plDb += 40 * mlog(d3d/r.dBP)
		valid = (err == nil)
		return plDb, valid
	} else {
		log.Printf("\nDistance not supported in this model")
		return 0, false

	}

}

func (r RMa) FnP1(d, fGHz float64) (pdDb float64, valid bool) {
	r.Check()
	return r.p1(d, fGHz)
}

func (r RMa) FnP2(d, fGHz float64) (pdDb float64, valid bool) {
	r.Check()
	p1BP, valid := r.p1(r.dBP, fGHz)
	p2 := p1BP + 40.0*mlog(d/r.dBP)

	return p2, valid

}
