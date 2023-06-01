package main

import (
	"chess_analyzer/analyzer"
	"chess_analyzer/httpservice"
	"fmt"
	"github.com/samlotti/blip/blipUtil"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Chess Analyzer")

	analyzer.Environment.EnginePath = "./engines/"

	fmt.Printf("Running the server:  http://localhost:8181\n")

	// Show timings of the individual template renders
	blipUtil.Instance().SetMonitor(&blipUtil.DebugBlipMonitor{})

	//http.HandleFunc("/", Index)
	//http.HandleFunc("/favicon.ico", http.NotFound)
	//http.HandleFunc("/users/listAll", UListAll)
	//http.HandleFunc("/users/listActive", UListActive)
	//http.HandleFunc("/users/view/", UView)

	http.HandleFunc("/chess/ai/fen", httpservice.AnalyzeFen)

	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}

}
