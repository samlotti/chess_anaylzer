package analyzer

import (
	"fmt"
	ai "github.com/samlotti/chess_anaylzer/chessboard"
)

/*
Thoughts
Make multiple passes,
one with a max depth of 3?
then a second with max depth of 5
then max time.

Also dont keep creating instances of the engine

*/

var pgnChan = make(chan *PgnData, 1)

// PgnResponse is from worker to consumer
// Will send these over the channel
type PgnResponse struct {
	RCode      RCode       `json:"rcode"`
	Error      string      `json:"error"` // IF there is an error
	ARBestMove *ARBestMove `json:"bestMove"`
	ARInfo     *ARInfo     `json:"info"`
	MoveNum    int         `json:"moveNum"`
	PlayedMove string      `json:"playedmove"`
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

	// Do a quick run to make sure no errors.
	brd := createNewBoard(wrapper)
	for _, mv := range wrapper.InternalMoves {
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
		// fen := ai.BoardToFen(brd, i)
		// fmt.Printf("%s \n", fen)

		// brd.PrintBoard("dd")
	}
	//
	//msg.RChannel <- &PgnResponse{
	//	RCode: RCODE_ERROR,
	//	Error: fmt.Sprintf("Code Not complete! moves %v", wrapper.Moves),
	//	Done:  true,
	//}

	var fenAnalyzer = NewFenAnalyzer()
	fenAnalyzer.KeepProcess = true
	defer fenAnalyzer.Close()

	brd = createNewBoard(wrapper)
	// The initial board fen ... will analyze using fen
	fen := ai.BoardToFen(brd, 0)
	for i, mv := range wrapper.InternalMoves {

		ims := ai.MoveToInputString(mv)
		algstr := ai.MoveToString(mv)

		fenAnalyzer.NumPVLines = msg.NumLines
		fenAnalyzer.MaxTimeSec = msg.MaxTimeSec
		fenAnalyzer.Depth = msg.Depth
		fenAnalyzer.Fen = fen
		fenAnalyzer.UserMove = algstr
		fenAnalyzer.MoveNum = i + 1

		fmt.Printf("In: %s = %s \n", ims, fen)
		err = f.doAnalyzeThisMove(fenAnalyzer, msg)
		if err != nil {
			msg.RChannel <- &PgnResponse{
				RCode: RCODE_ERROR,
				Error: err.Error(),
				Done:  true,
			}
			return
		}
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

func createNewBoard(wrapper *ai.PgnWrapper) *ai.Board {
	brd := ai.NewBoard()
	if len(wrapper.StartFen) > 0 {
		ai.ParseFen(brd, wrapper.StartFen)
	} else {
		ai.ParseFen(brd, ai.StartFen)
	}
	return brd
}

func (f *PgnAnalyzer) doAnalyzeThisMove(fenAnalyzer *FenAnalyzer, msg *PgnData) error {
	dchan := make(chan struct{}, 10)
	rchan := make(chan *AResults, 10)

	var err error = nil

	go func() {
		for {
			m := <-rchan

			fr := &PgnResponse{}
			fr.ARInfo = m.Info
			fr.ARBestMove = m.BestMode
			if m.Err != nil {
				fr.Error = m.Err.Error()
				err = m.Err
			}

			fr.PlayedMove = m.UserMove
			fr.MoveNum = m.MoveNumber
			fr.Done = false // m.Done
			fr.RCode = m.RCode

			//if fr.RCode == RCODE_DONE {
			//	// Just ignore these, as they're from the fen search
			//	continue
			//}

			// fmt.Printf("send from fen: %+v\n", fr)

			if fr.RCode != RCODE_DONE {
				msg.RChannel <- fr
			}

			if m.Done {
				dchan <- struct{}{}
				return
			}

		}
	}()

	fenAnalyzer.Analyze(rchan)

	// wait for complete
	<-dchan

	return err

}
