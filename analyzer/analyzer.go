package analyzer

import (
	. "github.com/samlotti/chess_anaylzer/uci"
	"strings"
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
	RCODE_DONE     RCode = "done"
	RCODE_ERROR          = "err"
	RCODE_BESTMOVE       = "bm"
)

// AResults - Results of the analyzer.
type AResults struct {
	RCode    RCode
	Err      error
	Score    int32
	Done     bool
	BestMove string
	Ponder   string
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
			println("CB: ", cbc.Raw)

			answer := &AResults{}

			rchan <- answer

			if strings.HasPrefix(cbc.Raw, "best") {
				println("done!!")
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
