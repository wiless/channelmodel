package modelB

import (
	"math"

	"github.com/wiless/cellular/pathloss"
)

var mlog = math.Log10
var mpow = math.Pow
var max = math.Max
var mexp = math.Exp

func LoadInHConfA() *pathloss.ModelSetting {
	//. Loads
	ms := pathloss.NewModelSetting()
	ms.SetFGHz(4)
	ms.Name = "InH_ConfA"
	// ms.AddParam("dMax_Los", 100.0).AddParam("dMax_NLos", 86)
	return ms
}

func LoadInHConfB() *pathloss.ModelSetting {
	//. Loads
	ms := pathloss.NewModelSetting()
	ms.SetFGHz(30)
	ms.Name = "InH_ConfB"
	// ms.AddParam("dMax_Los", 100.0).AddParam("dMax_NLos", 86)
	return ms
}

func LoadInHConfC() *pathloss.ModelSetting {
	//. Loads

	ms := pathloss.NewModelSetting()
	ms.SetFGHz(70)
	ms.Name = "InH_ConfC"
	// ms.AddParam("dMax_Los", 100.0).AddParam("dMax_NLos", 86)
	return ms
}
