package ai

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitialPosition(t *testing.T) {

	var fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	brd := NewBoard()
	ParseFen(brd, fen)
	brd.printBoard( fen )
	// brd.PrintPieceList()
	brd.isSqAttacked(fr_to_SQ120(1, 1), brd.side)
	brd.printSquareAttacked(brd.side)

	moves := make([]Move,0)
	var revfen = BoardToFen(brd, len(moves))
	assert.Equal(t, fen, revfen)

}

func TestInitialPositionE3(t *testing.T) {

	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"

	brd := NewBoard()
	ParseFen(brd, fen)
	brd.printBoard("TestInitialPositionE3")

	// brd.PrintPieceList()
	att := brd.isSqAttacked(fr_to_SQ120(1, 1), brd.side)
	fmt.Println("Sq Attacked: ", att)
	brd.printSquareAttacked(brd.side)

	moves := make([]Move,0)
	var revfen = BoardToFen(brd, len(moves))
	assert.Equal(t, fen, revfen)

}

/*
8    .    .    .    .    .    .    .    k
7    .    .    .    .    .    .    .    .
6    .    .    .    .    .    .    .    .
5    p    .    .    .    .    .    .    .
4    .    P    .    .    .    .    .    .
3    .    .    .    .    .    .    .    .
2    .    .    .    .    .    .    .    .
1    .    .    .    .    .    .    .    K
*/
func TestSQAttacked1(t *testing.T) {

	brd := NewBoard()
	ParseFen(brd, "7k/8/8/p7/1P6/8/8/7K w - - 0 1")
	brd.printBoard("TestSQAttacked1")

	// brd.PrintPieceList()
	att := brd.isSqAttacked(fr_to_SQ120(1, 1), brd.side)
	fmt.Println("Sq Attacked: ", att)
	brd.printSquareAttacked(Color_WHITE)

	var sq120 = fr_to_SQ120(FILE_A, RANK_6)
	expectAttacked(brd, sq120, Color_WHITE, t, false)

	sq120 = fr_to_SQ120(FILE_A, RANK_5)
	expectAttacked(brd, sq120, Color_WHITE, t, true)

	sq120 = fr_to_SQ120(FILE_H, RANK_2)
	expectAttacked(brd, sq120, Color_WHITE, t, true)

	sq120 = fr_to_SQ120(FILE_H, RANK_2)
	expectAttacked(brd, sq120, Color_BLACK, t, false)

}

func expectAttacked(brd *Board, sq120 SQ120, color Color, t *testing.T, expected bool) {
	if brd.isSqAttacked(sq120, color) != expected {
		t.Error("Expected ", expected, "  but got opposite")
	}
}

func TestSQAttacked2(t *testing.T) {

	brd := NewBoard()
	ParseFen(brd, "q3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1")
	brd.printBoard("TestSQAttacked2")

	// brd.PrintPieceList()
	att := brd.isSqAttacked(fr_to_SQ120(1, 1), brd.side)
	fmt.Println("Sq Attacked: ", att)
	brd.printSquareAttacked(Color_WHITE)
	brd.printSquareAttacked(Color_BLACK)

	var sq120 = fr_to_SQ120(FILE_A, RANK_6)
	expectAttacked(brd, sq120, Color_WHITE, t, true)

	sq120 = fr_to_SQ120(FILE_A, RANK_5)
	expectAttacked(brd, sq120, Color_BLACK, t, true)

	sq120 = fr_to_SQ120(FILE_H, RANK_2)
	expectAttacked(brd, sq120, Color_WHITE, t, true)

	sq120 = fr_to_SQ120(FILE_H, RANK_2)
	expectAttacked(brd, sq120, Color_BLACK, t, true)

	sq120 = fr_to_SQ120(FILE_B, RANK_7)
	expectAttacked(brd, sq120, Color_BLACK, t, true)

	sq120 = fr_to_SQ120(FILE_B, RANK_7)
	expectAttacked(brd, sq120, Color_WHITE, t, false)

	brd.PrintPieceList()

}

func TestSQAttackedKnight(t *testing.T) {

	brd := NewBoard()
	ParseFen(brd, "q3k2r/8/8/3N4/8/8/8/R3K2R b KQkq - 0 1")
	brd.printBoard("TestSQAttackedKnight")

	// brd.PrintPieceList()
	att := brd.isSqAttacked(fr_to_SQ120(1, 1), brd.side)
	fmt.Println("Sq Attacked: ", att)
	brd.printSquareAttacked(Color_WHITE)
	brd.printSquareAttacked(Color_BLACK)

	var sq120 = fr_to_SQ120(FILE_B, RANK_6)
	expectAttacked(brd, sq120, Color_WHITE, t, true)
	expectAttacked(brd, sq120, Color_BLACK, t, false)

}
