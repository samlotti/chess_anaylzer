package ai

import (
	"fmt"
	"github.com/samlotti/chess_anaylzer/chessboard/minilex"
	"strings"
)

const (
	tag        = "\\[.*?\\]"       // Capture tags as a unit
	comment    = "\\{(.|\\n)*?\\}" // Capture comments as a unit
	resumption = "\\d+\\.\\.\\."   // Resume moves after comment
	moveNumber = "\\d+\\."
	endOfGame  = "0-1|1-0|0-0|1/2-1/2"
	nag        = "\\$\\d+"           //  Numeric annotation glyph
	move       = "[-+\\w\\.(=.)?/]+" // # Anything else is a move
	newline    = "\n"
	//whitespace = "\\s+"
)

const (
	TAG minilex.TokenId = minilex.UserTokeId + iota
	COMMENT
	RESUMPTION
	MOVENUMBER
	ENDOFGAME
	NAG
	MOVE
	NEWLINE
	//WHITESPACE
)

func AddPgnLexMap(lexer *minilex.MiniLexer) error {
	err := lexer.AddPattern(TAG, tag)
	if err == nil {
		err = lexer.AddPattern(COMMENT, comment)
	}
	if err == nil {
		err = lexer.AddPattern(COMMENT, comment)
	}

	if err == nil {
		err = lexer.AddPattern(RESUMPTION, resumption)
	}

	if err == nil {
		err = lexer.AddPattern(MOVENUMBER, moveNumber)
	}

	if err == nil {
		err = lexer.AddPattern(ENDOFGAME, endOfGame)
	}

	if err == nil {
		err = lexer.AddPattern(NAG, nag)
	}

	if err == nil {
		err = lexer.AddPattern(MOVE, move)
	}

	if err == nil {
		err = lexer.AddPattern(NEWLINE, newline)
	}
	//
	//if err == nil {
	//	err = lexer.AddPattern(WHITESPACE, whitespace)
	//}
	return err
}

type PgnWrapper struct {
	lex           *minilex.MiniLexer
	Moves         []string
	InternalMoves []Move
	Attributes    map[string]string
	board         *Board
	StartFen      string
}

func NewPgnWrapper(pgn string) *PgnWrapper {

	p := &PgnWrapper{}

	mo := minilex.NewMiniLexOptions()
	mo.RemoveAsWhiteSpace("\n")
	p.lex = minilex.NewMiniLexer(pgn, mo)
	AddPgnLexMap(p.lex)

	p.resetForNextPgn()

	return p
}

func (p *PgnWrapper) Parse() error {

	p.resetForNextPgn()

	for {
		tk, err := p.lex.PeekToken()
		if err != nil {
			return err
		}
		if tk.Id != NEWLINE {
			break
		}

		// Advance it
		p.lex.NextToken()
	}

	err := p.loadAttributes()
	if err != nil {
		return err
	}

	err = p.loadMoves()

	return err
}

// loadAttributes - reads the attributed until blank line
// [Event "F/S Return Match"]
// Will be after the NewLine, ready for moves
func (p *PgnWrapper) loadAttributes() error {
	for {
		tk, err := p.lex.NextToken()
		if err != nil {
			return err
		}

		if tk.Is(NEWLINE) {
			break
		}

		err = tk.AssertIs(TAG)
		if err != nil {
			return err
		}

		tagKV := tk.Literal

		// Advance Tag
		tk, err = p.lex.NextToken()
		if err != nil {
			return err
		}

		err = tk.AssertIs(NEWLINE)
		if err != nil {
			return err
		}

		tagKV = strings.TrimSuffix(tagKV, "]")
		tagKV = strings.TrimPrefix(tagKV, "[")
		sects := strings.SplitN(tagKV, " ", 2)
		key := sects[0]
		data := strings.Trim(sects[1], "\"")

		fmt.Printf("Tag: %s\n", tagKV)

		fmt.Println("Found:", key, "  ", data)
		p.Attributes[key] = data

		//if key == "FEN" {
		//	p.StartFen = data
		//}

	}
	return nil
}

// UP until the end or a newline, read all data.
func (p *PgnWrapper) loadMoves() error {

	if len(p.StartFen) > 0 {
		ParseFen(p.board, p.StartFen)
	} else {
		ParseFen(p.board, StartFen)
	}

	for {
		if p.lex.IsEOF() {
			return nil
		}

		tk, err := p.lex.PeekToken()
		if err != nil {
			return err
		}

		if tk.Is(minilex.TKEof) {
			return nil
		}

		tk, err = p.lex.NextToken()
		if err != nil {
			return err
		}
		if tk.Is(NEWLINE) {
			// blank line is end of this game
			if tk.Pos == 0 {
				break
			}
			continue
		}
		if tk.Is(NAG) {
			continue
		}
		if tk.Is(COMMENT) {
			continue
		}
		if tk.Is(MOVENUMBER) {
			continue
		}
		if tk.Is(ENDOFGAME) {
			continue
		}
		if tk.Is(RESUMPTION) {
			// ?
			continue
		}
		err = tk.AssertIs(MOVE)
		if err != nil {
			return err
		}

		err = p.applyMoveSAN(tk.Literal, false)
		if err != nil {
			return err
		}

	}
	return nil
}

func (p *PgnWrapper) applyMoveSAN(sanMove string, debug bool) error {
	// fmt.Printf("Move: %s\n", sanMove)

	// Get rid of check indicator
	sanMove = strings.TrimSuffix(sanMove, "+")

	// Check for available moves.
	// vmoves := GetAllValidMoves(p.board)
	vmoves := GetAllMoves(p.board)
	for _, vm := range vmoves {
		san := PgnForMove(p.board, vm)
		if debug {
			fmt.Printf(" possible move: %s for %s\n", san, sanMove)
		}
		if san == sanMove {
			p.board.MakeMove(vm, sanMove)

			p.InternalMoves = append(p.InternalMoves, vm)
			p.Moves = append(p.Moves, MoveToString(vm))

			return nil
		}
	}

	p.board.PrintBoard(fmt.Sprintf("Move not found for: %s", sanMove))
	//if !debug {
	//	// So we can see it.
	//	p.applyMoveSAN(sanMove, true)
	//}

	return fmt.Errorf(fmt.Sprintf("Move not found for: %s", sanMove))
}

func (p *PgnWrapper) IsEof() bool {
	return p.lex.IsEOF()
}

func (p *PgnWrapper) resetForNextPgn() {
	p.Moves = make([]string, 0)
	p.Attributes = make(map[string]string)
	p.board = NewBoard()

}

// PgnForMove
// Returns the pgn format for the move
// Note: This must be before the move is made so we can see the possible moves
func PgnForMove(b *Board, m Move) string {

	// moves := GetAllValidMoves(b)
	moves := GetAllMoves(b)

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
	fromRank := ranksBrd[fSq120]

	// toFile := filesBrd[getToSq120(m)]
	tSq120 := getToSq120(m)
	// toRank := ranksBrd[tSq120]

	sep := ""
	if capturedPiece != Piece_EMPTY {
		sep = "x"
	}

	// Pawn can enpassant so.
	if fromPiece == Piece_WPAWN || fromPiece == Piece_BPAWN {
		if fromFile != filesBrd[tSq120] {
			sep = "x"
		}
	}

	trail := ""
	//if ( m.isCheck() ) trail = "+";
	//if ( m.isMate() ) trail = "#";

	promote := ""
	if promotedPiece != Piece_EMPTY {
		promote = strings.ToUpper(fmt.Sprintf("=%s", string(PceCharLetter[promotedPiece])))
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
		if sep != "" {
			// if capturedPiece != Piece_EMPTY {
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
		rankString = fmt.Sprintf("%s", string(RankChar[fromRank]))
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
