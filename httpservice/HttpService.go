package httpservice

import (
	"fmt"
	"github.com/samlotti/chess_anaylzer/analyzer"
	"log"
	"net/http"
)

// AnalyzeFen
// http://localhost:8181/chess/ai/fen?fen=2q1rr1k/3bbnnp/p2p1pp1/2pPp3/PpP1P1P1/1P2BNNP/2BQ1PRK/7R%20b%20-%20-

func AnalyzeFen(w http.ResponseWriter, r *http.Request) {
	fen, ok := r.URL.Query()["fen"]
	if ok {
		// w.Write([]byte(fen[0]))
	} else {
		w.Write([]byte("please enter ?fen= on the url"))
		return
	}

	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error processing request"))
			log.Printf("error %s\n", err)
		}
	}()

	fd := &analyzer.FenData{
		Fen:      fen[0],
		UserMove: "",
		RChannel: make(chan *analyzer.FenResponse),
	}
	analyzer.AnalyzeFenChannelSender(fd)
	fresp := <-fd.RChannel

	if fresp.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error processing request"))
		log.Printf("error %s\n", fresp.Error)
	} else {
		w.Write([]byte(fmt.Sprintf("Score: %d\n", fresp.Score)))
		w.Write([]byte(fmt.Sprintf("Worker: %d\n", fresp.Worker)))
		w.Write([]byte("request completed"))
	}

}
