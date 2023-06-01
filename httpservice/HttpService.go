package httpservice

import (
	"chess_analyzer/analyzer"
	"log"
	"net/http"
)

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

	analyzer := analyzer.NewAnalyzer()
	analyzer.Fen = fen[0]
	analyzer.UserMove = ""
	_, err := analyzer.AnalyzeFen()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error processing request"))
		log.Printf("error %s\n", err)
	} else {
		w.Write([]byte("request completed"))
	}

}
