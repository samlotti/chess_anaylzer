package main

import (
	"fmt"
	"github.com/samlotti/blip/blipUtil"
	"github.com/samlotti/chess_anaylzer/analyzer"
	. "github.com/samlotti/chess_anaylzer/chessboard/common"
	"github.com/samlotti/chess_anaylzer/httpservice"
	"log"
	"net/http"
)

//go:generate blip -dir ../template
func main() {
	fmt.Println("Chess FenAnalyzer")

	// The number of workers
	analyzer.CreateFenWorkers(5)
	analyzer.CreatePgnWorkers(5)

	Environment.EnginePath = "../engines/"

	fmt.Printf("Running the server:  http://localhost:8181\n")

	fmt.Printf("http://localhost:8181/chess/ai/fen?fen=" +
		"2r3k1/p4p2/3Rp2p/1p2P1pK/8/1P4P1/P3Q2P/1q6 b - - 0 1" + "\n")

	// Show timings of the individual template renders
	blipUtil.Instance().SetMonitor(&blipUtil.DebugBlipMonitor{})

	//http.HandleFunc("/", Index)
	//http.HandleFunc("/favicon.ico", http.NotFound)
	//http.HandleFunc("/users/listAll", UListAll)
	//http.HandleFunc("/users/listActive", UListActive)
	//http.HandleFunc("/users/view/", UView)

	http.HandleFunc("/chess/ai/pgn", httpservice.AnalyzePgn)
	http.HandleFunc("/chess/ai/fen", httpservice.AnalyzeFen)

	http.Handle("/", http.FileServer(http.Dir("../public")))
	// http.HandleFunc("/", httpservice.Index)

	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}

}
