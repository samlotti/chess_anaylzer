package analyzer

import (
	"fmt"
	"github.com/samlotti/chess_anaylzer/chessboard/common"
	"github.com/samlotti/chess_anaylzer/uci"
	"testing"
)
import "github.com/stretchr/testify/assert"

func collectPgnResults(in chan *PgnResponse) []*PgnResponse {

	r := make([]*PgnResponse, 0)

	for {
		m := <-in
		fmt.Printf("AnswerPgn: %+v \n", m)
		r = append(r, m)
		if m.Done {
			break
		}
	}
	return r
}

func TestPgn1(t *testing.T) {

	a := &PgnData{}
	assert.NotNil(t, a)
	a.Pgn = `
[Event "F/S Return Match"]
[Site "Belgrade, Serbia JUG"]
[Date "1992.11.04"]
[Round "29"]
[White "Fischer, Robert J."]
[Black "Spassky, Boris V."]
[Result "1/2-1/2"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 {This opening is called the Ruy Lopez.}
4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7
11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15. Nb1 h6 16. Bh4 c5 17. dxe5
Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21. Nc4 Nxc4 22. Bxc4 Nb6
23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7 27. Qe3 Qg5 28. Qxg5
hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33. f3 Bc8 34. Kf2 Bf5
35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5 40. Rd6 Kc5 41. Ra6
Nf2 42. g4 Bd3 43. Re6 1/2-1/2
`
	a.NumLines = 5
	rchan := make(chan *PgnResponse, 10)
	a.RChannel = rchan
	pa := NewPgnAnalyzer()
	go pa.DoAnalyze(a)
	ar := collectPgnResults(rchan)
	assert.NotNil(t, ar)
}

// Takes 1min 10sec
func TestPgn2(t *testing.T) {

	Verbose = false
	common.Verbose = false
	uci.Verbose = false

	a := &PgnData{}
	assert.NotNil(t, a)
	a.Pgn =
		`
[Event "Superbet Rapid 2023"]
[Site "Warsaw POL"]
[Date "2023.05.21"]
[Round "3"]
[White "So,W"]
[Black "Rapport,R"]
[Result "1-0"]
[WhiteElo "2760"]
[BlackElo "2745"]
[EventDate "2023.05.19"]
[ECO "A61"]

1. d4 Nf6 2. c4 e6 3. Nf3 c5 4. d5 d6 5. Nc3 exd5 6. cxd5 g6 7. Bf4 Bg7 8.
e3 O-O 9. h3 Ne8 10. Be2 Nd7 11. O-O Ne5 12. Nd2 f5 13. Bh2 Nc7 14. a4 Qe7
15. Re1 Bd7 16. Qb3 Na6 17. Qxb7 Nb4 18. Bxe5 Bxe5 19. Bb5 Rfd8 20. Red1
Rab8 21. Qc7 Rbc8 22. Qb7 Rb8 23. Qc7 Rbc8 24. Qa5 a6 25. Bxd7 Rxd7 26. Nc4
Rb8 27. Na2 Ra7 28. Nxb4 Rxb4 29. Nxe5 Qxe5 30. Qd8+ Kg7 31. Rd2 Rab7 32.
Qa8 Rxb2 33. Rxb2 Qxb2 34. Rf1 c4 35. Qxa6 c3 36. Kh2 c2 37. Qc4 Rb4 38.
Qc7+ Kh6 39. Qxd6 Rc4 40. Qf8+ Kh5 41. Qe7 h6 42. d6 c1=Q 43. Rxc1 Rxc1 44.
g4+ fxg4 45. hxg4+ Kxg4 46. Qe4+ Kg5 47. Qf4+ Kh5 48. d7 Qd2 49. Qf3+ Kg5
50. Qf4+ Kh5 51. Qe5+ g5 52. Qe8+ Kh4 53. Qe4+ Kh5 54. Qe8+ Kh4 55. Qf8
Rh1+ 56. Kxh1 Qxd7 57. Qxh6+ Kg4 58. Qa6 Kf3 59. Qf6+ Ke2 60. Kh2 Qh7+ 61.
Kg2 Qe4+ 62. f3 Qxa4 63. Qxg5 Qa2 64. e4 Qb3 65. Qg3 Qe6 66. e5 Ke3 67. f4+
Kd4 68. Qf2+ Ke4 69. Qf3+ Kd4 70. Kf2 Qf5 71. Qd1+ Ke4 72. Qf3+ Kd4 73. Kg3
Qg6+ 74. Qg4 Qd3+ 75. Qf3 Qg6+ 76. Qg4 Qd3+ 77. Kh4 Qh7+ 78. Qh5 Qe7+ 79.
Kg4 Qe6+ 80. Qf5 Qg8+ 81. Qg5 Qe6+ 82. f5 Qa2 83. Qf4+ Kd5 84. Kg5 Qb2 85.
f6 Qg2+ 86. Kh6 1-0
`

	a.NumLines = 2
	a.MaxTimeSec = 1
	a.Depth = 5
	rchan := make(chan *PgnResponse, 10)
	a.RChannel = rchan
	pa := NewPgnAnalyzer()
	go pa.DoAnalyze(a)
	ar := collectPgnResults(rchan)
	assert.NotNil(t, ar)
}
