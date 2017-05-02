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
	var pl CM.InH
	var rma CM.RMa

	pl.Set(CM.InHDefault().SetFGHz(1.8))

	rma.Set(CM.RMADefault().SetFGHz(.7))

	var d1 float64
	d1 = 1
	// var dist, ploss vlib.VectorF
	dist := vlib.NewVectorF(100)
	ploss := vlib.NewVectorF(100)
	ploss2 := vlib.NewVectorF(100)
	counter := 0
	for d1 < 100.0 {
		pDb, L, err := pl.PL(d1)
		pDb2, _, _ := rma.PL(d1)
		_ = L
		if err == nil {
			dist[counter] = d1
			ploss[counter] = pDb
			ploss2[counter] = pDb2
			counter++
		}
		d1++
	}

	dist.Resize(counter)
	ploss.Resize(counter)

	var m vlib.MatrixF
	m.AppendColumn(dist).AppendColumn(ploss).AppendColumn(ploss2)

	pf.StartX()

	pf.Fig("Indoor HotSpot Model")
	pf.Plot(&m)

	pf.SetXlabel("distance (m)")
	pf.SetYlabel("PathLoss (dBm)")

	pf.HoldOn()
	log.Print("Next image starts... ")
	time.Sleep(2 * time.Second)
	pf.Plot(&m, 0, 2)

	pf.SetXlabel("distance (m)")
	pf.SetYlabel("PathLoss (dBm)")

	pf.Wait()
}

func init() {
	now := time.Now()

	fmt.Printf("=======   Starting code at %v  ======== \n", now)

}
