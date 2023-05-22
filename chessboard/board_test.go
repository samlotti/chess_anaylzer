package ai

import (
	"testing"
)

func TestOffBoard(t *testing.T) {
	brd := NewBoard()
	ParseFen(brd, StartFen)

	expectBool(t, "1", brd.isSqOffBoard(fr_to_SQ120(FILE_A, RANK_6)), false)
	expectBool(t, "2", brd.isSqOffBoard(fr_to_SQ120(FILE_A, RANK_6)-1), true)
}

