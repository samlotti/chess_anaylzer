package analyzer

import (
	"fmt"
	"github.com/samlotti/chess_anaylzer/uci"
	"strconv"
	"time"
)

/*


Green: good   +val
Yellow: inaccuracy - -> -.3
Orange: mistake -.301 -> -.9
Red: blunder  -> -.9

*/

const (
	DefaultAnalyzePerMoveSec = 15
	Verbose                  = true
)

// FenAnalyzer - can analyze a position.
type FenAnalyzer struct {
	Fen        string
	MoveNum    int    // halfmove count
	UserMove   string // move that is being analyzed
	Depth      int
	MaxTimeSec int
	NumPVLines int
}

func NewFenAnalyzer() *FenAnalyzer {
	return &FenAnalyzer{
		MaxTimeSec: DefaultAnalyzePerMoveSec,
		NumPVLines: 5,
	}
}

// AResults - Results of the analyzer.
type AResults struct {
	RCode RCode
	Err   error
	Done  bool

	UserMove string

	BestMode *ARBestMove
	Info     *ARInfo
}

func AResultsError(err error) *AResults {
	return &AResults{
		RCode: RCODE_ERROR,
		Err:   err,
		Done:  true,
	}
}

// Analyze the position
// Returns Best move, Your move, top X good moves.
// Top X the least losing moves
func (a *FenAnalyzer) Analyze(rchan chan *AResults) {

	// The best move value used as the baseline.
	// The diff between the best move and players move
	// Is used determine the response of inaccuracy / blunder ...
	u, err := uci.UciManager().GetUci("zahak")
	// u, err := uci.UciManager().GetUci("stockfish")
	if err != nil {
		rchan <- AResultsError(err)
		return
	}
	defer uci.UciManager().Return(u)

	cb := make(chan *uci.UciCallback, 10)
	u.SetAsyncChannel(cb)

	cbf := func() {
		if Verbose {
			println("Waiting for UCI responses")
		}

		for {
			cbc := <-cb
			if Verbose {
				fmt.Printf("From UCI: %v\n", cbc)
			}

			answer := &AResults{}

			answer.UserMove = a.UserMove

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

			if cbc.Err != nil {
				answer.Err = cbc.Err
			}

			if answer.Err != nil {
				answer.RCode = RCODE_ERROR
				answer.Done = true
			}

			rchan <- answer

			if answer.Done {
				if Verbose {
					println("analyzer done")
				}
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

	err = u.SetOption("MultiPV", strconv.Itoa(a.NumPVLines))
	if err != nil {
		rchan <- AResultsError(err)
		return
	}

	opts := &uci.GoOptions{
		Depth:      a.Depth,
		SearchMove: "",
		// Fen:        a.Fen,
	}
	err = u.SendGo(opts)
	if err != nil {
		rchan <- AResultsError(err)
		return
	}

	if a.MaxTimeSec <= 0 {
		a.MaxTimeSec = DefaultAnalyzePerMoveSec
	}

	err = u.WaitMoveUpTo(time.Duration(a.MaxTimeSec) * time.Second)
	if err != nil {
		rchan <- AResultsError(err)
		return
	}

}
