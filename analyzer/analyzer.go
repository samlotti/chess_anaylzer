package analyzer

import (
	"time"
)

/*


Green: good   +val
Yellow: inaccuracy - -> -.3
Orange: mistake -.301 -> -.9
Red: blunder  -> -.9

*/

// Analyzer - can analyze a position.
type Analyzer struct {
	Fen      string
	UserMove string
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AResults - Results of the analyzer.
type AResults struct {
}

// AnalyzeFen the position
// Returns Best move, Your move, top X good moves.
// Top X the least losing moves
func (a *Analyzer) AnalyzeFen() (*AResults, error) {

	// The best move value used as the baseline.
	// The diff between the best move and players move
	// Is used determine the response of inaccuracy / blunder ...

	u, err := UciManager().GetUci("zahak")
	if err != nil {
		return nil, err
	}
	defer UciManager().Return(u)

	answer := &AResults{}

	err = u.SetPositionFen(a.Fen)
	if err != nil {
		return nil, err
	}

	err = u.SetOption("MultiPV", "10")
	if err != nil {
		return nil, err
	}

	opts := &GoOptions{
		Depth:      15,
		SearchMove: "",
	}
	err = u.SendGo(opts)
	if err != nil {
		return nil, err
	}

	err = u.WaitMoveUpTo(5 * time.Second)
	if err != nil {
		return nil, err
	}

	return answer, nil

}
