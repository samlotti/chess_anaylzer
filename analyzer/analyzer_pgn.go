package analyzer

import (
	"fmt"
	ai "github.com/samlotti/chess_anaylzer/chessboard"
)

var pgnChan = make(chan *PgnData, 1)

// PgnResponse is from worker to consumer
// Will send these over the channel
type PgnResponse struct {
	RCode      RCode       `json:"rcode"`
	Error      string      `json:"error"` // IF there is an error
	ARBestMove *ARBestMove `json:"bestMove"`
	ARInfo     *ARInfo     `json:"info"`
	ARYourMove *ARBestMove `json:"yourMove"`
	MoveNum    int         `json:"moveNum"`
	Done       bool        `json:"done"` // The end of the messages
}

type PgnData struct {
	Pgn   string
	Depth int

	RChannel chan *PgnResponse

	MaxTimeSec int
	NumLines   int
}

// PgnAnalyzer - can analyze a position.
type PgnAnalyzer struct {
	Pgn        string
	Depth      int
	MaxTimeSec int
	NumPVLines int
}

func NewPgnAnalyzer() *PgnAnalyzer {
	return &PgnAnalyzer{
		MaxTimeSec: DefaultAnalyzePerMoveSec,
		NumPVLines: 5,
	}
}

func (f *PgnAnalyzer) DoAnalyze(msg *PgnData) {

	/***
	Need to loop through each move, passing the fen of the board
	*/

	wrapper := ai.NewPgnWrapper(msg.Pgn)
	err := wrapper.Parse()
	if err != nil {
		msg.RChannel <- &PgnResponse{
			RCode: RCODE_ERROR,
			Error: err.Error(),
			Done:  true,
		}
		return
	}

	brd := ai.NewBoard()
	if len(wrapper.StartFen) > 0 {
		ai.ParseFen(brd, wrapper.StartFen)
	} else {
		ai.ParseFen(brd, ai.StartFen)
	}
	for i, mv := range wrapper.InternalMoves {

		ims := ai.MoveToInputString(mv)
		err = brd.MakeMove(mv, ims)
		if err != nil {
			msg.RChannel <- &PgnResponse{
				RCode: RCODE_ERROR,
				Error: err.Error(),
				Done:  true,
			}
			return
		}
		fen := ai.BoardToFen(brd, i)
		fmt.Printf("%s \n", fen)
	}

	//msg.RChannel <- &PgnResponse{
	//	RCode: RCODE_ERROR,
	//	Error: fmt.Sprintf("Code Not complete! moves %v", wrapper.Moves),
	//	Done:  true,
	//}

	var fenAnalyzer = NewFenAnalyzer()

	brd = ai.NewBoard()
	if len(wrapper.StartFen) > 0 {
		ai.ParseFen(brd, wrapper.StartFen)
	} else {
		ai.ParseFen(brd, ai.StartFen)
	}

	// The initial board fen
	fen := ai.BoardToFen(brd, 0)
	for i, mv := range wrapper.InternalMoves {

		ims := ai.MoveToInputString(mv)

		fenAnalyzer.NumPVLines = msg.NumLines
		fenAnalyzer.MaxTimeSec = msg.MaxTimeSec
		fenAnalyzer.Depth = msg.Depth
		fenAnalyzer.Fen = fen
		fenAnalyzer.UserMove = ims
		fenAnalyzer.MoveNum = i + 1

		fmt.Printf("In: %s = %s \n", ims, fen)
		f.doAnalyzeThisMove(fenAnalyzer)
		// fmt.Printf("Out: %s = %s \n", ims, fen)

		err = brd.MakeMove(mv, ims)
		if err != nil {
			msg.RChannel <- &PgnResponse{
				RCode: RCODE_ERROR,
				Error: err.Error(),
				Done:  true,
			}
			return
		}

		fen = ai.BoardToFen(brd, i)

	}
	msg.RChannel <- &PgnResponse{
		RCode: RCODE_DONE,
		Done:  true,
	}

}

func (f *PgnAnalyzer) doAnalyzeThisMove(fenAnalyzer *FenAnalyzer) {
	dchan := make(chan struct{}, 10)
	rchan := make(chan *AResults, 10)
	go func() {
		for {
			m := <-rchan

			fr := &PgnResponse{}
			fr.ARInfo = m.Info
			fr.ARBestMove = m.BestMode
			if m.Err != nil {
				fr.Error = m.Err.Error()
			}

			fr.Done = m.Done
			fr.RCode = m.RCode
			fmt.Printf("send from fen: %+v\n", fr)
			// msg.RChannel <- fr

			if m.Done {
				dchan <- struct{}{}
				return
			}

		}
	}()

	fenAnalyzer.Analyze(rchan)

	// wait for complete
	<-dchan

}
