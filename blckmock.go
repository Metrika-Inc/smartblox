package blckmock

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/copier"
)

const (
	genesis  = 1752101294
	blckFreq = 5
)

var (
	defaultBlck = map[string]interface{}{
		"txs": []map[string]interface{}{
			{
				"sig": "deadbeef",
				"tx": map[string]interface{}{
					"amount":     1000,
					"sender":     1234,
					"receipient": 2349,
					"type":       "txfer",
				},
			},
			{
				"sig": "badbeef",
				"tx": map[string]interface{}{
					"type": "void",
				},
			},
		},
	}
	lastRound   = int64(-1)
	lastRoundMu = &sync.RWMutex{}
)

var timeNow = func() int64 {
	return time.Now().UTC().Unix()
}

func updateLastRound() {
	lastRoundMu.Lock()
	defer lastRoundMu.Unlock()

	lastRound = (timeNow() - genesis) / blckFreq
}

func GetBlock(round int64) ([]byte, error) {
	updateLastRound()

	if round < 0 {
		return nil, fmt.Errorf("block round must be >= 0")
	}

	lastRoundMu.RLock()
	defer lastRoundMu.RUnlock()

	if round > lastRound {
		return nil, fmt.Errorf("block round must be <= lastRound")
	}

	var blck map[string]interface{}
	err := copier.Copy(&blck, defaultBlck)
	if err != nil {
		return nil, err
	}
	blck["round"] = round

	blckBody, err := json.Marshal(blck)
	if err != nil {
		return nil, err
	}

	return blckBody, nil
}

func GetStatus() (int64, error) {
	updateLastRound()

	return lastRound, nil
}
