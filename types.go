// CM implements the functions & types for implementation of Channel models in 38-901-e00
package CM

import (
	"fmt"
	"math"

	clr "github.com/fatih/color"
	"github.com/wiless/vlib"
)

// // dist3D returns the 3D distance between two Nodes considering d_out & d_in as in Eq. 7.4.1
// func dist3D(src, dest deployment.Node) (d float64) {
// 	  d3d:=dest.Location.DistanceFrom(src)
//
// }
const C float64 = 3.0e8 // Speed of light
var DEFAULTERR_PL float64 = 99999

var mlog = math.Log10
var mpow = math.Pow
var max = math.Max
var mexp = math.Exp

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

type PLModel interface {
	Init(fGHz float64)
	IsSupported(fGHz float64) bool
	PLbetween(node1, node2 vlib.Location3D) (plDb float64, isNLOS bool, err error)
	IsLOS(dist float64) bool /// Given a distance returns LoS or NLOS statistical (it need not always return same value)
	PLnlos(dist float64) (plDb float64, e error)
	PLlos(dist float64) (plDb float64, e error)
	PL(dist float64) (plDb float64, isNLOS bool, err error)
}

// EvaluatePL returns the pathloss values for all the distances for the frequency, when fGHz or dist is not Supported, it returns error and corresponding values with DEFAULTERR_PL=99999 (not communicatable channel)
func EvaluatePL(pm *PLModel, fGHz float64, dists ...float64) {

}
