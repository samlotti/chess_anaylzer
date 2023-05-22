package ai

import (
	"fmt"
	"strings"
)

// LoadPgn -- returns a list of moves for the pgn
// or error if could not load the pgn
func LoadPgn(pgn string) ([]Move, error) {
	return nil, nil
}

// PgnForMove
// Returns the pgn format for the move
// Note: This must be before the move is made so we can see the possible moves
func PgnForMove(b *Board, m Move) string {

	moves := GetAllValidMoves(b)

	isAmbigRank := false
	isAmbigFile := false

	for _, availMove := range moves {
		if b.pieces[getFromSq120(m)] == b.pieces[getFromSq120(availMove)] {
			if getFromSq120(m) != getFromSq120(availMove) &&
				getToSq120(m) == getToSq120(availMove) {
				// Same piece type but not the same starting square

				// detection of rank / file
				if filesBrd[getFromSq120(availMove)] != filesBrd[getFromSq120(m)] {
					isAmbigFile = true
				} else {
					if ranksBrd[getFromSq120(availMove)] != ranksBrd[getFromSq120(m)] {
						isAmbigRank = true
					}
				}
				break
			}
		}
	}

	return asMoveString(b, m, isAmbigRank, isAmbigFile)

}

func asMoveString(b *Board, m Move, isAmbigRank bool, isAmbigFile bool) string {

	//log.Println( MoveToString(m) )
	//fmt.Println("Fen:", BoardToFen(b,0))

	capturedPiece := getCapturedPiece(m)
	promotedPiece := getPromoted(m)

	// This must be before move is make
	fSq120 := getFromSq120(m)
	fromPiece := b.pieces[fSq120]
	fromFile := filesBrd[fSq120]
	// fromRank := ranksBrd[getFromSq120(m)]

	// toFile := filesBrd[getToSq120(m)]
	tSq120 := getToSq120(m)
	// toRank := ranksBrd[tSq120]

	sep := ""
	if capturedPiece != Piece_EMPTY {
		sep = "x"
	}

	trail := ""
	//if ( m.isCheck() ) trail = "+";
	//if ( m.isMate() ) trail = "#";

	promote := ""
	if promotedPiece != Piece_EMPTY {
		promote = strings.ToUpper(fmt.Sprintf("= %s", string(PceCharLetter[promotedPiece])))
	}

	if fromPiece == Piece_BKING {
		if getFromSq120(m) == SQUARES_E8 && getToSq120(m) == SQUARES_G8 {
			return "O-O"
		}
		if getFromSq120(m) == SQUARES_E8 && getToSq120(m) == SQUARES_C8 {
			return "O-O-O"
		}
	}

	if fromPiece == Piece_WKING {
		if getFromSq120(m) == SQUARES_E1 && getToSq120(m) == SQUARES_C1 {
			return "O-O-O"
		}
		if getFromSq120(m) == SQUARES_E1 && getToSq120(m) == SQUARES_G1 {
			return "O-O"
		}
	}

	pieceIdent := ""
	switch fromPiece {
	case Piece_WBISHOP:
		pieceIdent = "B"
		break
	case Piece_WKING:
		pieceIdent = "K"
		break
	case Piece_WKNIGHT:
		pieceIdent = "N"
		break
	case Piece_WROOK:
		pieceIdent = "R"
		break
	case Piece_WQUEEN:
		pieceIdent = "Q"
		break
	case Piece_BBISHOP:
		pieceIdent = "B"
		break
	case Piece_BKING:
		pieceIdent = "K"
		break
	case Piece_BKNIGHT:
		pieceIdent = "N"
		break
	case Piece_BROOK:
		pieceIdent = "R"
		break
	case Piece_BQUEEN:
		pieceIdent = "Q"
		break
	}

	// For winboard... if Pawn move and Capture... and not ambig, Add the 'P'
	if !isAmbigFile {
		if capturedPiece != Piece_EMPTY {
			if len(pieceIdent) == 0 {
				// pieceIdent = "P";
				// Per spec, should be the file
				pieceIdent = fmt.Sprintf("%s", string(FileChar[fromFile]))
			}
		}
	}

	/*
	 * Some viewers had issues even when not ambig.
	 */
	rankString := ""
	fileString := ""
	if isAmbigFile {
		fileString = fmt.Sprintf("%s", string(FileChar[fromFile]))
	}
	if isAmbigRank {
		rankString = fmt.Sprintf("%s", string(RankChar[fromFile]))
	}

	//log.Println(
	//	"Pgn:", "piece:",	pieceIdent,
	//	" file:", fileString,
	//	" rank:", rankString,
	//	" sep:", sep,
	//	" rankChar:", sqToString(tSq120),
	//	" promote:", promote,
	//	" trail:", trail)

	r := fmt.Sprintf(
		"%s%s%s%s%s%s%s",
		pieceIdent,
		fileString,
		rankString, sep,
		sqToString(tSq120), promote, trail)

	// log.Println("> ", r)
	return r

	/*return "" + m.from.toString().toLowerCase() + sep +
	m.to.toString().toLowerCase() + trail + promote;
	*/
}
