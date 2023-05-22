package ai

import (
	"testing"
)

func TestFromFileRankToSq(t *testing.T) {

	if fr_to_SQ120(FILE_A, RANK_1) != SQUARES_A1 {
		t.Errorf("fr_to_SQ120( FILE_A, RANK_1) != SQUARES_A1")
	}
	if fr_to_SQ120(FILE_A, RANK_8) != SQUARES_A8 {
		t.Errorf("fr_to_SQ120( FILE_A, RANK_8)")
	}

	if fr_to_SQ120(FILE_G, RANK_1) != SQUARES_G1 {
		t.Errorf("fr_to_SQ120( FILE_G, RANK_1)")
	}

	if fr_to_SQ120(FILE_H, RANK_1) != SQUARES_H1 {
		t.Errorf("fr_to_SQ120( FILE_H, RANK_1)")
	}
}

func TestMapSq120And64(t *testing.T) {

	// A1 = 0 (64) and SQUARES_A1 (120)
	if sq_64(fr_to_SQ120(FILE_A, RANK_1)) != 0 {
		t.Errorf("sq_64( fr_to_SQ120( FILE_A, RANK_1) ) != 0")
	}

	if sq_120(0) != fr_to_SQ120(FILE_A, RANK_1) {
		t.Errorf("sq_120( 0  ) != fr_to_SQ120( FILE_A, RANK_1)")
	}

	// B2 = 9 (64) and SQUARES_B1+10 (120)
	if sq_64(fr_to_SQ120(FILE_B, RANK_2)) != 9 {
		t.Errorf("sq_64( fr_to_SQ120( FILE_B, RANK_2) ) != 9")
	}

	if sq_120(9) != fr_to_SQ120(FILE_B, RANK_2) {
		t.Errorf("sq_120( 9  ) != fr_to_SQ120( FILE_B, RANK_2)")
	}

}

func TestMapping1FileRanks(t *testing.T) {
	if sq_to_FILE(SQUARES_A8) != FILE_A {
		t.Errorf("sq_to_FILE( SQUARES_A8)")
	}

	if sq_to_RANK(SQUARES_A8) != RANK_8 {
		t.Errorf("sq_to_RANK( SQUARES_A8)")
	}

	if filesBrd[0] != 100 {
		t.Errorf("filesBrd[ 0 ] != 100")
	}

	if ranksBrd[0] != 100 {
		t.Errorf("ranksBrd[ 0 ] != 100")
	}

	if filesBrd[SQUARES_A1] != 0 {
		t.Errorf("filesBrd[ SQUARES_A1 ] != 0")
	}

	if ranksBrd[SQUARES_A1] != 0 {
		t.Errorf("ranksBrd[ SQUARES_A1 ] != 0")
	}

	if ranksBrd[SQUARES_A1] != 0 {
		t.Errorf("ranksBrd[ SQUARES_A1 ] != 0")
	}

	if filesBrd[SQUARES_E8] != 4 {
		t.Errorf("filesBrd[ SQUARES_E8 ] != 4")
	}
	if ranksBrd[SQUARES_E8] != 7 {
		t.Errorf("ranksBrd[ SQUARES_E8 ] != 7")
	}

}
