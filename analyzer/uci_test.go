package analyzer

import (
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"

func TestSegments(t *testing.T) {

	u := NewUci("zahak")
	assert.NotNil(t, u)

	err := u.Start()
	assert.Nil(t, err)
	err = u.WaitOk(1 * time.Second)
	assert.Nil(t, err)

	err = // u.SetPositionFen("r2qnrnk/p2b2b1/1p1p2pp/2pPpp2/1PP1P3/PRNBB3/3QNPPP/5RK1 w - -")
		u.SetPositionFen("2q1rr1k/3bbnnp/p2p1pp1/2pPp3/PpP1P1P1/1P2BNNP/2BQ1PRK/7R b - -")
	assert.Nil(t, err)

	err = u.SetOption("MultiPV", "15")
	assert.Nil(t, err)

	// u.send("setoption name MultiPV value 5")

	o := NewGoOptions()
	o.Depth = 15

	err = u.SendGo(o)
	assert.Nil(t, err)

	err = u.WaitMoveUpTo(5 * time.Second)
	assert.Nil(t, err)

	assert.True(t, u.IsReadyForMove())

	u.Terminate()

}

func TestParseInfo1(t *testing.T) {

	u, err := UciInfoParse(
		"info depth 11 seldepth 25 hashfull 17 tbhits 0 nodes 391537 nps 769164 score cp -305 time 509 multipv 3 pv a5b4",
	)
	assert.Nil(t, err)
	assert.Equal(t, 11, u.Depth)
	assert.Equal(t, 769164, u.Nps)
	assert.Equal(t, -305, u.ScoreCP)
	assert.Equal(t, 3, u.MPv)
	assert.Equal(t, 1, len(u.Moves))
	assert.Equal(t, 0, u.MateIn)

}

func TestParseInfo2(t *testing.T) {

	u, err := UciInfoParse(
		"info depth 15 seldepth 27 hashfull 47 tbhits 0 nodes 1312678 nps 1035602 score mate +3 time 1267 multipv 1 pv b1g6 h5g4 g6f5 g4h5",
	)
	assert.Nil(t, err)
	assert.Equal(t, 15, u.Depth)
	assert.Equal(t, 1035602, u.Nps)
	assert.Equal(t, 15003, u.ScoreCP)
	assert.Equal(t, 1, u.MPv)
	assert.Equal(t, 4, len(u.Moves))
	assert.Equal(t, 3, u.MateIn)

}
