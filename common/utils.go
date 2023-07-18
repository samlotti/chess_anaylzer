package common

import (
	"fmt"
	"net/url"
	"strconv"
	"sync/atomic"
)

type utils struct {
	seq            int64
	availFenWorker int64
}

var Utils = &utils{
	seq:            0,
	availFenWorker: 0,
}

func (u *utils) NextSeq() int64 {
	return atomic.AddInt64(&u.seq, 1)
}

// AdjustFenWorker -- inc or dec the number of available fen workerss
func (u *utils) AdjustFenWorker(num int64) {
	atomic.AddInt64(&u.availFenWorker, num)
}

// GetFenWorkers - the number of available fen workers.
func (u *utils) GetFenWorkers() int64 {
	return atomic.LoadInt64(&u.availFenWorker)
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
