package ai

import (
	"fmt"
	"math/rand"
	"time"
)

//
//
const DEBUG = false

// ----------- CONSTANTS FOR CHESS AI SERVER WITH LIMITED USERS -------------
// How many moves we can support in a game
//const maxGameMoves = 2048

// This cannot be reducedfor 50 move calculation
const maxGameMovesHistory = 2048

// Max number of moves that can be available for a position
const maxPositionMoves = 256

// How deap we can search (max)
//const maxDepth = 64

// ----------- CONSTANTS FOR CHESS SERVER WITH MANY USERS NO AI -------------
// Minimal size... ai may not be strong .. for chess server with no AI
const maxGameMoves = 10
const maxDepth = 45

// 64 = 520k memory, 4-272k
// 5,2 = 182k  => 100k games ~2gig -- Test cases fail (search)
// 2048, 64 => 100k games ~10g
// 10,30 = 379592, 100k ~= 2.3 gig -- This is the min for test cases to pass.
//const MaxGameMoves = maxGameMoves

/*
Type information
*/
type Square = int
type File = int
type Rank = int
type Color = int
type SQ = int
type SQ120 = int
type SQ64 = int
type Piece = int
type Side = int
type Move = int

const StartFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
const RankChar = "12345678"
const FileChar = "abcdefgh"
const PceCharLetter = ".PNBRQKpnbrqk"

// Note that due to utf8 storage, this cannot be indexed by a simple offset
const PceChar = ".♙♘♗♖♕♔♟♞♝♜♛♚"

const SideChar = "wb."

const PVENTRIES = 10000

const INFINITE = 30000
const MATE = 29000

/**
The piece and values
*/
const EMPTY_SQ120 SQ120 = 0
const PIECES_EMPTY = 0
const BRD_SQ_NUM = 120

const (
	Piece_EMPTY   = 0
	Piece_WPAWN   = 1
	Piece_WKNIGHT = 2
	Piece_WBISHOP = 3
	Piece_WROOK   = 4
	Piece_WQUEEN  = 5
	Piece_WKING   = 6

	Piece_BPAWN   = 7
	Piece_BKNIGHT = 8
	Piece_BBISHOP = 9
	Piece_BROOK   = 10
	Piece_BQUEEN  = 11
	Piece_BKING   = 12
)

var mvvLvaScores [14 * 14]int

var mvvLvaValue = []int{
	0, 100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600,
}

/**
File constants
*/
const ( // iota is reset to 0
	FILE_A File = iota // c0 == 0
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
	FILE_NONE
)

/**
Rank constants
*/
const ( // iota is reset to 0
	RANK_1 Rank = iota // c0 == 0
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
	RANK_NONE
)

const (
	Color_WHITE Color = iota
	Color_BLACK
	Color_BOTH
)

// Castling bits
const (
	WKCA int = 1
	WQCA int = 2
	BKCA int = 4
	BQCA int = 8
)

const (
	SQUARES_A1 Square = iota + 21
	SQUARES_B1
	SQUARES_C1
	SQUARES_D1
	SQUARES_E1
	SQUARES_F1
	SQUARES_G1
	SQUARES_H1
)

const (
	SQUARES_A8 Square = iota + 91
	SQUARES_B8
	SQUARES_C8
	SQUARES_D8
	SQUARES_E8
	SQUARES_F8
	SQUARES_G8
	SQUARES_H8
	SQUARES_NO_SQ    = 99
	SQUARES_OFFBOARD = 100
)

// Are the pieces big?
var PieceBig = []bool{
	false, // none
	false, // pawn
	true,  // knight
	true,  // bishop
	true,  // rook
	true,  // queen
	true,  // king
	false, // pawn
	true,  // knight
	true,  // bishop
	true,  // rook
	true,  // queen
	true,  // king
}

var PieceMaj = []bool{
	false, // none
	false, // pawn
	false, // knight
	false, // bishop
	true,  // rook
	true,  // queen
	true,  // king
	false, // pawn
	false, // knight
	false, // bishop
	true,  // rook
	true,  // queen
	true,  // king
}

var PieceMin = []bool{
	false, // none
	false, // pawn
	true,  // knight
	true,  // bishop
	false, // rook
	false, // queen
	false, // king
	false, // pawn
	true,  // knight
	true,  // bishop
	false, // rook
	false, // queen
	false, // king
}

var PieceVal = []int{
	0,     // none
	100,   // pawn
	325,   // knight
	325,   // bishop
	550,   // rook
	1000,  // queen
	50000, // king
	100,   // pawn
	325,   // knight
	325,   // bishop
	550,   // rook
	1000,  // queen
	50000, // king
}

var PieceCol = []int{
	Color_BOTH,  // none
	Color_WHITE, // pawn
	Color_WHITE, // knight
	Color_WHITE, // bishop
	Color_WHITE, // rook
	Color_WHITE, // queen
	Color_WHITE, // king
	Color_BLACK, // pawn
	Color_BLACK, // knight
	Color_BLACK, // bishop
	Color_BLACK, // rook
	Color_BLACK, // queen
	Color_BLACK, // king
}

var PiecePawn = []bool{
	false, // none
	true,  // pawn
	false, // knight
	false, // bishop
	false, // rook
	false, // queen
	false, // king
	true,  // pawn
	false, // knight
	false, // bishop
	false, // rook
	false, // queen
	false, // king
}

var PieceKnight = []bool{
	false, // none
	false, // pawn
	true,  // knight
	false, // bishop
	false, // rook
	false, // queen
	false, // king
	false, // pawn
	true,  // knight
	false, // bishop
	false, // rook
	false, // queen
	false, // king
}

var PieceKing = []bool{
	false, // none
	false, // pawn
	false, // knight
	false, // bishop
	false, // rook
	false, // queen
	true,  // king
	false, // pawn
	false, // knight
	false, // bishop
	false, // rook
	false, // queen
	true,  // king
}

var PieceRookQueen = []bool{
	false, // none
	false, // pawn
	false, // knight
	false, // bishop
	true,  // rook
	true,  // queen
	false, // king
	false, // pawn
	false, // knight
	false, // bishop
	true,  // rook
	true,  // queen
	false, // king
}

var PieceBishopQueen = []bool{
	false, // none
	false, // pawn
	false, // knight
	true,  // bishop
	false, // rook
	true,  // queen
	false, // king
	false, // pawn
	false, // knight
	true,  // bishop
	false, // rook
	true,  // queen
	false, // king
}

var PieceSlides = []bool{
	false, // none
	false, // pawn
	false, // knight
	true,  // bishop
	true,  // rook
	true,  // queen
	false, // king
	false, // pawn
	false, // knight
	true,  // bishop
	true,  // rook
	true,  // queen
	false, // king
}

/*
 This value bitwise anded to the castle permission setting.
 Basically will remove the bits to disable various castling.
 The rook and king positions line up with the bit removals
 to reflect the loss of castle permissions
 7 = 0111 So Black queenside is removed if the rook is to move
 3 = 0011 So Black queenside and kingside permissions are removed.
*/
var CastlePerm = []int{
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 13, 15, 15, 15, 12, 15, 15, 14, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 7, 15, 15, 15, 3, 15, 15, 11, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}

var Kings = []int{Piece_WKING, Piece_BKING}

// Piece movements
var KnightDirection = []int{-8, -19, -21, -12, 8, 19, 21, 12}
var RookDirection = []int{-1, -10, 1, 10}
var BishopDirection = []int{-11, -9, 9, 11}
var KingDirection = []int{-1, 1, -10, 10, -11, -9, 9, 11}
var NoMove []int

var PieceDirections = []([]int){NoMove,
	NoMove, KnightDirection, BishopDirection, RookDirection, KingDirection, KingDirection,
	NoMove, KnightDirection, BishopDirection, RookDirection, KingDirection, KingDirection}

var LoopNonSlidePce = []int{
	Piece_WKNIGHT, Piece_WKING, Piece_EMPTY,
	Piece_BKNIGHT, Piece_BKING, Piece_EMPTY}

var LoopSlidePce = []int{
	Piece_WROOK, Piece_WBISHOP, Piece_WQUEEN, Piece_EMPTY,
	Piece_BROOK, Piece_BBISHOP, Piece_BQUEEN, Piece_EMPTY}

var LoopNonSlideIndex = []int{0, 3}
var LoopSlideIndex = []int{0, 4}

// Hash keys
var PieceKeys = [14 * 120]int{0}

// One needed to represent the size to move. in and out as side changes
var SideKey = RAND_32()

// 16 different flag values
var CastleKeys = [16]int{0}

/**
Mapping that will tell us the file for each square in our 120 grid
*/
var filesBrd [] /*BRD_SQ_NUM*/ File
var ranksBrd [] /*BRD_SQ_NUM*/ Rank

// Mapping between 120 and 64 board versions
var Sq120To64 = [BRD_SQ_NUM]SQ64{0}
var Sq64To120 = [64]SQ120{0}

func FR2SQ(f File, r Rank) SQ120 {
	return ((21 + f) + (r * 10))
}

var randx *rand.Rand

func RAND_32() int {
	// r := randx.Int31()
	if randx == nil {
		randx = rand.New(rand.NewSource(time.Now().Unix()))
	}
	r := (randx.Intn(256) << 23) |
		(randx.Intn(256) << 16) |
		(randx.Intn(256) << 8) |
		(randx.Intn(256))

	// fmt.Println("R:", r)
	return int(r)
}

func initFileRanks() {
	// Create the files board
	filesBrd = make([]File, BRD_SQ_NUM)
	// Create the ranks board
	ranksBrd = make([]Rank, BRD_SQ_NUM)

	for index := 0; index < BRD_SQ_NUM; index++ {
		filesBrd[index] = SQUARES_OFFBOARD
		ranksBrd[index] = SQUARES_OFFBOARD
	}

	for ranks := 0; ranks <= RANK_8; ranks++ {
		for files := 0; files <= FILE_H; files++ {
			sq120 := FR2SQ(files, ranks)
			filesBrd[sq120] = files
			ranksBrd[sq120] = ranks
		}
	}

	// a := filesBrd
	// b := ranksBrd
	// println("Done with file ranks", len(a), len(b) )

}

/**
This is static, done once
*/
func init() {
	fmt.Println("Init called in definitions")
	NoMove = make([]int, 0)
	randx = rand.New(rand.NewSource(time.Now().Unix()))
	initFileRanks()
	initHashKeys()
	init120To64Mappings()
	initMVVLVV()
}

func initMVVLVV() {
	// Setup the attacker / victim array
	for attackers := Piece_WPAWN; attackers <= Piece_BKING; attackers++ {
		for victim := Piece_WPAWN; victim <= Piece_BKING; victim++ {
			score := mvvLvaValue[victim] + 6 - (mvvLvaValue[attackers] / 100)
			mvvLvaScores[victim*14+attackers] = score
			//fmt.Println("MVV: ", string(PceCharLetter[attackers]),
			//	" x ", string(PceCharLetter[victim]), " = ", score)
		}
	}
}

/**
 * Convert from a file and rank to the square index
 */
func fr_to_SQ120(file File, rank Rank) SQ120 {
	return (21 + (file)) + ((rank) * 10)
}

func sq_to_RANK(sq SQ120) int {
	return ranksBrd[sq]
}

func sq_to_FILE(sq SQ120) int {
	return filesBrd[sq]
}

func sq_64(sq120 SQ120) int {
	return Sq120To64[sq120]
}

func sq_120(sq64 SQ64) SQ120 {
	return Sq64To120[sq64]
}

func init120To64Mappings() {
	var file = FILE_A
	var rank = RANK_1
	var sq = SQUARES_A1
	var sq64 = 0

	// RESET
	for index := 0; index < BRD_SQ_NUM; index++ {
		// println("index = ${index}")
		Sq120To64[index] = 65
	}

	// RESET
	for index := 0; index < 64; index++ {
		Sq64To120[index] = 120
	}

	// Loop in rank and file order.
	for rank = RANK_1; rank <= RANK_8; rank++ {
		for file = FILE_A; file <= FILE_H; file++ {
			sq = fr_to_SQ120(file, rank)
			// println("Setting ${rank}, ${file} = ${sq} <-> ${sq64}")
			Sq64To120[sq64] = sq
			Sq120To64[sq] = sq64
			sq64++
		}
	}
}

func initHashKeys() {

	for x := 0; x < len(PieceKeys); x++ {
		PieceKeys[x] = RAND_32()
	}

	SideKey = RAND_32()

	for x := 0; x < len(CastleKeys); x++ {
		CastleKeys[x] = RAND_32()
	}

}

/*
Convert a sq120 index into the a1 .. h8 string
*/
func sqToString(sq120 SQ120) string {

	if sq120 == SQUARES_NO_SQ {
		return "  "
	}

	file := filesBrd[sq120]
	rank := ranksBrd[sq120]
	// fmt.Println("F:" , file, " R:", rank, " SQ120:", sq120)
	// return FileChar[file:file+1] + RankChar[rank:rank+1]
	return string(FileChar[file]) + string(RankChar[rank])
}
