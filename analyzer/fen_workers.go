package analyzer

import (
	"fmt"
	. "github.com/samlotti/chess_anaylzer/chessboard/common"
	"sync"
)

// FenResponse is from worker to consumer
type FenResponse struct {
	RCode      RCode       `json:"rcode"`
	Error      string      `json:"error"` // IF there is an error
	ARBestMove *ARBestMove `json:"bestMove"`
	ARInfo     *ARInfo     `json:"info"`
	Done       bool        `json:"done"` // The end of the messages
}

type FenData struct {
	Fen      string
	UserMove string
	Depth    int

	RChannel chan *FenResponse

	MaxTimeSec int
	NumLines   int
}

var mutex = sync.Mutex{}

var fenChan = make(chan *FenData, 1)

// AnalyzeFenChannelSender - analyze this data.
// Returns false if the queue are busy
// true if the fen was sent
func AnalyzeFenChannelSender(fd *FenData) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if len(fenChan) == cap(fenChan) {
		return false
	}
	fenChan <- fd
	return true
}

type FenWorker struct {
	seq int64

	analyzer *Analyzer
}

var fenWorkers = make([]*FenWorker, 0)

func CreateFenWorkers(num int) {
	for i := 1; i <= num; i++ {
		var fw = &FenWorker{
			seq: Utils.NextSeq(),
		}
		go fw.runLoop()
		fenWorkers = append(fenWorkers, fw)
	}
}

// runLoop -- Worker that waits for fen requests
func (f *FenWorker) runLoop() {
	println("Fen worker started: ", f.seq)
	Utils.AdjustFenWorker(1)
	f.analyzer = NewAnalyzer()
	for {
		msg := <-fenChan
		Utils.AdjustFenWorker(-1)
		fmt.Printf("Fen worker started %d: %s\n", f.seq, msg.Fen)

		f.analyzer.NumPVLines = msg.NumLines
		f.analyzer.MaxTimeSec = msg.MaxTimeSec
		f.analyzer.Depth = msg.Depth
		f.analyzer.Fen = msg.Fen
		f.analyzer.UserMove = msg.UserMove

		dchan := make(chan struct{}, 10)
		rchan := make(chan *AResults, 10)
		go func() {
			for {
				m := <-rchan

				fr := &FenResponse{}
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

		// whait for complete
		<-dchan

		Utils.AdjustFenWorker(1)

	}

}
