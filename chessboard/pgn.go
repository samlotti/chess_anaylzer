package ai

import (
	"fmt"
	"strings"
)

const (
	tag        = "\\[.*?\\]"     // Capture tags as a unit
	comment    = "\\{.*?\\}"     // Capture comments as a unit
	resumption = "\\d+\\.\\.\\." // Resume moves after comment
	moveNumber = "\\d+\\."
	endOfGame  = "0-1|1-0|0-0|1/2-1/2"
	nag        = "\\$\\d+"      //  Numeric annotation glyph
	move       = "[-+\\w\\./]+" // # Anything else is a move
	newline    = "\n"
	whitespace = "\\s+"
)

func GetPgnLexMap() map[string]string {
	m := make(map[string]string)
	m["tag"] = tag
	m["comment"] = comment
	m["resumption"] = resumption
	m["moveNumber"] = moveNumber
	m["endOfGame"] = endOfGame
	m["nag"] = nag
	m["move"] = move
	m["newline"] = newline
	m["whitespace"] = whitespace
	return m
}

type PgnWrapper struct {
	Pgn        string
	pos        int
	line       int
	Moves      []Move
	Attributes map[string]string
	board      *Board
	StartFen   string
}

func NewPgnWrapper(pgn string) *PgnWrapper {
	p := &PgnWrapper{Pgn: pgn}
	p.Moves = make([]Move, 0)
	p.Attributes = make(map[string]string)
	p.pos = 0
	p.line = 1
	p.board = NewBoard()
	return p
}

func (p *PgnWrapper) Parse() error {
	// Remove leading NL
	for {
		if p.isNewLine() {
			p.advance()
		} else {
			break
		}
	}

	err := p.loadAttributes()
	if err != nil {
		return err
	}

	p.advanceNL()
	err = p.loadMoves()

	return err
}

// loadAttributes - reads the attributed until blank line
// [Event "F/S Return Match"]
func (p *PgnWrapper) loadAttributes() error {
	for {
		// at end of attributes
		if p.isNewLine() {
			break
		}

		if !p.isChar('[') {
			return fmt.Errorf(fmt.Sprintf("Expected '[' at pos: %d", p.pos))
		}

		p.advance()
		key, err := p.readTil(' ')
		if err != nil {
			return err
		}
		p.advanceSpaces()

		if !p.isChar('"') {
			return fmt.Errorf(fmt.Sprintf("Expected '\"' at pos: %d", p.pos))
		}
		p.advance()
		data, err := p.readTil('"')
		if err != nil {
			return err
		}
		p.advance()
		p.advanceSpaces()
		if !p.isChar(']') {
			return fmt.Errorf(fmt.Sprintf("Expected ']' at pos: %d", p.pos))
		}
		p.advance()
		p.advanceSpaces()
		p.advanceNL()

		fmt.Println("Found:", key, "  ", data)
		p.Attributes[key] = data

		if key == "FEN" {
			p.StartFen = data
		}

	}
	return nil
}

func (p *PgnWrapper) advance() {
	p.pos += 1

	if p.isEof() {
		return
	}

	if p.Pgn[p.pos] == '\n' {
		p.line++
	}
}

func (p *PgnWrapper) readTil(chr uint8) (string, error) {
	st := p.pos
	for {
		if p.Pgn[p.pos] == chr {
			break
		}
		p.advance()

		if chr != '\n' {
			if p.isNewLine() || p.isEof() {
				return "", fmt.Errorf("invalid pgn at line %d, unexpected new line", p.line)
			}
		}

		if p.isEof() {
			return "", fmt.Errorf("invalid pgn at line %d, unexpected end of text", p.line)
		}
	}
	str := p.Pgn[st:p.pos]
	return str, nil
}

func (p *PgnWrapper) isNewLine() bool {
	return p.Pgn[p.pos] == '\n'
}

func (p *PgnWrapper) isEof() bool {
	return p.pos >= len(p.Pgn)
}

func (p *PgnWrapper) advanceSpaces() {
	for {
		if p.isSpace() {
			p.advance()
		} else {
			break
		}
	}
}

func (p *PgnWrapper) isSpace() bool {
	return p.Pgn[p.pos] == ' '
}

func (p *PgnWrapper) isChar(chr uint8) bool {
	return p.Pgn[p.pos] == chr
}

func (p *PgnWrapper) advanceNL() {
	p.advanceSpaces()
	if p.isNewLine() {
		p.advance()
	}
}

// UP until the end or a newline, read all data.
func (p *PgnWrapper) loadMoves() error {

	if len(p.StartFen) > 0 {
		ParseFen(p.board, p.StartFen)
	} else {
		ParseFen(p.board, StartFen)
	}

	for {
		if p.isEof() {
			return nil
		}

		if p.isNewLine() {
			return nil
		}
		line, err := p.readTil('\n')
		if err != nil {
			return err
		}
		p.advanceNL()
		fmt.Println("Line: ", line)

		entries := strings.Split(line, " ")
		for _, s := range entries {
			s = strings.TrimSpace(s)
			if strings.HasSuffix(s, ".") {
				fmt.Printf("move# %s\n", s)
			} else {
				if strings.Contains(s, "-") {
					// if s == "1" || s == "-1" || s == "1/2" {
					fmt.Printf("end   %s\n", s)
				} else {
					fmt.Printf("move   %s\n", s)

					err := p.applyMoveSAN(s)
					if err != nil {
						return err
					}

				}
			}
		}
	}
	return nil
}

func (p *PgnWrapper) applyMoveSAN(sanMove string) error {

	// Check for available moves.
	vmoves := GetAllValidMoves(p.board)
	for _, vm := range vmoves {
		san := PgnForMove(p.board, vm)
		if san == sanMove {
			p.board.makeMove(vm)
			return nil
		}
	}

	p.board.printBoard(fmt.Sprintf("Move not found for: %s", sanMove))

	return fmt.Errorf(fmt.Sprintf("Move not found for: %s", sanMove))
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
