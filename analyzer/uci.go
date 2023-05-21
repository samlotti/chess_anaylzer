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
	return nil
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

		if strings.HasPrefix(txt, "bestmove") {
			p.SetEState(EOk)
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

func newUci(engine string) *UciProcess {
	p := UciProcess{Engine: engine}
	p.state = UciNotStarted
	p.estate = ENotReady
	return &p
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

func (p *UciProcess) Stop() error {
	if p.GetEState() != EOk {
		return fmt.Errorf("engine not in ready state")
	}
	return nil
}

func (p *UciProcess) SendStop() error {
	return p.send("stop")
}
