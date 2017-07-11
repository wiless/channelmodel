package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wiless/channelmodel"
	"github.com/wiless/plotutils"
	"github.com/wiless/vlib"
)

func ByeBye() {
	fmt.Print("\n ============= Bye bye  ========= \n")
}

func main() {
	defer ByeBye()
	rand.Seed(time.Now().Unix())
	var rma CM.RMa

	rma.Set(CM.RMADefault().SetFGHz(.7))
	rma.SetExtended(false)
	// var dist, ploss vlib.VectorF
	// dist := vlib.NewVectorF(100)
	counter := 0

	var dist vlib.VectorF
	var p1, p2, plos, pnlos vlib.VectorF

	_, _ = p1, p2
	dmax := rma.DMax()
	for d := 10.0; d < dmax; d += 100 {

		v1, _ := rma.FnP1(d, rma.FGHz())
		v2, _ := rma.FnP2(d, rma.FGHz())
		p1db, e1 := rma.PLlos(d)
		p2db, e2 := rma.PLnlos(d)
		dist.AppendAtEnd(d)

		p1.AppendAtEnd(v1)
		p2.AppendAtEnd(v2)
		plos.AppendAtEnd(p1db)
		pnlos.AppendAtEnd(p2db)

		if e1 == nil && e2 == nil {
			// dist.AppendAtEnd(d)
			// plos.AppendAtEnd(p1db)
			// pnlos.AppendAtEnd(p2db)
			// counter++
		} else {
			log.Printf("Err [%v]: %v", d, e1, e2)
		}
		counter++

	}

	var m vlib.MatrixF
	m.AppendColumn(dist).AppendColumn(plos).AppendColumn(pnlos)

	pf.StartX()

	pf.Fig("Path Loss Rma - LOS")
	pf.SemilogX(dist, p1)
	pf.SetYLim(0, 140)
	pf.HoldOn()
	pf.SemilogX(dist, p2)
	pf.SemilogX(dist, plos)
	// pf.Plot(&m, 0, 2)
	pf.SemilogX(dist, pnlos)
	pf.SetYLim(0, 140)
	pf.SetXlabel("distance (m)")
	pf.SetYlabel("PathLoss (dBm)")
	pf.Legends("P1", "P2", "LOS", "NLOS")

	pf.Fig("LOS Components P1 & P2")
	pf.SemilogX(dist, p1)
	pf.HoldOn()
	pf.SemilogX(dist, p2)
	pf.SetYLim(0, 140)
	pf.SetXlabel("distance (m)")
	pf.SetYlabel("PathLoss (dBm)")
	pf.Legends("P1", "P2")

	pf.Wait()
}

func init() {
	now := time.Now()

	fmt.Printf("=======   Starting code at %v  ======== \n", now)

}
