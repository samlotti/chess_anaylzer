package analyzer

import (
	"fmt"
	. "github.com/samlotti/chess_anaylzer/uci"
	"time"
)

/*


Green: good   +val
Yellow: inaccuracy - -> -.3
Orange: mistake -.301 -> -.9
Red: blunder  -> -.9

*/

// Analyzer - can analyze a position.
type Analyzer struct {
	Fen      string
	UserMove string
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

type RCode string

const (
	RCODE_INFO     RCode = "info"
	RCODE_BESTMOVE       = "bm"
	RCODE_ERROR          = "error"
)

type ARBestMove struct {
	BestMove string
	Ponder   string
}

type ARInfo struct {
	Depth   int      // the depth of the move
	MPv     int      // the PV number
	ScoreCP int      // score in centipawns  100 = one pawn, > 15000 = mate in 15001=1, - = mated in
	MateIn  int      // 0=no mate, + = current player mates,  - other player mates
	Moves   []string // the moves
	Nps     int      // the nodes per sec
}

// AResults - Results of the analyzer.
type AResults struct {
	RCode RCode
	Err   error
	Done  bool

	BestMode *ARBestMove
	Info     *ARInfo
}

func AResultsError(err error) *AResults {
	return &AResults{
		RCode: RCODE_ERROR,
		Err:   err,
	}
}

// AnalyzeFen the position
// Returns Best move, Your move, top X good moves.
// Top X the least losing moves
func (a *Analyzer) AnalyzeFen(rchan chan *AResults) {

	// The best move value used as the baseline.
	// The diff between the best move and players move
	// Is used determine the response of inaccuracy / blunder ...

	u, err := UciManager().GetUci("zahak")
	if err != nil {
		rchan <- AResultsError(err)
		return
	}
	defer UciManager().Return(u)

	cb := make(chan *UciCallback, 10)
	u.SetAsyncChannel(cb)

	cbf := func() {
		println("Waiting for CB")
		for {
			cbc := <-cb
			fmt.Printf("CB: %v\n", cbc)

			answer := &AResults{}
			if cbc.BestMove != nil {
				answer.RCode = RCODE_BESTMOVE
				answer.Err = cbc.BestMove.Err
				answer.BestMode = &ARBestMove{}
				answer.BestMode.BestMove = cbc.BestMove.BestMove
				answer.BestMode.Ponder = cbc.BestMove.Ponder
				answer.Done = true
			}
			if cbc.Info != nil {
				answer.RCode = RCODE_INFO
				answer.Err = cbc.Info.Err
				answer.Info = &ARInfo{}
				answer.Info.Moves = cbc.Info.Moves
				answer.Info.Depth = cbc.Info.Depth
				answer.Info.Nps = cbc.Info.Nps
				answer.Info.ScoreCP = cbc.Info.ScoreCP
				answer.Info.MPv = cbc.Info.MPv
				answer.Info.MateIn = cbc.Info.MateIn
			}

			if answer.Err != nil {
				answer.RCode = RCODE_ERROR
			}

			rchan <- answer

			if answer.Done {
				println("DONE")
				return
			}
		}
	}
	go cbf()

	err = u.SetPositionFen(a.Fen)
	if err != nil {
		rchan <- AResultsError(err)
		return
	}

	err = u.SetOption("MultiPV", "10")
	if err != nil {
		rchan <- AResultsError(err)
		return
	}

	opts := &GoOptions{
		Depth:      15,
		SearchMove: "",
	}
	err = u.SendGo(opts)
	if err != nil {
		rchan <- AResultsError(err)
		return
	}

	err = u.WaitMoveUpTo(5 * time.Second)
	if err != nil {
		rchan <- AResultsError(err)
		return
	}

}
