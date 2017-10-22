package main

import (
	"log"

	"github.com/wiless/vlib"

	plotutils "github.com/wiless/plotutils"

	"github.com/wiless/channelmodel/modelB"
)

func main() {
	testInH()
	// compareConfs()
}

func testInH() {

	var fc float64 = 4.0 // 4GHZ
	inh := modelB.NewInH(fc)

	_, y := inh.IsValid(fc, 500, true)
	if y != nil {
		log.Println(y)

	}
	plotutils.StartX()

	plotutils.Fig("Indoor HotSpot LOS vs NLOS")
	xx := vlib.ToVectorF("1:5:100")
	yy := vlib.NewVectorF(xx.Len())
	yy1 := vlib.NewVectorF(xx.Len())
	yy2 := vlib.NewVectorF(xx.Len())
	for i, v := range xx {
		yy[i], _ = inh.PLos(v)
		yy1[i], _ = inh.PNLos(v)
		yy2[i], _ = inh.PNLosOpt(v)
	}
	plotutils.PlotXY(xx, yy)
	plotutils.HoldOn()
	plotutils.PlotXY(xx, yy1)
	plotutils.PlotXY(xx, yy2)
	plotutils.SetLabel("Distance (m)", "Pathloss (dB)")
	plotutils.Legends("LOS", "NLOS", "NLOSopt")
	plotutils.SetYLim(0, 200)
	plotutils.SavePDF()
	// fmt.Println(inh.PNLos(10))

	plotutils.Wait()

}

func compareConfs() {
	var inhA, inhB, inhC modelB.InH
	inhA.Set(modelB.LoadInHConfA())
	inhB.Set(modelB.LoadInHConfB())
	inhC.Set(modelB.LoadInHConfC())

	plotutils.StartX()

	plotutils.Fig("Indoor HotSpot Conf A vs B vs C compared")
	xx := vlib.ToVectorF("1:5:100")
	yy1 := vlib.NewVectorF(xx.Len())
	yy2 := vlib.NewVectorF(xx.Len())
	yy3 := vlib.NewVectorF(xx.Len())
	var err error
	for i, v := range xx {
		yy1[i], err = inhA.PLos(v)
		if err != nil {
			log.Println("A", i, v, err)
		}
		yy2[i], err = inhB.PLos(v)
		if err != nil {
			log.Println("B", i, v, err)
		}
		yy3[i], err = inhC.PLos(v)
		if err != nil {
			log.Println("C", i, v, err)
		}
	}
	plotutils.PlotXY(xx, yy1)
	plotutils.HoldOn()
	plotutils.PlotXY(xx, yy2)
	plotutils.PlotXY(xx, yy3)
	plotutils.SetLabel("Distance (m)", "Pathloss (dB)")
	plotutils.Legends("Conf A 4GHz", "Conf B 30GHz", "Conf C 70GHz")
	plotutils.SetYLim(0, 200)
	plotutils.SavePDF()
	// fmt.Println(inh.PNLos(10))

	plotutils.Wait()

}
