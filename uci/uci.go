package uci

import (
	"bufio"
	"fmt"
	"github.com/samlotti/chess_anaylzer/chessboard/common"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**

https://www.wbec-ridderkerk.nl/html/UCIProtocol.html

http://chess.grantnet.us/index/3/

// setoption name MultiPV value 5


*/

// AIRequest - The request to run the analysis
type AIRequest struct {
	Pgn string `json:"pgn"`
	Fen string `json:"fen"`
}

// AIMoveData - The data for a single move in the game.
// Contains the move string and the move score.
type AIMoveData struct {
	MoveStr   string `json:"moveStr"`
	MoveScore int    `json:"MoveScore"`
}

// AIMoveResponse - The response for a single move,
// contains all the move scores
// moved is the actual move in the pgn
// moves are the calculated moves
type AIMoveResponse struct {
	MoveNum int           `json:"moveNum"`
	Moved   *AIMoveData   `json:"moved"`
	Moves   []*AIMoveData `json:"moves"`
}

type AIResponse struct {
}

type UciState int16

const (
	UciNotStarted UciState = 0
	UciRunning             = 1
	UciStopped             = 2
	UciFailed              = 3
)

type EngineState int16

const (
	ENotReady    EngineState = 0
	EOk                      = 1
	ECalculating             = 2
)

type UciBestMove struct {
	Err      error
	BestMove string
	Ponder   string
}

type UciCallback struct {
	Raw      string
	BestMove *UciBestMove
	Info     *UciInfo
}

// UciBestMoveParse - parses the best move line
//
// bestmove f4h6 ponder g7h6
func UciBestMoveParse(bm string) *UciBestMove {
	var sections = strings.Split(bm, " ")
	if sections[0] != "bestmove" {
		return nil
	}

	res := &UciBestMove{}
	pos := 0
	for {
		if pos >= len(sections) {
			break
		}
		cmd := sections[pos]
		switch cmd {
		case "bestmove":
			res.BestMove = sections[pos+1]
			pos += 2
		case "ponder":
			res.Ponder = sections[pos+1]
			pos += 2
		default:
			// Skip the command
			pos += 1
		}
	}
	return res
}

type UciInfo struct {
	Err     error
	Depth   int      // the depth of the move
	MPv     int      // the PV number
	ScoreCP int      // score in centipawns  100 = one pawn, > 15000 = mate in 15001=1, - = mated in
	MateIn  int      // 0=no mate, + = current player mates,  - other player mates
	Moves   []string // the moves
	Nps     int      // the nodes per sec
}

func UciInfoParse(info string) *UciInfo {

	var sections = strings.Split(info, " ")
	if sections[0] != "info" {
		return nil
	}

	var err error
	res := &UciInfo{}
	pos := 1
	for {
		if pos >= len(sections) {
			break
		}
		cmd := sections[pos]
		switch cmd {
		case "depth":
			res.Depth, err = strconv.Atoi(sections[pos+1])
			if err != nil {
				res.Err = err
				return res
			}
			pos += 2
		case "nps":
			res.Nps, err = strconv.Atoi(sections[pos+1])
			if err != nil {
				res.Err = err
				return res
			}
			pos += 2
		case "multipv":
			res.MPv, err = strconv.Atoi(sections[pos+1])
			if err != nil {
				res.Err = err
				return res
			}
			pos += 2
		case "score":
			pos += 1
		case "cp":
			res.ScoreCP, err = strconv.Atoi(sections[pos+1])
			if err != nil {
				res.Err = err
				return res
			}
			pos += 2
		case "mate":
			res.MateIn, err = strconv.Atoi(sections[pos+1])
			res.ScoreCP, err = strconv.Atoi(sections[pos+1])
			if res.ScoreCP > 0 {
				res.ScoreCP = 15000 + res.ScoreCP
			} else {
				res.ScoreCP = -15000 + res.ScoreCP
			}
			if err != nil {
				res.Err = err
				return res
			}
			pos += 2
		case "pv":
			pos++
			for {
				if pos >= len(sections) {
					break
				}
				res.Moves = append(res.Moves, sections[pos])
				pos++
			}
		default:
			// Skip the command
			pos += 1
		}
	}
	return res
}

type UciProcess struct {
	lock   sync.Mutex
	Engine string
	epath  string
	cmd    *exec.Cmd
	state  UciState
	stdin  io.ReadCloser
	stdout io.WriteCloser
	estate EngineState

	callback chan *UciCallback

	// options - The engine options.
	Options map[string]string
}

// NewUci - creates a new instance of the engine!
func NewUci(engine string) *UciProcess {
	p := UciProcess{Engine: engine, epath: common.Environment.EnginePath}
	p.state = UciNotStarted
	p.estate = ENotReady
	p.Options = make(map[string]string)
	return &p
}

func (p *UciProcess) Start() error {
	// p.cmd = exec.Command("./engines/" + p.Engine)
	p.cmd = exec.Command(p.epath + p.Engine)
	var err error

	p.stdin, err = p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	p.stdout, err = p.cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err = p.cmd.Start(); err != nil {
		p.state = UciFailed
		return err
	}

	fmt.Println("Command running")

	p.state = UciRunning

	go p.monitor()

	p.send("uci")

	p.send("isready")

	err = p.WaitOk(10 * time.Second)

	return err
}

func (p *UciProcess) WaitOk(timeout time.Duration) error {
	st := time.Now()
	for {
		if p.GetEState() == EOk {
			return nil
		}

		if time.Now().Sub(st) > timeout {
			return fmt.Errorf("timeout waiting for engine ok")
		}

		time.Sleep(100 * time.Millisecond)

	}
}

func (p *UciProcess) monitor() {
	defer func() {
		p.state = UciStopped
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			p.state = UciFailed
		}
	}()
	scanner := bufio.NewScanner(p.stdin)
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Printf("Have: %s\n", txt)

		if p.callback != nil {
			p.callback <- &UciCallback{
				Raw:      txt,
				BestMove: UciBestMoveParse(txt),
				Info:     UciInfoParse(txt),
			}
		}

		if txt == "uciok" {
			p.SetEState(EOk)
		}

		if txt == "readyok" {
			p.SetEState(EOk)
		}

		if strings.HasPrefix(txt, "bestmove") {
			p.setBestMove(txt)
		}

		if strings.HasPrefix(txt, "option") {
			p.addOption(txt)
		}

		if strings.HasPrefix(txt, "info") {
			p.addMoveInfo(txt)
		}

	}
}

// send - Sends a line to the chess engine.
func (p *UciProcess) send(line string) error {
	_, err := fmt.Fprint(p.stdout, line)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(p.stdout, "\n")
	if err != nil {
		return err
	}
	return nil
}

// Terminate - Terminate the process
func (p *UciProcess) Terminate() {
	p.state = UciStopped
	_ = p.stdout.Close()
	_ = p.stdin.Close()
	_ = p.cmd.Process.Kill()
}

func (p *UciProcess) GetEState() EngineState {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.estate
}

func (p *UciProcess) SetEState(state EngineState) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.estate = state
}

// SetPositionFen - Will set the position
// The engine must be in the ready state
func (p *UciProcess) SetPositionFen(fen string) error {
	if err := p.checkReady(); err != nil {
		return err
	}
	return p.send("position fen " + fen)
}

func (p *UciProcess) checkReady() error {
	if p.GetEState() != EOk {
		return fmt.Errorf("engine not in ready state")
	}
	return nil
}

type GoOptions struct {
	Depth      int
	SearchMove string
}

func NewGoOptions() *GoOptions {
	return &GoOptions{}
}

// SendGo - Send the go command to the engine, pass in options to
// configure the command.  default will run indefinitely.
func (p *UciProcess) SendGo(opts *GoOptions) error {
	if err := p.checkReady(); err != nil {
		return err
	}

	p.SetEState(ECalculating)

	str := "go "
	if opts.Depth > 0 {
		str = fmt.Sprintf("%s depth %d", str, opts.Depth)
	}
	if len(opts.SearchMove) > 0 {
		str = fmt.Sprintf("%s searchmoves %s", str, opts.SearchMove)
	}

	return p.send(str)

}

//// Stop - Tell the engine to stop calculating.
//func (p *UciProcess) Stop() error {
//	if p.GetEState() != EOk {
//		return fmt.Errorf("engine not in ready state")
//	}
//	return nil
//}

// SendStop - Tell the engine to stop calculating.
func (p *UciProcess) SendStop() error {
	return p.send("stop")
}

// SetOption - set the options for the engine.
// The options are stores in the engine Options map.
func (p *UciProcess) SetOption(name string, val string) error {
	return p.send(fmt.Sprintf("setoption name %s value %s", name, val))
}

// WaitMoveUpTo - wait for the move to complete, or send stop then wait a bit more
func (p *UciProcess) WaitMoveUpTo(timeout time.Duration) error {
	err := p.WaitOk(timeout)
	if err == nil {
		return nil

	}

	p.SendStop()
	err = p.WaitOk(500 * time.Millisecond)
	return err
}

func (p *UciProcess) IsReadyForMove() bool {
	return p.GetEState() == EOk
}

// addOption - parses the option line from the engine
// ex: option name Ponder type check default false
func (p *UciProcess) addOption(txt string) {
	txt = strings.TrimPrefix(txt, "option name ")
	sects := strings.SplitN(txt, " ", 2)
	p.Options[sects[0]] = sects[1]
}

/**
info depth 15 seldepth 30 hashfull 58 tbhits 0 nodes 1720962 nps 1103672 score cp -124 time 1559 multipv 5 pv c8c7 f3e1 f8g8 h1g1 a6a5 f2f3 e7d8 d2f2 d8e7 e1d3 e8a8 f2d2 g7e8 h3h4 g8g7
bestmove f8g8 ponder h1f1
*/

// setBestMove - parses the option line from the engine
// ex: option name Ponder type check default false
func (p *UciProcess) setBestMove(txt string) {
	// bestmove f8g8 ponder h1f1
	//txt = strings.TrimPrefix(txt, "option name ")
	//sects := strings.SplitN(txt, " ", 2)
	//p.Options[sects[0]] = sects[1]

	p.SetEState(EOk)
}

// addMoveInfo - parses the option line from the engine
// ex: option name Ponder type check default false
func (p *UciProcess) addMoveInfo(txt string) {
	// info depth 15 seldepth 30 hashfull 58 tbhits 0 nodes 1720962 nps 1103672 score cp -124 time 1559 multipv 5 pv c8c7 f3e1 f8g8 h1g1 a6a5 f2f3 e7d8 d2f2 d8e7 e1d3 e8a8 f2d2 g7e8 h3h4 g8g7
	// info depth 4 seldepth 12 hashfull 0 tbhits 0 nodes 28043 nps 1016662 score mate +4 time 27 multipv 1 pv d7d1 e2d1 a3a2 b1c1 c5a3 e5b2
	//txt = strings.TrimPrefix(txt, "option name ")
	//sects := strings.SplitN(txt, " ", 2)
	//p.Options[sects[0]] = sects[1]

}

func (p *UciProcess) SetAsyncChannel(callbacks chan *UciCallback) {
	p.callback = callbacks
}
