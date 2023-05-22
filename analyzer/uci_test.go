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
