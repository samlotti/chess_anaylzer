package minilex

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
	//whitespace = "\\s+"
)

const (
	TAG TokenId = UserTokeId + iota
	COMMENT
	RESUMPTION
	MOVENUMBER
	ENDOFGAME
	NAG
	MOVE
	NEWLINE
	//WHITESPACE
)

func AddPgnLexMap(lexer *MiniLexer) error {
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

func TestMinLex1(t *testing.T) {

	var pgn = `[Event "F/S Return Match"] 
[Site "Belgrade, Serbia JUG"] 

1. e4 e5 
2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 
9. h3 Nb8 10. d4 Nbd7 11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 
15. Nb1 h6 16. Bh4 c5 17. dxe5 Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21.
Nc4 Nxc4 22. Bxc4 Nb6 23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7
27. Qe3 Qg5 28. Qxg5 hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33.
f3 Bc8 34. Kf2 Bf5 35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5
40. Rd6 Kc5 41. Ra6 Nf2 42. g4 Bd3 43. Re6 1/2-1/2
`

	mo := NewMiniLexOptions()
	mo.RemoveAsWhiteSpace("\n")
	ml := NewMiniLexer(pgn, mo)
	AddPgnLexMap(ml)
	tk, err := ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, TAG, tk.id)
	assert.Equal(t, "[Event \"F/S Return Match\"]", tk.literal)

	tk, err = ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, NEWLINE, tk.id)
	assert.Equal(t, "\n", tk.literal)
	assert.True(t, tk.pos > 5)
	assert.True(t, tk.line == 1)

	tk, err = ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, TAG, tk.id)
	assert.Equal(t, "[Site \"Belgrade, Serbia JUG\"]", tk.literal)

	tk, err = ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, NEWLINE, tk.id)
	assert.Equal(t, "\n", tk.literal)
	assert.True(t, tk.pos > 5)
	assert.True(t, tk.line == 2)

	tk, err = ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, NEWLINE, tk.id)
	assert.Equal(t, "\n", tk.literal)
	assert.True(t, tk.pos == 0)
	assert.True(t, tk.line == 3)

	tk, err = ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, MOVENUMBER, tk.id)
	assert.Equal(t, "1.", tk.literal)
	assert.True(t, tk.pos == 0)
	assert.True(t, tk.line == 4)

	tk, err = ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, MOVE, tk.id)
	assert.Equal(t, "e4", tk.literal)
	assert.True(t, tk.pos == 3)
	assert.True(t, tk.line == 4)

	tk, err = ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, MOVE, tk.id)
	assert.Equal(t, "e5", tk.literal)
	assert.True(t, tk.pos == 6)
	assert.True(t, tk.line == 4)

	tk, err = ml.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, NEWLINE, tk.id)
	assert.Equal(t, "\n", tk.literal)
	assert.True(t, tk.pos > 6)
	assert.True(t, tk.line == 4)

}

func TestMinLex2(t *testing.T) {

	var pgn = `
[Event "F/S Return Match"] 
[Site "Belgrade, Serbia JUG"] 

1. e4 e5 
2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 
9. h3 Nb8 10. d4 Nbd7 11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 
15. Nb1 h6 16. Bh4 c5 17. dxe5 Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21.
Nc4 Nxc4 22. Bxc4 Nb6 23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7
27. Qe3 Qg5 28. Qxg5 hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33.
f3 Bc8 34. Kf2 Bf5 35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5
40. Rd6 Kc5 41. Ra6 Nf2 42. g4 Bd3 43. Re6 1/2-1/2
`

	mo := NewMiniLexOptions()
	mo.RemoveAsWhiteSpace("\n")
	ml := NewMiniLexer(pgn, mo)
	AddPgnLexMap(ml)
	tkList, err := ml.ReadAllTokens()
	assert.Nil(t, err)

	for _, tk := range tkList {
		fmt.Println(tk)
	}

	assert.Equal(t, 142, len(tkList))

}
