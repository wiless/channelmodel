package modelB

import (
	"errors"
	"fmt"

	"github.com/wiless/cellular/pathloss"
)

type UMi struct {
	*pathloss.ModelSetting
	SFlos, SFNlos float64 // ShadowFading variance for LOS & NLOS in dB
}

func (i *UMi) Set(ms *pathloss.ModelSetting) {

	i.ModelSetting = ms
	i.SFlos = 3     //dB
	i.SFNlos = 8.03 //dB

}

func NewUMi(fcGHz float64) *UMi {
	result := new(UMi)
	result.ModelSetting = new(pathloss.ModelSetting)
	result.SetFGHz(fcGHz)
	return result
}

func (i *UMi) ok() (bool, error) {
	if i.ModelSetting == nil {
		return false, errors.New("No MODEL set!! Call Init()")
	}

	if i.FGHz() == 0 {
		return false, errors.New("Freq set to ZERO !!")
	}

	return true, nil
}

func (i *UMi) IsValid(fcGHz float64, dist_m float64, los bool) (bool, error) {
	// 0.5GHz≤ fc ≤100GHz

	if !(fcGHz >= 0.5 && fcGHz <= 100) {
		return false, fmt.Errorf("UMi : Freq %f outside range 0.5 < fc < 100GHz", fcGHz)
	}

	if los {
		// 1 m d3D  100 m
		if !(dist_m >= 1 && dist_m <= 150) {
			return false, fmt.Errorf("UMi : LOS dist  %f outside  range 1 < d3D < 150m ", dist_m)
		}
	} else {
		// NLOS
		if !(dist_m >= 1 && dist_m <= 150) { /// UPDATE TO 150 from 86 as WP5D#28
			return false, fmt.Errorf("UMi : LOS dist %f outside  range 1 < d3D < 86m ", dist_m)
		}
	}

	return true, nil
}

func (i *UMi) PLos(d3D float64) (pldB float64, err error) {

	if ok, err := i.ok(); ok != true {
		return 0, err
	}

	fc := i.FGHz()
	// PLUMi-LOS =32.4+17.3log10(d3D)+20log10(fc),σSF =3dB, 1 m d3D  100 m
	valid, err := i.IsValid(fc, d3D, true)
	if !valid {
		return 0, err
	}
	pldB = 32.4 + 17.3*mlog(d3D) + 20*mlog(fc)
	return pldB, nil
}

func (i *UMi) PNLos(d3D float64) (pldB float64, e error) {
	if ok, err := i.ok(); ok != true {
		return 0, err
	}
	// PL =max(PL ,PLʹnlos )
	// PLnlosʹ =38.3log10(d )+17.30+24.9log10(f ),σ =8.03dB,1 m d3D  86 m
	fc := i.FGHz()
	valid, err := i.IsValid(fc, d3D, false)
	if !valid {
		return 0, err
	}

	pldB1, err := i.PLos(d3D)
	if err != nil {
		return 0, err
	}

	pldB2 := 38.3*mlog(d3D) + 17.30 + 24.9*mlog(fc)

	return max(pldB1, pldB2), nil
}

func (i *UMi) PNLosOpt(d3D float64) (pldB float64, e error) {
	if ok, err := i.ok(); ok != true {
		return 0, err
	}
	// PL =max(PL ,PLʹnlos )
	//  Optional:PLʹ =32.4+20log (f )+31.9log (d ), σ =8.29dB
	i.SFNlos = 8.29 // if option called

	fc := i.FGHz()
	valid, err := i.IsValid(fc, d3D, false)
	if !valid {
		return 0, err
	}

	pldB1, err := i.PLos(d3D)
	if err != nil {
		return 0, err
	}

	// pldB2 := 38.3*mlog(d3D) + 17.30 + 24.9*mlog(fc)
	pldB2 := 32.4 + 20*mlog(fc) + 31.9*mlog(d3D)
	return max(pldB1, pldB2), nil
}
