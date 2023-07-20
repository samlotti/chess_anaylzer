package ai

import (
	"fmt"
	"testing"
)

func TestE2E4(t *testing.T) {

	move := buildMOVE(fr_to_SQ120(FILE_E, RANK_2),
		fr_to_SQ120(FILE_E, RANK_4),
		0,
		0,
		0,
	)

	fmt.Println(moveToDebugString(move))

	expect(t, MoveToString(move), "e2e4")

}

func TestE7E8Promote(t *testing.T) {

	move := buildMOVE(fr_to_SQ120(FILE_E, RANK_7),
		fr_to_SQ120(FILE_E, RANK_8),
		0,
		Piece_BQUEEN,
		0,
	)

	fmt.Println(moveToDebugString(move))

	expect(t, MoveToString(move), "e7e8q")

}

func TestMovePull(t *testing.T) {

	move := buildMOVE(fr_to_SQ120(FILE_E, RANK_7),
		fr_to_SQ120(FILE_E, RANK_8),
		Piece_WBISHOP,
		Piece_BQUEEN,
		MFLAG_PawnStart,
	)

	expectInt(t, "1", getCapturedPiece(move), Piece_WBISHOP)
	expectInt(t, "2", getPromoted(move), Piece_BQUEEN)
	expectInt(t, "3", move&MFLAG_PawnStart, MFLAG_PawnStart)

}

func TestMakeMove(t *testing.T) {

	var b = genAndPrint("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	fmt.Println("About to make a move")
	generateMoves(b)

	if !b.checkBoard() {
		panic(fmt.Errorf("Check board failed!!"))
	}

	b.printMoveList()
	b.makeMove(b.moveList[0]) // sb a2 - a3
	b.checkBoard()

	b.PrintBoard("Made a move")
	b.checkBoard()

	b.takeMove()
	b.PrintBoard("Take back a move")
	b.checkBoard()

	//
	//
	//println("Reverse move")
	//b.takeMove()
	//b.printBoard()
	//b.checkBoard()

}

func genAndPrint(fen string) *Board {
	fmt.Println("Testing FEN: ", fen)
	b := NewBoard()
	ParseFen(b, fen)
	generateMoves(b)
	b.PrintBoard(fen)
	// b.printMoveList()
	// b.checkBoard()
	return b
}
