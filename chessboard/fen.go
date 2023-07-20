package ai

import (
	"bytes"
	"strconv"
)

func BoardToFen(board *Board, numMoves int) string {

	var buffer bytes.Buffer

	var rank = RANK_8
	var file = FILE_A

	var first = true

	for rank = RANK_8; rank >= RANK_1; rank-- {
		if !first {
			buffer.WriteString("/")
		}
		first = false
		var blanks = 0
		for file = FILE_A; file <= FILE_H; file++ {
			sq120 := FR2SQ(file, rank)
			piece := board.pieces[sq120]
			if piece == Piece_EMPTY {
				blanks++
				continue
			} else {
				if blanks > 0 {
					buffer.WriteString(strconv.Itoa(blanks))
				}
				// log.Print("Piece: ", piece)
				buffer.WriteString(string(PceCharLetter[piece]))
				blanks = 0
			}
		}
		if blanks > 0 {
			buffer.WriteString(strconv.Itoa(blanks))
		}
	}

	buffer.WriteString(" ")
	buffer.WriteString(string(SideChar[board.side]))

	buffer.WriteString(" ")

	cdash := true
	if board.castlePermFlag&WKCA != 0 {
		buffer.WriteString("K")
		cdash = false
	}
	if board.castlePermFlag&WQCA != 0 {
		buffer.WriteString("Q")
		cdash = false
	}
	if board.castlePermFlag&BKCA != 0 {
		buffer.WriteString("k")
		cdash = false
	}
	if board.castlePermFlag&BQCA != 0 {
		buffer.WriteString("q")
		cdash = false
	}
	if cdash {
		buffer.WriteString("-")
	}

	buffer.WriteString(" ")
	if board.enPassantSq120 != SQUARES_NO_SQ {
		buffer.WriteString(sqToString(board.enPassantSq120))
	} else {
		buffer.WriteString("-")
	}

	// moves since last capture
	buffer.WriteString(" ")
	buffer.WriteString(strconv.Itoa(board.fiftyMove))

	// full moves
	buffer.WriteString(" ")
	buffer.WriteString(strconv.Itoa(numMoves/2 + 1))

	return buffer.String()

}

// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
func ParseFen(board *Board, fen string) {

	board.ResetBoard()

	var rank = RANK_8
	var file = FILE_A
	var piece = 0
	var count = 0
	var i = 0
	var sq120 = 0
	var fenCnt = 0

	//fmt.Println("8 = ", '8')
	//fmt.Println("1 = ", '1')

	for rank = RANK_8; rank >= RANK_1 && fenCnt < len(fen); {
		count = 1

		//fmt.Println("Process:" , fen[fenCnt:fenCnt+1], "  fen pos: " , fenCnt,  )
		switch fen[fenCnt] {
		case 'p':
			piece = Piece_BPAWN
		case 'r':
			piece = Piece_BROOK
		case 'n':
			piece = Piece_BKNIGHT
		case 'b':
			piece = Piece_BBISHOP
		case 'k':
			piece = Piece_BKING
		case 'q':
			piece = Piece_BQUEEN

		case 'P':
			piece = Piece_WPAWN
		case 'R':
			piece = Piece_WROOK
		case 'N':
			piece = Piece_WKNIGHT
		case 'B':
			piece = Piece_WBISHOP
		case 'K':
			piece = Piece_WKING
		case 'Q':
			piece = Piece_WQUEEN

		case '1', '2', '3', '4', '5', '6', '7', '8':
			//fmt.Println("Here")
			piece = Piece_EMPTY
			count = int(fen[fenCnt] - '0')
			//fmt.Println("for ", fen[fenCnt], " = " , count)

		case '/', ' ':
			//fmt.Println("process /")
			rank--
			file = FILE_A
			fenCnt++
			count = 0
			continue
		default:
			//fmt.Println("Fen error", fen)
			return
		}

		//fmt.Println("---- " , count)
		for i = 0; i < count; i++ {
			sq120 = FR2SQ(file, rank)
			//fmt.Println("Set piece: ", piece, " into", sq120, "File:", file, " Rank:", rank)
			board.pieces[sq120] = piece
			file++
		}
		fenCnt++
	}

	if fenCnt < len(fen) {
		// w or b
		switch fen[fenCnt] {
		case 'w':
			board.side = Color_WHITE
		case 'b':
			board.side = Color_BLACK
		}
	}

	fenCnt += 2

	var brk = false
	for i = 0; (i < 4) && !brk; i++ {
		if fenCnt >= len(fen) {
			break
		}
		// println("Castle fen ${fenChar[fenCnt]}")
		switch fen[fenCnt] {
		case 'K':
			board.castlePermFlag |= WKCA
		case 'Q':
			board.castlePermFlag |= WQCA
		case 'k':
			board.castlePermFlag |= BKCA
		case 'q':
			board.castlePermFlag |= BQCA
		default:
			brk = true
			// fenCnt += 1
			// break
		}

		if !brk {
			fenCnt += 1
		}

	}

	for fenCnt < len(fen) {
		if fen[fenCnt] == ' ' {
			fenCnt++
			break
		}
		fenCnt++
	}

	// fenCnt++

	if fenCnt < len(fen) {
		// En Passant
		//fmt.Println("At " , fenCnt, "    ", fen[fenCnt], "   ", fen)
		//fmt.Println("At " , fenCnt, "    ", fen[fenCnt], "   ", fen[fenCnt:])
		if fen[fenCnt] != '-' {
			file = int(fen[fenCnt] - 'a')
			fenCnt++
			rank = int(fen[fenCnt] - '1')
			board.enPassantSq120 = fr_to_SQ120(file, rank)
		}
		fenCnt += 2

	}

	board.posKey = board.generatePosKey()

	board.updateListsMaterial()

}
