// CM implements the functions & types for implementation of Channel models in 38-901-e00
package CM

import (
	"fmt"
	"math"
	"math/rand"

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

var msqrt = math.Sqrt
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
	PLbetweenIndoor(node1, node2 vlib.Location3D, dIn float64) (plDb float64, isLOS bool, err error)
	IsLOS(dist float64) bool /// Given a distance returns LoS or NLOS statistical (it need not always return same value)
	PLnlos(dist float64) (plDb float64, e error)
	PLlos(dist float64) (plDb float64, e error)
	PL(dist float64) (plDb float64, isNLOS bool, err error)
	O2ILossDb(fGHz float64, d2Din float64) (o2ilossdB float64)
}

// EvaluatePL returns the pathloss values for all the distances for the frequency, when fGHz or dist is not Supported, it returns error and corresponding values with DEFAULTERR_PL=99999 (not communicatable channel)
func EvaluatePL(pm *PLModel, fGHz float64, dists ...float64) {

}

// https://wikimedia.org/api/rest_v1/media/math/render/svg/2748541aa04938707a3d25923da2290f4d32ab59
func FreeSpace(dmtrs float64, fGHz float64) (plDb float64) {
	// return 20*mlog(dmtrs) + 20*mlog(fGHz*1e9) - 27.55
	// return 20*mlog(dmtrs) + 20*mlog(fGHz)+20*mlog(1e9)-147.55
	return 20*mlog(dmtrs) + 20*mlog(fGHz) + 32.45
}

// O2ICarLossDb returns the Car penetration loss in dB
// Ref M.2412 Section 3.3 μ = 9, and σP = 5
func O2ICarLossDb() float64 {
	// μ = 9, and σP = 5
	mean := 9.0
	sigmaP := 5.0
	return rand.NormFloat64()*sigmaP + mean
}

// All loss below are in dB
func Lglass(fGHz float64) float64 {
	return 2 + 0.2*fGHz
}
func LIRRglass(fGHz float64) float64 {
	return 23 + 0.3*fGHz
}
func Lconcrete(fGHz float64) float64 {
	return 5 + 4*fGHz
}
func Lwood(fGHz float64) float64 {
	return 4.85 + 0.12*fGHz
}

// RandLogNorm returns a log-normal distributed random variable
// with mean = u and std-devation = sigma
func RandLogNorm(meanDb float64, sigmaDb float64) float64 {

	// m := meanDb
	// v := sigmaDb * sigmaDb // mean and variance of Y
	// phi := msqrt(v + m*m)
	// mu := mlog(m * m / phi)                   // mean of log(Y)
	// sigma := msqrt(mlog(phi * phi / (m * m))) // std dev of log(Y)
	// result := mexp(mu + sigma*rand.Float64())
	result := mexp(meanDb + sigmaDb*rand.Float64())
	return result
}
