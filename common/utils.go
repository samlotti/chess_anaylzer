package common

import (
	"fmt"
	"net/url"
	"strconv"
	"sync/atomic"
)

const Verbose = true

type utils struct {
	seq            int64
	availFenWorker int64
	availPgnWorker int64
}

var Utils = &utils{
	seq:            0,
	availFenWorker: 0,
	availPgnWorker: 0,
}

func (u *utils) NextSeq() int64 {
	return atomic.AddInt64(&u.seq, 1)
}

func (u *utils) debug() {
	if Verbose {
		p := atomic.LoadInt64(&u.availPgnWorker)
		f := atomic.LoadInt64(&u.availFenWorker)
		fmt.Printf("Num workers: Pgn %d,  Fen: %d\n", p, f)
	}
}

// AdjustPgnWorker -- inc or dec the number of available fen workers
func (u *utils) AdjustPgnWorker(num int64) {
	atomic.AddInt64(&u.availPgnWorker, num)
	u.debug()
}

// GetPgnWorkers - the number of available fen workers.
func (u *utils) GetPgnWorkers() int64 {
	r := atomic.LoadInt64(&u.availPgnWorker)
	u.debug()
	return r
}

// AdjustFenWorker -- inc or dec the number of available fen workers
func (u *utils) AdjustFenWorker(num int64) {
	atomic.AddInt64(&u.availFenWorker, num)
	u.debug()
}

// GetFenWorkers - the number of available fen workers.
func (u *utils) GetFenWorkers() int64 {
	r := atomic.LoadInt64(&u.availFenWorker)
	u.debug()
	return r
}

// ArgInt - returns the int value of a string from the map
func (u *utils) ArgInt(data url.Values, key string, dflt int) (int, error) {
	val, ok := data[key]
	if !ok {
		return dflt, nil
	}
	if len(val) != 1 {
		return 0, fmt.Errorf("invalid value for %s", key)
	}
	ival, err := strconv.Atoi(string(val[0]))
	return ival, err
}

func (u *utils) AToI(val string, dflt int) (int, error) {
	if len(val) == 0 {
		return dflt, nil
	}
	ival, err := strconv.Atoi(val)
	return ival, err

}
