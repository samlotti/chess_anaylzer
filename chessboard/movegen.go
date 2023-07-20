package ai

import (
	"fmt"
	s "strings"
)

/**
 *
 * Info to store
 *
 * FromSq
 * ToSQ
 * isEnPass Capture
 * Captured Piece
 * Promoted Piece
 * Pawn Start
 * Castle Move
 *
 * Will store in integers of 31 bits
 *
 *  A-F
 * 0000 0000 0000 0000   0000 0000 0000 0000
 * 0000 0000 0000 0000   0000 0000 0001 0000 = 16
 * 0000 0000 0000 0000   0000 0001 0000 0000 = 256
 *
 * 0000 0000 0000 0000   0000 0001 0000 0000
 *
 * Square values from 21 to 98 -> fit in 7 bits
 *
 * 0000 0000 0000 0000   0000 0000 0111 1111 - From Square  ( d & 0x7f )=From square
 * 0000 0000 0000 0000   0011 1111 1000 0000 - To Sq        ( >> 7  & 0x7f = To Sq
 * 0000 0000 0000 0011   1100 0000 0000 0000 - Captures     ( >> 7+7  & 0xf = Captured piece
 * 0000 0000 0000 0100   0000 0000 0000 0000 - Is Enpassant ( & 0x40000 )
 * 0000 0000 0000 1000   0000 0000 0000 0000 - Is Pawn strt ( & 0x80000 )
 * 0000 0000 1111 0000   0000 0000 0000 0000 - Promoted Piece ( >> 20, 0xF ) = Promoted piece
 * 0000 0001 0000 0000   0000 0000 0000 0000 - Castling      ( &0x1000000)
 * 0000 0000 0000 0000   0000 0000 0000 0000
 *
 *
 * For performance these are static at the MoveGenKt
 */
const MFLAG_EnPassant = 0x40000
const MFLAG_PawnStart = 0x80000
const MFLAG_Castling = 0x1000000

// This is a combination of enpassant and capture piece
// Not really used... appears to be offset and overlays other bits.
// const MFLAG_Captured = 0x7c000

const MFLAG_Promoted = 0xF00000

const NOMOVE Move = 0

func buildMOVE(from120 SQ120, to120 SQ120, captured Piece, promoted Piece, flag int) Move {
	var move = from120 | (to120 << 7) | (captured << 14) | (promoted << 20) | (flag)
	return move
}

func getFromSq120(move Move) SQ120 {
	return move & 0x7f
}

func getToSq120(move Move) SQ120 {
	return (move >> 7) & 0x7f
}

func getCapturedPiece(move Move) int {
	return (move >> 14) & (0xf)
}

func getPromoted(move Move) int {
	return (move >> 20) & (0xf)
}

func MoveToString(move Move) string {

	if move == 0 {
		return "?"
	}

	var fSq120 = getFromSq120(move)
	var tSq120 = getToSq120(move)
	var promoted = getPromoted(move)
	var promotedStr = ""
	if promoted != Piece_EMPTY {
		promotedStr = s.ToLower(string(PceCharLetter[promoted]))
		// promotedStr = ("" + pceCharLetter[promoted])
	}
	return fmt.Sprintf("%s%s%s", sqToString(fSq120), sqToString(tSq120), promotedStr)
}

func moveToDebugString(move Move) string {

	var fSq120 = getFromSq120(move)
	var tSq120 = getToSq120(move)
	var promoted = getPromoted(move)
	var promotedStr = ""
	if Piece_EMPTY != promoted {
		promotedStr = s.ToLower(string(PceCharLetter[promoted]))
		// promotedStr = ("" + pceCharLetter[promoted])
	}

	var capPiece = getCapturedPiece(move)
	var capPStr = ""
	if Piece_EMPTY != capPiece {
		capPStr = "  cap:" + s.ToLower(string(PceCharLetter[capPiece]))
	}

	var flagStr = "   F:"
	if move&MFLAG_Castling != 0 {
		flagStr += "Cstl "
	}
	if move&(MFLAG_EnPassant) != 0 {
		flagStr += "Ep "
	}
	if move&(MFLAG_PawnStart) != 0 {
		flagStr += "PwnSt "
	}
	if move&MFLAG_Promoted != 0 {
		flagStr += "Prmt "
	}
	//if move&MFLAG_Captured != 0 {
	//	flagStr += "Cpt "
	//}

	return fmt.Sprintf("%s%s%s%s%s", sqToString(fSq120), sqToString(tSq120), promotedStr, capPStr, flagStr)
}

//
//func MoveExists(move Move, board *Board) bool {
//	board.GenerateMoves()
//
//	var moveFound = NOMOVE
//	for idx := board.moveListStart[board.ply]; idx < board.moveListStart[board.ply + 1] {
//		moveFound = move
//		// board.MoveList[idx]
//		if !board.makeMove(moveFound) {
//			continue
//		}
//		board.takeMove()
//		if (move == moveFound) {
//			return true
//		}
//	}
//	return false
//}

/*
*

	 // Where the move starts at a given Depth. (?? index lookup ??)
		 moveListStart[] -> index for the first move at a given ply

	 // The moves at each ply (searching)
		moveList[index]

		So movelist contains a block of all move ply0, then block all moves ply 1 ....
		moveliststart just tells the index at the ply for the beginning of the moves.

		dynamic index structure.  as moves added to a ply, increments the start of the next ply to point to current pos + 1
	  for (i in moveListStart(ply)..moveListStart(ply+1)-1) {
		 processmove.
	  }
*/
func generateMoves(brd *Board) {

	if DEBUG {
		brd.PrintBoard("Generate moves")
	}

	brd.moveListStart[brd.ply+1] = brd.moveListStart[brd.ply]

	var cur_pceType = 0
	var cur_pceNum = 0
	var sq SQ120 = 0
	var pceIndex = 0
	var pce = 0
	var t_sq SQ120 = 0
	var dir = 0

	if brd.side == Color_WHITE {
		cur_pceType = Piece_WPAWN
		for cur_pceNum = 0; cur_pceNum < brd.pceNum[cur_pceType]; cur_pceNum++ {
			sq = brd.pList[pieceIndex(cur_pceType, cur_pceNum)]

			if brd.pieces[sq+10] == Piece_EMPTY {
				brd.addWhitePawnQuietMove(sq, sq+10)

				if ranksBrd[sq] == RANK_2 && brd.pieces[sq+20] == Piece_EMPTY {
					// Add quiet move here
					brd.addQuietMove(buildMOVE(sq, sq+20, Piece_EMPTY, Piece_EMPTY, MFLAG_PawnStart))
				}
			}

			if brd.isSqOffBoard(sq+9) == false && PieceCol[brd.pieces[sq+9]] == Color_BLACK {
				brd.addWhitePawnCaptureMove(sq, (sq + 9), brd.pieces[sq+9])
			}

			if brd.isSqOffBoard(sq+11) == false && PieceCol[brd.pieces[sq+11]] == Color_BLACK {
				brd.addWhitePawnCaptureMove(sq, sq+11, brd.pieces[sq+11])
			}

			if brd.enPassantSq120 != SQUARES_NO_SQ {
				if sq+9 == brd.enPassantSq120 {
					brd.addEnPassantMove(buildMOVE(sq, sq+9, Piece_EMPTY, Piece_EMPTY, MFLAG_EnPassant))
				}
				if sq+11 == brd.enPassantSq120 {
					brd.addEnPassantMove(buildMOVE(sq, sq+11, Piece_EMPTY, Piece_EMPTY, MFLAG_EnPassant))
				}
			}
		}

		// Note that the king ending up in check is done later as part of normal filtering
		// Checks for open squares and not in check no squares
		if brd.castlePermFlag&(WKCA) != 0 {
			// F1 G1 must be empty, e1, f1 ! attacked
			if brd.pieces[SQUARES_F1] == Piece_EMPTY && brd.pieces[SQUARES_G1] == Piece_EMPTY {
				// check attached
				if brd.isSqAttacked(SQUARES_F1, Color_BLACK) == false &&
					brd.isSqAttacked(SQUARES_E1, Color_BLACK) == false {
					// printBoard("Castling")
					brd.addQuietMove(buildMOVE(SQUARES_E1, SQUARES_G1, Piece_EMPTY, Piece_EMPTY, MFLAG_Castling))
				}
			}
		}

		if brd.castlePermFlag&(WQCA) != 0 {
			// F1 G1 must be empty, e1, f1 ! attacked
			if brd.pieces[SQUARES_D1] == Piece_EMPTY &&
				brd.pieces[SQUARES_C1] == Piece_EMPTY &&
				brd.pieces[SQUARES_B1] == Piece_EMPTY {
				// check attached
				if brd.isSqAttacked(SQUARES_D1, Color_BLACK) == false &&
					brd.isSqAttacked(SQUARES_E1, Color_BLACK) == false {
					// printBoard("Castling")
					brd.addQuietMove(buildMOVE(SQUARES_E1, SQUARES_C1, Piece_EMPTY, Piece_EMPTY, MFLAG_Castling))
				}
			}
		}

	} else {
		cur_pceType = Piece_BPAWN

		for cur_pceNum = 0; cur_pceNum < brd.pceNum[cur_pceType]; cur_pceNum++ {

			sq = brd.pList[pieceIndex(cur_pceType, cur_pceNum)]

			if brd.pieces[sq-10] == Piece_EMPTY {
				brd.addBlackPawnQuietMove(sq, (sq - 10))
				if ranksBrd[sq] == RANK_7 && brd.pieces[sq-20] == Piece_EMPTY {
					brd.addQuietMove(buildMOVE(sq, (sq - 20), Piece_EMPTY, Piece_EMPTY, MFLAG_PawnStart))
				}
			}

			if brd.isSqOffBoard((sq-9)) == false && PieceCol[brd.pieces[sq-9]] == Color_WHITE {
				brd.addBlackPawnCaptureMove(sq, (sq - 9), brd.pieces[sq-9])
			}

			if brd.isSqOffBoard((sq-11)) == false && PieceCol[brd.pieces[sq-11]] == Color_WHITE {
				brd.addBlackPawnCaptureMove(sq, (sq - 11), brd.pieces[sq-11])
			}

			if brd.enPassantSq120 != SQUARES_NO_SQ {
				if sq-9 == brd.enPassantSq120 {
					brd.addEnPassantMove(buildMOVE(sq, (sq - 9), Piece_EMPTY, Piece_EMPTY, MFLAG_EnPassant))
				}
				if sq-11 == brd.enPassantSq120 {
					brd.addEnPassantMove(buildMOVE(sq, (sq - 11), Piece_EMPTY, Piece_EMPTY, MFLAG_EnPassant))
				}
			}
		}

		if brd.castlePermFlag&(BKCA) != 0 {
			if brd.pieces[SQUARES_F8] == Piece_EMPTY && brd.pieces[SQUARES_G8] == Piece_EMPTY {
				// check attached
				if brd.isSqAttacked(SQUARES_F8, Color_WHITE) == false &&
					brd.isSqAttacked(SQUARES_E8, Color_WHITE) == false {
					// printBoard("Castling")
					brd.addQuietMove(buildMOVE(SQUARES_E8, SQUARES_G8, Piece_EMPTY, Piece_EMPTY, MFLAG_Castling))
				}
			}
		}

		if brd.castlePermFlag&(BQCA) != 0 {
			if brd.pieces[SQUARES_D8] == Piece_EMPTY &&
				brd.pieces[SQUARES_C8] == Piece_EMPTY &&
				brd.pieces[SQUARES_B8] == Piece_EMPTY {
				// check attached
				if brd.isSqAttacked(SQUARES_D8, Color_WHITE) == false &&
					brd.isSqAttacked(SQUARES_E8, Color_WHITE) == false {
					// printBoard("Castling")
					brd.addQuietMove(buildMOVE(SQUARES_E8, SQUARES_C8, Piece_EMPTY, Piece_EMPTY, MFLAG_Castling))
				}
			}
		}

	}
	// End castle and pawn

	// Non sliding
	// get pce for side.   wN and wK
	// look all directions for piece.
	pceIndex = LoopNonSlideIndex[brd.side]
	pce = LoopNonSlidePce[pceIndex]
	pceIndex++

	//        println(" non sliding pieces for side    ${side}")
	for pce != Piece_EMPTY {
		//            println(" Piece to check for moves:   ${pce} ${pceChar[pce]}   # on board: ${pceNum[pce]}")
		cur_pceNum = 0
		for cur_pceNum < brd.pceNum[pce] {
			sq = brd.pList[pieceIndex(pce, cur_pceNum)]
			cur_pceNum++

			//                println("   at sq ${sq}")

			// Loop all directions
			var dirArray = PieceDirections[pce]
			for _, dir = range dirArray {
				t_sq = (sq + dir)
				//                    println("    Direction: ${dir}  target: ${t_sq}")

				if brd.isSqOffBoard(t_sq) {
					continue
				}

				if brd.pieces[t_sq] != Piece_EMPTY {
					// May be a capture
					if PieceCol[brd.pieces[t_sq]] != brd.side {
						// add capture move
						brd.addCaptureMove(buildMOVE(sq, t_sq, brd.pieces[t_sq], Piece_EMPTY, 0))
					}
				} else {
					brd.addQuietMove(buildMOVE(sq, t_sq, Piece_EMPTY, Piece_EMPTY, 0))
				}
			}
		}
		pce = LoopNonSlidePce[pceIndex]
		pceIndex++
	}

	// Sliding pieces

	// Sliding pieces
	pceIndex = LoopSlideIndex[brd.side]
	pce = LoopSlidePce[pceIndex]
	pceIndex++
	//println(" sliding pieces for side    ${side}")
	for pce != Piece_EMPTY {
		//println(" Piece to check for moves:   ${pce} ${pceChar[pce]}   # on board: ${pceNum[pce]}")
		cur_pceNum = 0
		for cur_pceNum < brd.pceNum[pce] {
			sq = brd.pList[pieceIndex(pce, cur_pceNum)]
			cur_pceNum++

			//println("   at sq ${sq}")

			// Loop all directions
			var dirArray = PieceDirections[pce]
			for _, dir = range dirArray {
				t_sq = (sq + dir)

				for !brd.isSqOffBoard(t_sq) {
					// println("    Direction: ${dir}  target: ${t_sq}")

					if brd.pieces[t_sq] != Piece_EMPTY {
						// May be a capture
						if PieceCol[brd.pieces[t_sq]] != brd.side {
							// add capture move
							brd.addCaptureMove(buildMOVE(sq, t_sq, brd.pieces[t_sq], Piece_EMPTY, 0))
						}
						// Exit the loop
						break
					}
					brd.addQuietMove(buildMOVE(sq, t_sq, Piece_EMPTY, Piece_EMPTY, 0))
					t_sq = t_sq + dir
				}
			}
		}
		pce = LoopSlidePce[pceIndex]
		pceIndex++
	}

}

/*
*
Clone of generate moves... but only does captures
*/
func generateCaptureMoves(brd *Board) {

	if DEBUG {
		brd.PrintBoard("Generate moves")
	}
	brd.moveListStart[brd.ply+1] = brd.moveListStart[brd.ply]

	var cur_pceType = 0
	var cur_pceNum = 0
	var sq SQ120 = 0
	var pceIndex = 0
	var pce = 0
	var t_sq SQ120 = 0
	var dir = 0

	if brd.side == Color_WHITE {
		cur_pceType = Piece_WPAWN
		for cur_pceNum = 0; cur_pceNum < brd.pceNum[cur_pceType]; cur_pceNum++ {
			sq = brd.pList[pieceIndex(cur_pceType, cur_pceNum)]

			//if brd.pieces[sq+10] == Piece_EMPTY {
			//	brd.addWhitePawnQuietMove(sq, sq+10)
			//
			//	if ranksBrd[sq ] == RANK_2 && brd.pieces[sq+20] == Piece_EMPTY {
			//		// Add quiet move here
			//		brd.addQuietMove(buildMOVE(sq, sq+20, Piece_EMPTY, Piece_EMPTY, MFLAG_PawnStart))
			//	}
			//}

			if brd.isSqOffBoard(sq+9) == false && PieceCol[brd.pieces[sq+9]] == Color_BLACK {
				brd.addWhitePawnCaptureMove(sq, (sq + 9), brd.pieces[sq+9])
			}

			if brd.isSqOffBoard(sq+11) == false && PieceCol[brd.pieces[sq+11]] == Color_BLACK {
				brd.addWhitePawnCaptureMove(sq, sq+11, brd.pieces[sq+11])
			}

			if brd.enPassantSq120 != SQUARES_NO_SQ {
				if sq+9 == brd.enPassantSq120 {
					brd.addEnPassantMove(buildMOVE(sq, sq+9, Piece_EMPTY, Piece_EMPTY, MFLAG_EnPassant))
				}
				if sq+11 == brd.enPassantSq120 {
					brd.addEnPassantMove(buildMOVE(sq, sq+11, Piece_EMPTY, Piece_EMPTY, MFLAG_EnPassant))
				}
			}
		}

		// Note that the king ending up in check is done later as part of normal filtering
		// Checks for open squares and not in check no squares
		//if brd.castlePermFlag & (WKCA) != 0 {
		//	// F1 G1 must be empty, e1, f1 ! attacked
		//	if  brd.pieces[SQUARES_F1 ] == Piece_EMPTY &&  brd.pieces[SQUARES_G1 ] == Piece_EMPTY  {
		//		// check attached
		//		if brd.isSqAttacked(SQUARES_F1, Color_BLACK) == false &&
		//			brd.isSqAttacked(SQUARES_E1, Color_BLACK) == false  {
		//			// printBoard("Castling")
		//			brd.addQuietMove( buildMOVE(SQUARES_E1, SQUARES_G1, Piece_EMPTY, Piece_EMPTY, MFLAG_Castling ) )
		//		}
		//	}
		//}

		//if brd.castlePermFlag&(WQCA) != 0 {
		//	// F1 G1 must be empty, e1, f1 ! attacked
		//	if  brd.pieces[SQUARES_D1 ] == Piece_EMPTY &&
		//		brd.pieces[SQUARES_C1 ] == Piece_EMPTY &&
		//		brd.pieces[SQUARES_B1 ] == Piece_EMPTY  {
		//		// check attached
		//		if brd.isSqAttacked(SQUARES_D1, Color_BLACK) == false &&
		//			brd.isSqAttacked(SQUARES_E1, Color_BLACK) == false {
		//			// printBoard("Castling")
		//			brd.addQuietMove( buildMOVE(SQUARES_E1, SQUARES_C1, Piece_EMPTY, Piece_EMPTY, MFLAG_Castling ) )
		//		}
		//	}
		//}

	} else {
		cur_pceType = Piece_BPAWN

		for cur_pceNum = 0; cur_pceNum < brd.pceNum[cur_pceType]; cur_pceNum++ {

			sq = brd.pList[pieceIndex(cur_pceType, cur_pceNum)]

			//if (brd.pieces[sq-10] == Piece_EMPTY) {
			//	brd.addBlackPawnQuietMove(sq, (sq-10) )
			//	if (ranksBrd[sq ] == RANK_7 && brd.pieces[sq-20] == Piece_EMPTY ) {
			//		brd.addQuietMove(buildMOVE(sq, (sq-20) , Piece_EMPTY, Piece_EMPTY, MFLAG_PawnStart))
			//	}
			//}

			if brd.isSqOffBoard((sq-9)) == false && PieceCol[brd.pieces[sq-9]] == Color_WHITE {
				brd.addBlackPawnCaptureMove(sq, (sq - 9), brd.pieces[sq-9])
			}

			if brd.isSqOffBoard((sq-11)) == false && PieceCol[brd.pieces[sq-11]] == Color_WHITE {
				brd.addBlackPawnCaptureMove(sq, (sq - 11), brd.pieces[sq-11])
			}

			if brd.enPassantSq120 != SQUARES_NO_SQ {
				if sq-9 == brd.enPassantSq120 {
					brd.addEnPassantMove(buildMOVE(sq, (sq - 9), Piece_EMPTY, Piece_EMPTY, MFLAG_EnPassant))
				}
				if sq-11 == brd.enPassantSq120 {
					brd.addEnPassantMove(buildMOVE(sq, (sq - 11), Piece_EMPTY, Piece_EMPTY, MFLAG_EnPassant))
				}
			}
		}

		//if brd.castlePermFlag & (BKCA) != 0 {
		//	if brd.pieces[SQUARES_F8 ] == Piece_EMPTY &&  brd.pieces[SQUARES_G8 ] == Piece_EMPTY {
		//		// check attached
		//		if brd.isSqAttacked(SQUARES_F8, Color_WHITE) == false &&
		//			brd.isSqAttacked(SQUARES_E8, Color_WHITE) == false {
		//			// printBoard("Castling")
		//			brd.addQuietMove( buildMOVE(SQUARES_E8, SQUARES_G8, Piece_EMPTY, Piece_EMPTY, MFLAG_Castling ) )
		//		}
		//	}
		//}
		//
		//if brd.castlePermFlag & (BQCA) != 0 {
		//	if brd.pieces[SQUARES_D8 ] == Piece_EMPTY &&
		//		brd.pieces[SQUARES_C8 ] == Piece_EMPTY &&
		//		brd.pieces[SQUARES_B8 ] == Piece_EMPTY {
		//		// check attached
		//		if brd.isSqAttacked(SQUARES_D8, Color_WHITE) == false &&
		//			brd.isSqAttacked(SQUARES_E8, Color_WHITE) == false {
		//			// printBoard("Castling")
		//			brd.addQuietMove( buildMOVE(SQUARES_E8, SQUARES_C8, Piece_EMPTY, Piece_EMPTY, MFLAG_Castling ) )
		//		}
		//	}
		//}

	}
	// End castle and pawn

	// Non sliding
	// get pce for side.   wN and wK
	// look all directions for piece.
	pceIndex = LoopNonSlideIndex[brd.side]
	pce = LoopNonSlidePce[pceIndex]
	pceIndex++

	//        println(" non sliding pieces for side    ${side}")
	for pce != Piece_EMPTY {
		//            println(" Piece to check for moves:   ${pce} ${pceChar[pce]}   # on board: ${pceNum[pce]}")
		cur_pceNum = 0
		for cur_pceNum < brd.pceNum[pce] {
			sq = brd.pList[pieceIndex(pce, cur_pceNum)]
			cur_pceNum++

			//                println("   at sq ${sq}")

			// Loop all directions
			var dirArray = PieceDirections[pce]
			for _, dir = range dirArray {
				t_sq = (sq + dir)
				//                    println("    Direction: ${dir}  target: ${t_sq}")

				if brd.isSqOffBoard(t_sq) {
					continue
				}

				if brd.pieces[t_sq] != Piece_EMPTY {
					// May be a capture
					if PieceCol[brd.pieces[t_sq]] != brd.side {
						// add capture move
						brd.addCaptureMove(buildMOVE(sq, t_sq, brd.pieces[t_sq], Piece_EMPTY, 0))
					}
				} else {
					//brd.addQuietMove( buildMOVE(sq, t_sq, Piece_EMPTY, Piece_EMPTY, 0 ) )
				}
			}
		}
		pce = LoopNonSlidePce[pceIndex]
		pceIndex++
	}

	// Sliding pieces

	// Sliding pieces
	pceIndex = LoopSlideIndex[brd.side]
	pce = LoopSlidePce[pceIndex]
	pceIndex++
	//println(" sliding pieces for side    ${side}")
	for pce != Piece_EMPTY {
		//println(" Piece to check for moves:   ${pce} ${pceChar[pce]}   # on board: ${pceNum[pce]}")
		cur_pceNum = 0
		for cur_pceNum < brd.pceNum[pce] {
			sq = brd.pList[pieceIndex(pce, cur_pceNum)]
			cur_pceNum++

			//println("   at sq ${sq}")

			// Loop all directions
			var dirArray = PieceDirections[pce]
			for _, dir = range dirArray {
				t_sq = (sq + dir)

				for !brd.isSqOffBoard(t_sq) {
					// println("    Direction: ${dir}  target: ${t_sq}")

					if brd.pieces[t_sq] != Piece_EMPTY {
						// May be a capture
						if PieceCol[brd.pieces[t_sq]] != brd.side {
							// add capture move
							brd.addCaptureMove(buildMOVE(sq, t_sq, brd.pieces[t_sq], Piece_EMPTY, 0))
						}
						// Exit the loop
						break
					}
					//brd.addQuietMove( buildMOVE(sq, t_sq, Piece_EMPTY, Piece_EMPTY, 0 ) )
					t_sq = t_sq + dir
				}
			}
		}
		pce = LoopSlidePce[pceIndex]
		pceIndex++
	}

}

func moveExists(move Move, board *Board) bool {
	generateMoves(board)

	var moveFound = NOMOVE
	for idx := board.moveListStart[board.ply]; idx <= board.moveListStart[board.ply+1]-1; idx++ {
		moveFound = move
		//board.moveList[idx]
		if !board.makeMove(moveFound) {
			continue
		}
		board.takeMove()
		if move == moveFound {
			return true
		}
	}
	return false
}

// GetAllValidMoves
// Returns all the valid moves.
func GetAllValidMoves(board *Board) []int {

	generateMoves(board)
	// board.printMoveList()

	moves := make([]int, 0)

	for idx := board.moveListStart[board.ply]; idx < board.moveListStart[board.ply+1]; idx++ {
		move := board.moveList[idx]
		if board.makeMove(move) == false {
			// This is a valid case.  When king is in check
			// println("  make move ", MoveToString(move) , "   at ", p.perft_leafNodes, " Not valid")
			continue
		}
		moves = append(moves, move)
		board.takeMove()

	}
	return moves
}

// GetAllMoves
// Returns all the moves, some are invalid ex, if king goes into check
func GetAllMoves(board *Board) []int {

	generateMoves(board)
	// board.printMoveList()

	moves := make([]int, 0)

	for idx := board.moveListStart[board.ply]; idx < board.moveListStart[board.ply+1]; idx++ {
		move := board.moveList[idx]
		//if board.makeMove(move) == false {
		//	// This is a valid case.  When king is in check
		//	// println("  make move ", MoveToString(move) , "   at ", p.perft_leafNodes, " Not valid")
		//	continue
		//}
		moves = append(moves, move)
		//board.takeMove()

	}
	return moves
}

func MoveToSimpleString(move Move) string {
	var fSq120 = getFromSq120(move)
	var tSq120 = getToSq120(move)
	return fmt.Sprintf("%s%s", sqToString(fSq120), sqToString(tSq120))
}

func MoveToInputString(move Move) string {

	var fSq120 = getFromSq120(move)
	var tSq120 = getToSq120(move)
	var promoted = getPromoted(move)
	var promotedStr = ""
	if Piece_EMPTY != promoted {
		promotedStr = " (" + s.ToUpper(string(PceCharLetter[promoted])) + ")"
	}
	return fmt.Sprintf("%s%s%s", sqToString(fSq120), sqToString(tSq120), promotedStr)
}
