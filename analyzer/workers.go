package analyzer

import "fmt"

// FenResponse is from worker to consumer
type FenResponse struct {
	Pvar  int32  // PV number
	Score int32  // the score of the move
	Pos   string // The position

	Error error // IF there is an error

	Worker int64

	Done bool // The end of the messages
}

type FenData struct {
	Fen      string
	UserMove string

	RChannel chan *FenResponse
}

var fenChan = make(chan *FenData)

// AnalyzeFenChannelSender - analyze this data.
func AnalyzeFenChannelSender(fd *FenData) {
	fenChan <- fd
}

type FenWorker struct {
	seq      int64
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

		f.analyzer.Fen = msg.Fen
		f.analyzer.UserMove = msg.UserMove
		rst, err := f.analyzer.AnalyzeFen()

		fr := &FenResponse{}
		fr.Done = true
		if err != nil {
			fr.Error = err
		} else {
			fr.Score = rst.Score
		}
		fr.Worker = f.seq

		msg.RChannel <- fr

		Utils.AdjustFenWorker(1)
	}

}
