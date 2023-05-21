package analyzer

import (
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"

func TestSegments(t *testing.T) {

	u := newUci("zahak")
	assert.NotNil(t, u)

	err := u.Start()
	assert.Nil(t, err)
	err = u.WaitOk(1 * time.Second)
	assert.Nil(t, err)

	err = u.SetPositionFen("r2qnrnk/p2b2b1/1p1p2pp/2pPpp2/1PP1P3/PRNBB3/3QNPPP/5RK1 w - -")
	assert.Nil(t, err)

	err = u.SendGo(NewGoOptions())
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	err = u.SendStop()
	assert.Nil(t, err)

	err = u.WaitOk(1 * time.Second)
	assert.Nil(t, err)

	u.Terminate()

}
