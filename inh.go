package CM

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/wiless/cellular/pathloss"
)

type InH struct {
	*pathloss.ModelSetting
	sigmaSF_Los  float64
	sigmaSF_NLos float64
	isOK         bool
}

// InHDefault returns a default Indoor HotSpot Model settings
func InHDefault() *pathloss.ModelSetting {
	ms := pathloss.NewModelSetting()
	ms.SetFGHz(.7)
	ms.Name = "InH"
	ms.AddParam("dMax_Los", 100.0).AddParam("dMax_NLos", 86)
	return ms
}

func (i *InH) Init(fcGHz float64) {

	i.Set(InHDefault().SetFGHz(fcGHz))

}

func (r *InH) Check() {
	if !r.isOK {
		log.Panic("InH Model not initialized, Call .Init() ")
	}
}

func (i *InH) Set(ms *pathloss.ModelSetting) {

	i.ModelSetting = ms
	fGHz := i.FGHz()
	if fGHz == 0 {
		i.isOK = false
		log.Print("InH ModelSettings Frequency not set !!")
	}

	dMax_Los, dMax_NLos := i.Value("dMax_Los"), i.Value("dMax_NLos")
	if dMax_Los == 0 || dMax_NLos == 0 {
		log.Panic("InH : Unable to find Mandatory paramters : dMax_Los & dMax_NLos in the input model ", ms.Name)
	}
	i.sigmaSF_Los = 3
	i.sigmaSF_NLos = 8.03
	i.isOK = true
}

func (i *InH) PL(d float64) (plDB float64, IsLos bool, err error) {
	i.Check()

	dMax_Los, dMax_NLos := i.Value("dMax_Los"), i.Value("dMax_NLos")

	if d <= 0 {
		return 0, false, fmt.Errorf("InH Distance should be non-ZERO ")
	}
	fGhz := i.FGHz()
	IsLos = i.IsLOS(d)

	if IsLos {

		if d <= dMax_Los {
			plDB = 32.4 + 17.3*math.Log10(d) + 20*math.Log10(fGhz)
			return plDB, IsLos, nil
		} else {

			return 0, false, fmt.Errorf("InH LOS Distance not supported beyond %f ", dMax_Los)
		}
	}

	if !IsLos {

		if d <= dMax_NLos {
			p1 := 32.4 + 17.3*math.Log10(d) + 20*math.Log10(fGhz)
			p2 := 38.3*math.Log10(d) + 17.30 + 24.9*math.Log10(fGhz)
			plDB = math.Max(p1, p2)
			return plDB, IsLos, nil
		} else {
			return 0, IsLos, fmt.Errorf("InH NLOS Distance not supported beyond %f ", dMax_NLos)
		}

	}

	return 0, true, nil
}

func (i *InH) IsLOS(d float64) bool {

	i.Check()
	if d <= 5 {
		return true
	}
	var P_los float64
	if d > 5 && d <= 49 {
		P_los = math.Exp(-(d - 5.0) / 70.8)
	}
	if d > 49 {
		P_los = 0.54 * math.Exp(-(d-49.0)/211.7)
	}

	if rand.Float64() <= P_los {
		return true
	}

	return false

}
