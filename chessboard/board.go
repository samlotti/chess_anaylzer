package ai

import (
	"bytes"
	"fmt"
	_ "math/rand"
	"strings"
)

/**
Chess board that can tell valid moves, make the moves ...
*/

const (
	//  The square number for the piece.  Each piece takes 10 slots.  The # of each WP(example) is max of 10.
	//  for first 10 are empty, next ten will contain the positions of the pawns on the board. (0->8)
	//  the next 10 are for knights ..
	PLIST_SIZE = 14 * 10

	//  This an pList are so we dont have to iterate the entire board looking for the pieces.
	//  Just get the count in the pceNum and loop in the pList
	//  indexed by piece code (# of each piece)
	//  optimization when calulating the moves
	PCENUM_SIZE = 13

	MATERIAL_SIZE = 2
)

//
type HistoryData struct {
	move           Move
	castlePermFlag int
	enPassSq120    SQ120
	fiftyMode      int
	posKey         int
}

type Board struct {

	// Need a 120 cell slice
	pieces    []Piece
	side      Side
	fiftyMove int

	// how many half moves in the game. allows us to undo
	// how deap in the move search
	hisPly         int // 1/2 moves for entire game.
	ply            int // Number of 1/2 moves in search tree  0 at start of search
	castlePermFlag int
	material       []int // WHITE, BLACK material of piece (value of all material)

	history [maxGameMovesHistory]HistoryData

	// The capture square for enpassant
	enPassantSq120 SQ120

	// number by piece type (# on the board)
	pceNum []int

	//  The square number for the piece.  Each piece takes 10 slots.  The # of each WP(example) is max of 10.
	//  for first 10 are empty, next ten will contain the positions of the pawns on the board. (0->8)
	//  the next 10 are for knights ..
	pList []SQ120 // PLIST_SIZE

	//  The hash of this position
	posKey int

	// Where the move starts at a given Depth. (?? index lookup ??)
	moveListStart [maxDepth]int

	// The moves at each ply (searching)
	moveList [maxDepth * maxPositionMoves]Move

	// Scores for each move (searching)
	moveScores [maxDepth * maxPositionMoves]int

	searchKillers [3 * maxDepth]int

	searchHistory [14 * BRD_SQ_NUM]Move

	//// Table for the game (search)
	//// This is posKey % PVENTRIES.  It allows the engine to keep
	//// the Best move for each position.
	//// This allows iterative deepening to take Best from previous searches
	//pvtable [PVENTRIES]PVEntry

	// Best moves for each Depth
	// This is calculated after the search is completed.
	// var pvArray = Array(maxDepth) { NOMOVE }
	pvArray [maxDepth]Move
}

/**
Initialize everything about the board
*/
func (this *Board) init() {
	this.pieces = make([]Piece, BRD_SQ_NUM)

	for i := 0; i < BRD_SQ_NUM; i++ {
		this.pieces[i] = SQUARES_OFFBOARD
	}

	// Isnt even needed, go sets it to empty history data items.
	for i := 0; i < maxGameMoves; i++ {
		this.history[i] = HistoryData{}
	}

	this.side = Color_WHITE
	this.fiftyMove = 0
	this.hisPly = 0
	this.ply = 0
	this.castlePermFlag = 0
	this.material = make([]int, MATERIAL_SIZE)
	this.enPassantSq120 = 0
	this.pList = make([]SQ120, PLIST_SIZE)
	this.pceNum = make([]int, PCENUM_SIZE)
}

func NewBoard() *Board {

	if DEBUG {
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
		fmt.Println("WARNING: DEBUG flag is on in the definitions.go.")
	}

	b := new(Board)
	b.init()
	return b
}

func (this *Board) isSqOffBoard(sq120 SQ120) bool {
	return this.pieces[sq120] == SQUARES_OFFBOARD
}

/**
* Re-Generate a position key based on the values on the baord
 */
func (this *Board) generatePosKey() int {

	finalKey := 0
	piece := Piece_EMPTY
	pieces1 := this.pieces

	for sq := 0; sq < BRD_SQ_NUM; sq++ {
		piece = pieces1[sq]
		if piece != Piece_EMPTY && piece != SQUARES_OFFBOARD {
			finalKey = finalKey ^ (PieceKeys[(piece*120)+sq])
		}
	}

	if this.side == Color_WHITE {
		finalKey = finalKey ^ (SideKey)
	}

	if this.enPassantSq120 != SQUARES_NO_SQ {
		finalKey = finalKey ^ (PieceKeys[this.enPassantSq120])
	}

	finalKey = finalKey ^ (CastleKeys[this.castlePermFlag])

	return finalKey
}

func (this *Board) ResetBoard() {

	for sq120 := 0; sq120 < BRD_SQ_NUM; sq120++ {
		this.pieces[sq120] = SQUARES_OFFBOARD
	}

	// The inner 8x8 board to empty
	for i := 0; i < 64; i++ {
		this.pieces[sq_120(i)] = Piece_EMPTY
	}

	this.side = Color_BOTH
	this.enPassantSq120 = SQUARES_NO_SQ
	this.fiftyMove = 0
	this.ply = 0
	this.hisPly = 0
	this.castlePermFlag = 0
	this.posKey = 0
	this.moveListStart[this.ply] = 0

}

// chkSide = 0 or 1
func (this *Board) printSquareAttacked(chkSide Side) {

	// var chkSide = getOtherSide( side )

	fmt.Println("Squares attacked by ", string(SideChar[chkSide]))

	for rank := RANK_8; rank >= RANK_1; rank-- {

		fmt.Print(string(RankChar[rank]), "  ")

		for file := FILE_A; file <= FILE_H; file++ {
			sq120 := fr_to_SQ120(file, rank)
			if this.isSqAttacked(sq120, chkSide) {
				fmt.Print(" X ")
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Println("")
	}
	fmt.Println("    a  b  c  d  e  f  g  h")
}

func (this *Board) printBoard(msg string) {
	fmt.Print("\nGame Board: ", msg, "\n\n")
	fmt.Println("EP:", this.enPassantSq120, " 8x8:", sqToString(this.enPassantSq120))

	buffer := bytes.Buffer{}
	for rank := RANK_8; rank >= RANK_1; rank-- {
		buffer.Truncate(0)
		buffer.WriteString("" + string(RankChar[rank]) + "  ")
		for file := FILE_A; file <= FILE_H; file++ {
			sq120 := fr_to_SQ120(file, rank)
			piece := this.pieces[sq120]
			buffer.WriteString("  " + string(PceCharLetter[piece]) + "  ")
		}

		fmt.Println(buffer.String())
		//append(encodeT)
	}
	fmt.Println("     A    B    C    D    E    F    G    H")

	fmt.Println("side: ", SideChar[this.side])
	fmt.Println("enPass: ", this.enPassantSq120)

	buffer.Truncate(0)
	buffer.WriteString("Castle: ")
	if this.castlePermFlag&WKCA != 0 {
		buffer.WriteString("K")
	}
	if this.castlePermFlag&WQCA != 0 {
		buffer.WriteString("Q")
	}
	if this.castlePermFlag&BKCA != 0 {
		buffer.WriteString("k")
	}
	if this.castlePermFlag&BQCA != 0 {
		buffer.WriteString("q")
	}
	fmt.Println(buffer.String())

	fmt.Println("Pos Key:", this.posKey)
}

func (this *Board) PrintPieceList() {
	for piece := Piece_WPAWN; piece <= Piece_BKING; piece++ {
		for idx := 0; idx < this.pceNum[piece]; idx++ {
			fmt.Println("Piece ", string(PceCharLetter[piece]), " on ", sqToString(this.pList[pieceIndex(piece, idx)]))
		}
	}
	fmt.Println("Material W ", this.material[Color_WHITE], "   B ", this.material[Color_BLACK])
}

/**
* Updates the piecelist and materials. This will look at all the
* positions in the board. and fill in these arrays
 */
func (this *Board) updateListsMaterial() {
	// Reset all pieces by index
	for i := 0; i < PLIST_SIZE; i++ {
		this.pList[i] = EMPTY_SQ120
	}
	// Reset material value
	for i := 0; i < MATERIAL_SIZE; i++ {
		this.material[i] = 0
	}

	// Reset number of pieces
	for i := 0; i < PCENUM_SIZE; i++ {
		this.pceNum[i] = 0
	}

	for sq64 := 0; sq64 < 64; sq64++ {
		sq120 := Sq64To120[sq64]
		piece := this.pieces[sq120]
		if piece != Piece_EMPTY {
			//fmt.Println("Piece: ", piece, " on ", sqToString(sq120))
			col := PieceCol[piece]
			this.material[col] += PieceVal[piece]
			// fmt.Println("material: ", this.material)
			this.pList[pieceIndex(piece, this.pceNum[piece])] = sq120
			this.pceNum[piece]++
		}
	}

}

/*
Returns the piece index for the piece
in the plist table
*/
func pieceIndex(pieceCode int, pieceNum int) int {
	return (pieceCode * 10) + pieceNum
}

/**
 * Will check to see if this square is being attacked by the side
 *
 *   sq + 12 = up
 *   sq - 1 = left
 *   sq + 1 = right
 *   sq - 12 = down
 *
 *   +11 = upper right
 *   +9  = upper left
 *
 *   -11 = lower left
 *   -9 = lower right
 *
 *
 *  Note because of the way the board is created there is a border that has off side codes so we dont have to worry about
 *  out of bounds
 *
 *  This is called many time .. perfoamnce issue with the Integer.toInt
 *

(Running perf_test.go)

     2.43s 27.65% 27.65%      2.43s 27.65%  hapticapps/modules/chess/ai.(*Board).isSqAttacked
     1.18s 13.42% 41.07%      1.18s 13.42%  runtime.usleep
     0.71s  8.08% 49.15%      1.27s 14.45%  hapticapps/modules/chess/ai.(*Board).movePiece
     0.60s  6.83% 55.97%      3.91s 44.48%  hapticapps/modules/chess/ai.(*Board).makeMove

Move/Sec:  9598542   Nodes: 119060324  Duration Sec: 12.404
Move/Sec:  9628038   Nodes: 119060324  Duration Sec: 12.366
Move/Sec:  9612492   Nodes: 119060324  Duration Sec: 12.386
--- PASS: TestPerfs (12.95s)

--Remove ranges (slows down)  <reverted>
Move/Sec:  9245967   Nodes: 119060324  Duration Sec: 12.877
Move/Sec:  9225906   Nodes: 119060324  Duration Sec: 12.905
--- FAIL: TestPerfs (13.34s)

-- altered the conditionals to check color last.
Move/Sec:  9686001   Nodes: 119060324  Duration Sec: 12.292
Move/Sec:  9704158   Nodes: 119060324  Duration Sec: 12.269
Move/Sec:  9696255   Nodes: 119060324  Duration Sec: 12.279

*/
func (this *Board) isSqAttacked(sq120 SQ120, side Side) bool {

	// Pawn attacks
	if side == Color_WHITE {
		if this.pieces[sq120-11] == Piece_WPAWN ||
			this.pieces[sq120-9] == Piece_WPAWN {
			return true
		}
	} else {
		if this.pieces[sq120+11] == Piece_BPAWN ||
			this.pieces[sq120+9] == Piece_BPAWN {
			return true
		}
	}

	// check knight
	for _, dir := range KnightDirection {
		//for idx:=0; idx<len(KnightDirection); idx++  {
		//	var dir = KnightDirection[idx]
		//fmt.Println("knight: ", dir )
		var pce = this.pieces[sq120+dir]
		if pce != SQUARES_OFFBOARD && PieceKnight[pce] == true && PieceCol[pce] == side {
			return true
		}
	}

	// check rook / queen
	for _, dir := range RookDirection {
		//for idx:=0; idx<len(RookDirection); idx++  {
		//	var dir = RookDirection[idx]

		//fmt.Println("rook/queen: ", dir )
		var t_sq = sq120 + dir
		var pce = this.pieces[t_sq]
		for pce != SQUARES_OFFBOARD {
			if pce != Piece_EMPTY {
				if PieceCol[pce] == side && PieceRookQueen[pce] {
					return true
				}
				break
			}
			t_sq += dir
			pce = this.pieces[t_sq]
		}
	}

	// check bishop / queen
	for _, dir := range BishopDirection {
		//for idx:=0; idx<len(BishopDirection); idx++  {
		//	var dir = BishopDirection[idx]

		//fmt.Println("bishop/queen: ", dir )
		var t_sq = sq120 + dir
		var pce = this.pieces[t_sq]
		for pce != SQUARES_OFFBOARD {
			if pce != Piece_EMPTY {
				if PieceBishopQueen[pce] && PieceCol[pce] == side {
					return true
				}
				break
			}
			t_sq += dir
			pce = this.pieces[t_sq]
		}
	}

	// check king
	for _, dir := range KingDirection {
		//for idx:=0; idx<len(KingDirection); idx++  {
		//	var dir = KingDirection[idx]

		//fmt.Println("king: ", dir )
		var pce = this.pieces[sq120+dir]
		if pce != SQUARES_OFFBOARD && PieceKing[pce] == true && PieceCol[pce] == side {
			return true
		}
	}

	return false
}

func (brd *Board) addQuietMove(move Move) {
	// fmt.Printf("Add quiet move: %s\n", MoveToString(move))
	var idx = brd.moveListStart[brd.ply+1]
	brd.moveList[idx] = move
	brd.moveScores[idx] = 0

	// For quiet moves .. assign killer points
	if move == brd.searchKillers[brd.ply] {
		// killer 1 hit
		brd.moveScores[idx] = 900000
	} else {
		if move == brd.searchKillers[brd.ply+maxDepth] {
			brd.moveScores[idx] = 800000
		} else {
			brd.moveScores[idx] =
				brd.searchHistory[brd.pieces[getFromSq120(move)]*BRD_SQ_NUM+getToSq120(move)]
		}
	}

	brd.moveListStart[brd.ply+1]++
}

func (brd *Board) addWhitePawnQuietMove(fromSq120 SQ120, toSq120 SQ120) {
	if ranksBrd[fromSq120] == RANK_7 {
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_WQUEEN, 0))
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_WROOK, 0))
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_WBISHOP, 0))
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_WKNIGHT, 0))
	} else {
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_EMPTY, 0))
	}
}

func (brd *Board) addBlackPawnQuietMove(fromSq120 SQ120, toSq120 SQ120) {
	if ranksBrd[fromSq120] == RANK_2 {
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_BQUEEN, 0))
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_BROOK, 0))
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_BBISHOP, 0))
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_BKNIGHT, 0))
	} else {
		brd.addQuietMove(buildMOVE(fromSq120, toSq120, Piece_EMPTY, Piece_EMPTY, 0))
	}
}

func (brd *Board) addCaptureMove(move Move) {
	// println("Add capture move: ${MoveToString(move)}")
	var idx = brd.moveListStart[brd.ply+1]
	brd.moveList[idx] = move

	// To support move ordering
	brd.moveScores[idx] = mvvLvaScores[getCapturedPiece(move)*14+brd.pieces[getFromSq120(move)]] + 1000000

	brd.moveListStart[brd.ply+1]++
}

func (brd *Board) addBlackPawnCaptureMove(fromSq120 SQ120, toSq120 SQ120, capPiece Piece) {
	if ranksBrd[fromSq120] == RANK_2 {
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_BQUEEN, 0))
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_BROOK, 0))
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_BBISHOP, 0))
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_BKNIGHT, 0))
	} else {
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_EMPTY, 0))
	}
}

func (brd *Board) addWhitePawnCaptureMove(fromSq120 SQ120, toSq120 SQ120, capPiece Piece) {
	if ranksBrd[fromSq120] == RANK_7 {
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_WQUEEN, 0))
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_WROOK, 0))
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_WBISHOP, 0))
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_WKNIGHT, 0))
	} else {
		brd.addCaptureMove(buildMOVE(fromSq120, toSq120, capPiece, Piece_EMPTY, 0))
	}
}

func (brd *Board) addEnPassantMove(move Move) {
	// fmt.Println("Add EnPassant move: ", MoveToString(move) )
	var idx = brd.moveListStart[brd.ply+1]
	brd.moveList[idx] = move
	brd.moveScores[idx] = 105 + 1000000
	// mvvLvaScores[pawn takes pawn = 105] + 1_000_000
	brd.moveListStart[brd.ply+1]++
}

func (brd *Board) printMoveList() {
	for idx := brd.moveListStart[brd.ply]; idx < brd.moveListStart[brd.ply+1]; idx++ {
		move := brd.moveList[idx]
		// println("${idx+1}: ${MoveToString(move)}   ${move}")
		fmt.Println(idx+1, " = ", MoveToString(move))
	}
}

func (brd *Board) hashPiece(piece Piece, sq120 SQ120) {
	brd.posKey = brd.posKey ^ (PieceKeys[(piece*120)+sq120])
}

func (brd *Board) hashCastle() {
	brd.posKey = brd.posKey ^ (CastleKeys[brd.castlePermFlag])
}

func (brd *Board) hashSide() {
	brd.posKey = brd.posKey ^ (SideKey)
}

func (brd *Board) hashEnPassant() {
	brd.posKey = brd.posKey ^ (PieceKeys[brd.enPassantSq120])
}

/**
 * Remove a piece from the square
 */
//func (brd *Board) clearPiece(sq120 SQ120, move Move = 0 ) {
func (brd *Board) clearPiece(sq120 SQ120) {

	var pce = brd.pieces[sq120]
	var col = PieceCol[pce]
	var index = 0
	var t_pceNum = -1

	// remove the hash
	brd.hashPiece(pce, sq120)

	brd.pieces[sq120] = Piece_EMPTY

	//        if ( DEBUG ) {
	//            if (col == 2) {
	//                println("Move: ${moveToDebugString(move)}")
	//                printBoard()
	//            }
	//        }
	brd.material[col] -= PieceVal[pce]

	// remove from the piece list
	// remove from array by overwriting with the last entry.
	//
	for index < brd.pceNum[pce] {
		if brd.pList[pieceIndex(pce, index)] == sq120 {
			t_pceNum = index
			break
		}
		index++
	}

	brd.pceNum[pce] -= 1
	brd.pList[pieceIndex(pce, t_pceNum)] = brd.pList[pieceIndex(pce, brd.pceNum[pce])]

}

func (brd *Board) addPiece(sq120 SQ120, pce Piece) {
	var col = PieceCol[pce]

	// remove the hash
	brd.hashPiece(pce, sq120)

	brd.pieces[sq120] = pce

	brd.material[col] += PieceVal[pce]

	brd.pList[pieceIndex(pce, brd.pceNum[pce])] = sq120
	brd.pceNum[pce] += 1

}

func (brd *Board) movePiece(from_sq120 SQ120, to_sq120 SQ120) {
	var pce = brd.pieces[from_sq120]
	// var col = PieceCol[pce]

	// remove the hash
	brd.hashPiece(pce, from_sq120)
	brd.pieces[from_sq120] = Piece_EMPTY

	brd.hashPiece(pce, to_sq120)
	brd.pieces[to_sq120] = pce

	var index = 0
	for index < brd.pceNum[pce] {
		if brd.pList[pieceIndex(pce, index)] == from_sq120 {
			brd.pList[pieceIndex(pce, index)] = to_sq120
			break
		}
		index++
	}
}

func (brd *Board) checkBoard() bool {

	var t_pceNum = [13]int{0}
	var t_material = [3]int{0, 0, 0}
	var sq120 SQ120 = 0
	var ok = true

	// Check the piece lists
	for t_piece := Piece_WPAWN; t_piece <= Piece_BKING; t_piece++ {
		for t_pce_num := 0; t_pce_num < brd.pceNum[t_piece]; t_pce_num++ {
			sq120 = brd.pList[pieceIndex(t_piece, t_pce_num)]
			if brd.pieces[sq120] != t_piece {
				fmt.Println("Piece not same as piece list, board: ", sq120, " has ", brd.pieces[sq120], " expected ", t_piece)
				ok = false
			}
		}
	}

	// check board piece in piece list.  Gets covered from other checks

	// piece num and material
	var t_piece = 0
	for sq64 := 0; sq64 < 64; sq64++ {
		sq120 = sq_120(sq64)
		// fmt.Println("sq120:", sq120)
		t_piece = brd.pieces[sq120]
		// fmt.Println("T:", t_piece)
		t_pceNum[t_piece]++
		t_material[PieceCol[t_piece]] += PieceVal[t_piece]
	}

	for t_piece := Piece_WPAWN; t_piece <= Piece_BKING; t_piece++ {
		if t_pceNum[t_piece] != brd.pceNum[t_piece] {
			fmt.Println("Piece num not same, piece ", t_piece, " has ", t_pceNum[t_piece], " on the board, but pceNum has it as ", brd.pceNum[t_piece])
			ok = false
		}
	}

	if t_material[Color_WHITE] != brd.material[Color_WHITE] ||
		t_material[Color_BLACK] != brd.material[Color_BLACK] {
		ok = false
		fmt.Println("Material score difference WHITE ", t_material[Color_WHITE], " vs ", brd.material[Color_WHITE])
		fmt.Println("Material score difference WHITE ", t_material[Color_BLACK], " vs ", brd.material[Color_BLACK])
	}

	if brd.generatePosKey() != brd.posKey {
		ok = false
		fmt.Println("Position key generation not same as game board posKey.. This is the hash key.")
	}

	if brd.side != Color_WHITE && brd.side != Color_BLACK {
		ok = false
		fmt.Println("The side is not valid, found ", brd.side, " expected ", Color_BLACK, " or ", Color_WHITE)
	}

	if !ok {
		brd.printBoard("Board at time of check failure")
		panic("Errors in the check board")
	}
	//        println("Checkboard: ${ok}")
	return ok
}

// makeMove - make the move on the board
// can be undon with the takeMove
func (brd *Board) makeMove(move Move) bool {
	var from120 = getFromSq120(move)
	var to120 = getToSq120(move)
	var t_side = brd.side

	//if ( DEBUG ) {
	//printBoard("   about to move: ${moveToDebugString(move)}")
	//}

	// println("makeMove: ${moveToDebugString(move)}")

	// println("brd.hisPly ", brd.hisPly)
	// println("brd.PLY ", brd.ply)
	var historyData = &brd.history[brd.hisPly]

	// Store the poskey  (a)
	historyData.posKey = brd.posKey

	// Remove the enpassant capture
	if move&(MFLAG_EnPassant) != 0 {
		if brd.side == Color_WHITE {
			brd.clearPiece((to120 - 10))
		} else {
			brd.clearPiece((to120 + 10))
		}
		// (1)
	} else {
		// Castling move .. move the rook
		if move&(MFLAG_Castling) != 0 {
			switch to120 {
			case SQUARES_C1:
				brd.movePiece(SQUARES_A1, SQUARES_D1)
			case SQUARES_C8:
				brd.movePiece(SQUARES_A8, SQUARES_D8)
			case SQUARES_G1:
				brd.movePiece(SQUARES_H1, SQUARES_F1)
			case SQUARES_G8:
				brd.movePiece(SQUARES_H8, SQUARES_F8)
			default:
				panic("Error in castling move")
			}
		}
	}

	// Reset it .. can I do this in step (1)
	if brd.enPassantSq120 != SQUARES_NO_SQ {
		brd.hashEnPassant()
	}

	// Save this move info ... can this be done in (a)
	historyData.move = move
	historyData.fiftyMode = brd.fiftyMove
	historyData.enPassSq120 = brd.enPassantSq120
	historyData.castlePermFlag = brd.castlePermFlag

	brd.hashCastle()
	// handles rook moves, rook capture, king move
	brd.castlePermFlag = brd.castlePermFlag & (CastlePerm[from120]) & (CastlePerm[to120])
	brd.hashCastle()

	brd.enPassantSq120 = SQUARES_NO_SQ

	var capturedPiece = getCapturedPiece(move)

	brd.fiftyMove++

	if capturedPiece != Piece_EMPTY {
		brd.clearPiece(to120)
		brd.fiftyMove = 0
	}

	brd.hisPly++
	brd.ply++

	// Handle pawn moves
	if PiecePawn[brd.pieces[from120]] {
		brd.fiftyMove = 0

		// Pawn moves 2 make enpassant squares
		if move&(MFLAG_PawnStart) != 0 {
			if brd.side == Color_WHITE {
				brd.enPassantSq120 = (from120 + 10)
			} else {
				brd.enPassantSq120 = (from120 - 10)
			}
			brd.hashEnPassant()
		}
	}

	brd.movePiece(from120, to120)

	var prPiece = getPromoted(move)

	if prPiece != Piece_EMPTY {
		brd.clearPiece(to120)
		brd.addPiece(to120, prPiece)
	}

	brd.side = t_side ^ (1)
	brd.hashSide()

	// DEBUG.. REMOVE in FUTURE
	if DEBUG {
		brd.checkBoard()
	}

	// If I knew the king was not in check I should be able to optimise
	/*
		If king not in check,
			king move -- full atttack check
			Pawn move -- full attack check
			Others ... only do a check if the item is lined up .. horiz,vert or diaginals.
				no need to check otherwise
		If king is in check...
			moves that could not uncheck, auto fail.
				( * not a capture
				* not a king move
				* not end up on hor, vert or diag. ) = no need to check for attacked, it still is.
	*/
	if brd.isSqAttacked(brd.pList[pieceIndex(Kings[t_side], 0)], brd.side) {
		brd.takeMove()
		return false
	}
	return true
}

// Remove the move
func (brd *Board) takeMove() {

	if DEBUG {
		brd.checkBoard()
	}

	brd.hisPly--
	brd.ply--

	var move = brd.history[brd.hisPly].move
	var from120 = getFromSq120(move)
	var to120 = getToSq120(move)

	if brd.enPassantSq120 != SQUARES_NO_SQ {
		brd.hashEnPassant()
	}
	brd.hashCastle()

	brd.castlePermFlag = brd.history[brd.hisPly].castlePermFlag
	brd.fiftyMove = brd.history[brd.hisPly].fiftyMode
	brd.enPassantSq120 = brd.history[brd.hisPly].enPassSq120

	if brd.enPassantSq120 != SQUARES_NO_SQ {
		brd.hashEnPassant()
	}
	brd.hashCastle()

	brd.side = brd.side ^ (1)
	brd.hashSide()

	// reverse ep capture
	if move&(MFLAG_EnPassant) != 0 {
		if brd.side == Color_WHITE {
			brd.addPiece((to120 - 10), Piece_BPAWN)
		} else {
			brd.addPiece((to120 + 10), Piece_WPAWN)
		}
	} else {
		if move&(MFLAG_Castling) != 0 {
			switch to120 {

			case SQUARES_C1:
				{
					brd.movePiece(SQUARES_D1, SQUARES_A1)
				}
			case SQUARES_C8:
				{
					brd.movePiece(SQUARES_D8, SQUARES_A8)
				}
			case SQUARES_G1:
				{
					brd.movePiece(SQUARES_F1, SQUARES_H1)
				}
			case SQUARES_G8:
				{
					brd.movePiece(SQUARES_F8, SQUARES_H8)
				}
			default:
				{
					println("Error in castling move reversal")
				}
			}
		}
	}

	brd.movePiece(to120, from120)

	var capPiece = getCapturedPiece(move)
	if capPiece != Piece_EMPTY {
		brd.addPiece(to120, capPiece)
	}

	var promoted = getPromoted(move)
	if promoted != Piece_EMPTY {
		brd.clearPiece(from120)
		if PieceCol[promoted] == Color_WHITE {
			brd.addPiece(from120, Piece_WPAWN)
		} else {
			brd.addPiece(from120, Piece_BPAWN)
		}
	}

	if DEBUG {
		brd.checkBoard()
	}

}
func (this *Board) IsInCheck() bool {

	return this.isSqAttacked(this.pList[pieceIndex(Kings[this.side], 0)], this.side^(1))

}

/**
Makes the move based on the input string
FSTSp

Return the actual move if it is valid.  The move is applied to the board
If invalid, return error and NOMOVE for move
*/
func (this *Board) GetInternalMoveValueFromInputString(input string) (Move, error) {

	// Reset the ply ...
	this.ply = 0

	moves := GetAllValidMoves(this)
	for _, mv := range moves {
		moveStr := MoveToInputString(mv)
		// fmt.Println("M:",moveStr)
		//if moveStr == input {
		if strings.EqualFold(moveStr, input) {
			return mv, nil
		}
	}
	return NOMOVE, fmt.Errorf("Invalid move: %s", input)
}

//
// MakeMove
// Perform the move.
// The input is the string version of the move as an input - for error message
//
func (this *Board) MakeMove(mv Move, input string) error {

	// Reset the ply ...
	this.ply = 0

	// Needed for 50 move
	// wthis.hisPly = 0

	if !this.makeMove(mv) {
		return fmt.Errorf("Invalid move: %s", input)
	} else {
		return nil
	}

}

/**
 * Look for a repetition
 *
 * since fiftymove setting to 0 stops and chance of repetition (due to pawn or capture) just need to go back to that point in time
 *
 */
func (this *Board) IsRepetition() bool {
	count := 0
	for idx := (this.hisPly - this.fiftyMove); idx < (this.hisPly - 1); idx++ {
		// println("Rep: ", idx)
		if this.posKey == this.history[idx].posKey {
			count++
			if count == 3 {
				return true
			}
		}
	}
	return false
}

func (this *Board) IsFiftyMove() bool {
	return this.fiftyMove > 99
}
