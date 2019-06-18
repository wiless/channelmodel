package CM

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

type ShadowGrid struct {
	gridSize   float64
	maxX, maxY float64
	lsp        LSP
	grid       *mat.Dense
}

func NewShadowGrid(te TestEnv, gs float64, gsize float64) *ShadowGrid {
	var sg ShadowGrid
	sg.SetGridSize(gs, gsize)
	sg.lsp = loadLSP(te)
	return &sg
}

func (s *ShadowGrid) SetEnv(te TestEnv) {
	s.lsp = loadLSP(te)
}

// SetGridSize sets the size gs for both X & Y axis
// Set Model
func (s *ShadowGrid) SetGridSize(gs float64, maxsize float64) {
	s.gridSize = gs
	s.maxX, s.maxY = maxsize, maxsize
	N := int(math.Ceil(s.maxX / s.gridSize))
	s.grid = mat.NewDense(N, N, nil)
}

// AutoCorrFn implements AutoCorrelation between SF for links from same site
// M.2412 See  TABLE A1-22 (end) pg77
// 3.3.1 of WINNER II Channel Models [15]
// See Section 7.4.4 : 38.901
func AutoCorrFn(d float64, dcor float64) float64 {
	// var dcor float64 = 20 //
	d = math.Abs(d)
	return math.Exp(-d / dcor)
}

// loadLSP should return values from Table 4.5-6 in Pg 40 of 38.901
func loadLSP(te TestEnv) LSP {
	var result LSP
	if te == RMA {
		result = LSPLookupRMA()
	}
	return result
}

func (s *ShadowGrid) Create(lt LinkType) {
	// Currenly only ShadowGrid only for Shadow fading only
	dcor := s.lsp.SF[int(lt)] // LOS =0
	_ = dcor
}

type LinkType int

const (
	LOS LinkType = iota
	NLOS
	O2I
)

func (l LinkType) String() string {
	return [...]string{"LOS", "NLOS", "O2I"}[l]
}

type LSPname int

const (
	DS LSPname = iota
	ASD
	ASA
	SF
	K
	ZSA
	ZSD
)

func (l LSPname) String() string {
	return [...]string{"DS", "ASD", "ASA", "SF", "K", "ZSA", "ZSD"}[l]
}

type LSP struct {
	DS       []float64 // Fixed LOS, NLOS and O2I
	ASD, ASA []float64 // Fixed LOS, NLOS and O2I
	SF       []float64 // Fixed LOS, NLOS and O2I
	K        []float64 // Fixed LOS, NLOS and O2I
	ZSA, ZSD []float64 // Fixed LOS, NLOS and O2I
}

func (l LSP) Get(lsp LSPname) []float64 {
	switch lsp {
	case DS:
		return l.DS

	case SF:
		return l.SF
	default:
		return []float64{-1, -1, -1}
	}
}

// RMALSP returns the CorrDistance for RMA
//
// DS	ASD	ASA	SF	K	ZSA	ZSD
// 50	25	35	37	40	15	15  (LOS)
// 36	30	40	120	-1	50	50 	(NLOS)
// 36	30	40	120	-1	50	50  (O2I)

func LSPLookupRMA() LSP {
	var lspRMa LSP
	// a := [][]float64{
	// 	{50	25	35	37	40	15	15},
	// 	{36	30	40	120	-1	50	50},
	// 	{36	30	40	120	-1	50	50}
	// }
	myinfo := mat.NewDense(3, 7, []float64{
		50, 25, 35, 37, 40, 15, 15,
		36, 30, 40, 120, -1, 50, 50,
		36, 30, 40, 120, -1, 50, 50,
	})

	lspRMa.DS = mat.Col(nil, int(DS), myinfo)
	lspRMa.ASD = mat.Col(nil, int(ASD), myinfo)
	lspRMa.ASA = mat.Col(nil, int(ASA), myinfo)
	lspRMa.SF = mat.Col(nil, int(SF), myinfo)
	lspRMa.K = mat.Col(nil, int(K), myinfo)
	lspRMa.ZSA = mat.Col(nil, int(ZSA), myinfo)
	lspRMa.ZSD = mat.Col(nil, int(ZSD), myinfo)

	return lspRMa

}
