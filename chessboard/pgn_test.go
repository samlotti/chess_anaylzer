package ai

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPgnLoad1(t *testing.T) {

	var pgn = `[Event "F/S Return Match"] 
[Site "Belgrade, Serbia JUG"] 
[Date "1992.11.04"] 
[Round "29"] 
[White "Fischer, Robert J."]
[Black "Spassky, Boris V."] 
[Result "1/2-1/2"] 

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3
O-O 9. h3 Nb8 10. d4 Nbd7 11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15.
Nb1 h6 16. Bh4 c5 17. dxe5 Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21.
Nc4 Nxc4 22. Bxc4 Nb6 23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7
27. Qe3 Qg5 28. Qxg5 hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33.
f3 Bc8 34. Kf2 Bf5 35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5
40. Rd6 Kc5 41. Ra6 Nf2 42. g4 Bd3 43. Re6 1/2-1/2
`
	pw := NewPgnWrapper(pgn)
	err := pw.Parse()
	assert.Nil(t, err)

	assert.Equal(t, 7, len(pw.Attributes))

	assert.True(t, pw.isEof())

	//moves, err := LoadPgn(pgn)
	//assert.Nil(t, err)
	//assert.Equal(t, 42*2+1, len(moves))
}

func TestPgnWithManyTypesOfData(t *testing.T) {

	var pgn = `[[Event "Vilnius All-Russian Masters"]
[Site "Vilna (Vilnius) RUE"]
[Date "1912.08.23"]
[EventDate "1912.08.19"]
[Round "5"]
[Result "0-1"]
[White "Alexander Alekhine"]
[Black "Akiba Rubinstein"]
[ECO "C83"]
[WhiteElo "?"]
[BlackElo "?"]
[PlyCount "54"]

1. e4 {Notes by Dr. Savielly Tartakower.} 1... e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4
Nf6 5. O-O Nxe4 6. d4 b5 7. Bb3 d5 8. dxe5 Be6 9. c3 Be7 10. Nbd2 Nc5 11. Bc2
Bg4 12. h3 {The most reasonable course here is 12.Re1, guarding the e-pawn.}
12... Bh5 13. Qe1 $6 {Here again 13. Re1 ensured a very good game for White.}
13... Ne6 14. Nh2 $6 Bg6 $1 15. Bxg6 fxg6 {! Far seeing strategy! Black
recognizes that the f-file and not the e-file will be needed as a base for
action.} 16. Nb3 {Or 16.f4 d4!.} 16... g5 $1 17. Be3 O-O 18. Nf3 Qd7 19. Qd2
{White pays insufficient attention to the scope of his opponent's threats. A
better course is 19.Nfd4 (19...Nxe5 20.Bxg5) seeking to establish equality.}
19... Rxf3 $1 20. gxf3 Nxe5 21. Qe2 Rf8 22. Nd2 Ng6 23. Rfe1 Bd6 24. f4 Nexf4
25. Qf1 Nxh3+ 26. Kh1 g4 27. Qe2 Qf5 0-1
`
	pw := NewPgnWrapper(pgn)
	err := pw.Parse()
	assert.Nil(t, err)

	assert.Equal(t, 7, len(pw.Attributes))

	assert.True(t, pw.isEof())

	//moves, err := LoadPgn(pgn)
	//assert.Nil(t, err)
	//assert.Equal(t, 42*2+1, len(moves))
}

func TestRegex(t *testing.T) {

	rTag, _ := regexp.Compile(tag)
	r := rTag.Find([]byte("[Fen rfwefwef] test"))
	fmt.Println(r)
	assert.True(t, len(r) > 0)

	r = rTag.Find([]byte("1. [Fen rfwefwef] test"))
	fmt.Println(r)
	assert.True(t, len(r) == 0)

	match, _ := regexp.MatchString("p([a-z]+)ch", "peach cobbler")
	fmt.Println(match)
}
