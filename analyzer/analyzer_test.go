package analyzer

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestAnalyzer1(t *testing.T) {

	a := NewAnalyzer()
	assert.NotNil(t, a)

	a.Fen = "2q1rr1k/3bbnnp/p2p1pp1/2pPp3/PpP1P1P1/1P2BNNP/2BQ1PRK/7R b - -"
	a.Fen = "2r3k1/p4p2/3Rp2p/1p2P1pK/8/1P4P1/P3Q2P/1q6 b - - 0 1"
	a.Fen = "8/k2r4/p7/2b1Bp2/P3p3/qp4R1/4QP2/1K6 b - - 0 1"
	a.Fen = "6k1/pp4p1/2p5/2bp4/8/P5Pb/1P3rrP/2BRRN1K b - - 0 1"
	a.Fen = "3r4/pR2N3/2pkb3/5p2/8/2B5/qP3PPP/4R1K1 w - - 1 0"
	a.Fen = "rn3rk1/pbppq1pp/1p2pb2/4N2Q/3PN3/3B4/PPP2PPP/R3K2R w KQ - 7 11"
	a.Fen = "2r3k1/p4p2/3Rp2p/1p2P1pK/8/1P4P1/P3Q2P/1q6 b - - 0 1"
	a.Fen = "6k1/3b3r/1p1p4/p1n2p2/1PPNpP1q/P3Q1p1/1R1RB1P1/5K2 b - - 0-1"
	a.Fen = "2r3k1/p4p2/3Rp2p/1p2P1pK/8/1P4P1/P3Q2P/1q6 b - - 0 1"
	a.UserMove = ""

	r, err := a.AnalyzeFen()
	assert.Nil(t, err)

	assert.NotNil(t, r)

}
