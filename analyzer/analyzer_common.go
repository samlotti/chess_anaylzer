package analyzer

import "fmt"

const (
	RCODE_INFO     RCode = "info"
	RCODE_BESTMOVE       = "bm"
	RCODE_ERROR          = "error"
	RCODE_DONE           = "done"
)

type RCode string

type ARBestMove struct {
	BestMove string `json:"bestMove"`
	Ponder   string `json:"ponder"`
}

func (b *ARBestMove) String() string {
	return fmt.Sprintf("best:%s", b.BestMove)
}

type ARInfo struct {
	Depth      int      `json:"depth"`      // the depth of the move
	MPv        int      `json:"pv"`         // the PV number
	ScoreCP    int      `json:"score"`      // score in centipawns  100 = one pawn, > 15000 = mate in 15001=1, - = mated in
	MateIn     int      `json:"mateIn"`     // 0=no mate, + = current player mates,  - other player mates
	Moves      []string `json:"moves"`      // the moves
	Nps        int      `json:"nps"`        // the nodes per sec
	IsUserMove bool     `json:"isUserMove"` // this is the requsted user move
}

func (b *ARInfo) String() string {
	return fmt.Sprintf("%d: depth: %d Sc:%d  M:%d UM:%v > %v", b.MPv, b.Depth, b.ScoreCP, b.MateIn, b.IsUserMove, b.Moves)
}
