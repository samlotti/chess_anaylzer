package httpservice

import (
	"encoding/json"
	"github.com/samlotti/chess_anaylzer/analyzer"
	"github.com/samlotti/chess_anaylzer/chessboard/common"
	"log"
	"net/http"
)

// AnalyzeFen
// http://localhost:8181/chess/ai/fen?fen=2q1rr1k/3bbnnp/p2p1pp1/2pPp3/PpP1P1P1/1P2BNNP/2BQ1PRK/7R%20b%20-%20-
//
// args:  fen  required
//
//	depth optional
func AnalyzeFen(w http.ResponseWriter, r *http.Request) {
	fen, ok := r.URL.Query()["fen"]
	if ok {
		// w.Write([]byte(fen[0]))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("please enter ?fen= on the url"))
		return
	}

	depth, err := common.Utils.ArgInt(r.URL.Query(), "depth", 15)
	if err != nil || depth < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("depth invalid, please enter a valid number"))
		return
	}

	tsec, err := common.Utils.ArgInt(r.URL.Query(), "tsec", 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("tsec invalid, please enter a valid number"))
		return
	}

	pvlines, err := common.Utils.ArgInt(r.URL.Query(), "lines", 0)
	if err != nil || pvlines < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("line invalid, please enter a valid number"))
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
		Fen:        fen[0],
		UserMove:   "",
		Depth:      depth,
		NumLines:   pvlines,
		MaxTimeSec: tsec,
		RChannel:   make(chan *analyzer.FenResponse),
	}

	wasSent := analyzer.AnalyzeFenChannelSender(fd)
	if !wasSent {
		w.Write([]byte("server busy"))
		return
	}

	for {
		fresp := <-fd.RChannel

		json.NewEncoder(w).Encode(fresp)

		if fresp.Done {
			//w.Write([]byte("request completed"))
			return
		}

	}

}
