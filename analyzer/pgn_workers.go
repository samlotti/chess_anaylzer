package analyzer

import (
	"fmt"
	"github.com/samlotti/chess_anaylzer/chessboard/common"
)

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

var pgnChan = make(chan *PgnData, 1)

// AnalyzePgnChannelSender - analyze this data.
// Returns false if the queue are busy
// true if the fen was sent
func AnalyzePgnChannelSender(fd *PgnData) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if len(pgnChan) == cap(pgnChan) {
		return false
	}
	pgnChan <- fd
	return true
}

type PgnWorker struct {
	seq      int64
	analyzer *Analyzer
}

var pgnWorkers = make([]*PgnWorker, 0)

func CreatePgnWorkers(num int) {
	for i := 1; i <= num; i++ {
		var fw = &PgnWorker{
			seq: common.Utils.NextSeq(),
		}
		go fw.runLoop()
		pgnWorkers = append(pgnWorkers, fw)
	}
}

// runLoop -- Worker that waits for fen requests
func (f *PgnWorker) runLoop() {
	println("Pgn worker started: ", f.seq)
	common.Utils.AdjustPgnWorker(1)
	f.analyzer = NewAnalyzer()
	for {
		msg := <-pgnChan

		common.Utils.AdjustPgnWorker(-1)
		if Verbose {
			fmt.Printf("Pgn worker started %d: %s\n", f.seq, msg.Pgn)
		}

		f.doAnalyze(msg)
		common.Utils.AdjustPgnWorker(1)

	}
}

func (f *PgnWorker) doAnalyze(msg *PgnData) {

	/***
	Need to loop thru each move, passing the fen of the board
	*/

	msg.RChannel <- &PgnResponse{
		RCode: RCODE_ERROR,
		Error: "Code Not complete!",
		Done:  true,
	}
	return

	f.analyzer.NumPVLines = msg.NumLines
	f.analyzer.MaxTimeSec = msg.MaxTimeSec
	f.analyzer.Depth = msg.Depth
	f.analyzer.Fen = msg.Pgn
	// f.analyzer.UserMove = msg.UserMove

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
			msg.RChannel <- fr

			if m.Done {
				dchan <- struct{}{}
				return
			}

		}
	}()

	f.analyzer.AnalyzeFen(rchan)

	// wait for complete
	<-dchan

}
