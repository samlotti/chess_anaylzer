package httpservice

import (
	"encoding/json"
	"github.com/samlotti/chess_anaylzer/analyzer"
	"github.com/samlotti/chess_anaylzer/chessboard/common"
	"log"
	"net/http"
)

// AnalyzePgn
// args:  fen  required
//
//	depth optional
func AnalyzePgn(w http.ResponseWriter, r *http.Request) {
	pgn := r.PostFormValue("pgn")
	if len(pgn) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("please enter a pgn"))
		return
	}

	depth, err := common.Utils.AToI(r.PostFormValue("depth"), 15)
	if err != nil || depth < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("depth invalid, please enter a valid number"))
		return
	}

	tsec, err := common.Utils.AToI(r.PostFormValue("tsec"), 15)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("tsec invalid, please enter a valid number"))
		return
	}

	pvlines, err := common.Utils.AToI(r.PostFormValue("lines"), 5)
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

	fd := &analyzer.PgnData{
		Pgn:        pgn,
		Depth:      depth,
		NumLines:   pvlines,
		MaxTimeSec: tsec,
		RChannel:   make(chan *analyzer.PgnResponse),
	}

	wasSent := analyzer.AnalyzePgnChannelSender(fd)
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
