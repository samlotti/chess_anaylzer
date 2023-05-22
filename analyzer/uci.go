package analyzer

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
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

type UciProcess struct {
	lock   sync.Mutex
	Engine string
	cmd    *exec.Cmd
	state  UciState
	stdin  io.ReadCloser
	stdout io.WriteCloser
	estate EngineState

	// options - The engine options.
	Options map[string]string
}

// NewUci - creates a new instance of the engine!
func NewUci(engine string) *UciProcess {
	p := UciProcess{Engine: engine}
	p.state = UciNotStarted
	p.estate = ENotReady
	p.Options = make(map[string]string)
	return &p
}

func (p *UciProcess) Start() error {
	p.cmd = exec.Command("../engines/" + p.Engine)
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
	//txt = strings.TrimPrefix(txt, "option name ")
	//sects := strings.SplitN(txt, " ", 2)
	//p.Options[sects[0]] = sects[1]
}
