package analyzer

/*


Green: good   +val
Yellow: inaccuracy - -> -.3
Orange: mistake -.301 -> -.9
Red: blunder  -> -9

*/

// Analyzer - can analyze a position.
type Analyzer struct {
	Fen      string
	userMove string
}

// AResults - Results of the analyzer.
type AResults struct {
}

// Analyze the position
// Returns Best move, Your move, top X good moves.
// Top X the least losing moves
func (a *Analyzer) Analyze() {

}
