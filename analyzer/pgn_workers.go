package analyzer

import (
	"fmt"
	"github.com/samlotti/chess_anaylzer/chessboard/common"
)

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
	analyzer *PgnAnalyzer
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
	f.analyzer = NewPgnAnalyzer()
	for {
		msg := <-pgnChan

		common.Utils.AdjustPgnWorker(-1)
		if Verbose {
			fmt.Printf("Pgn worker started %d: %s\n", f.seq, msg.Pgn)
		}

		f.analyzer.DoAnalyze(msg)
		common.Utils.AdjustPgnWorker(1)

	}
}
