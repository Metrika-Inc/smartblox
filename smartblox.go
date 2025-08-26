package smartblox

import (
	"crypto/sha256"
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

	blck, err := generateBlock(round)
	if err != nil {
		return nil, err
	}

	blckBody, err := json.Marshal(blck)
	if err != nil {
		return nil, err
	}

	return blckBody, nil
}

func generateBlock(round int64) (map[string]interface{}, error) {
	var blck map[string]interface{}
	err := copier.Copy(&blck, defaultBlck)
	if err != nil {
		return nil, err
	}
	blck["round"] = round

	// add as many txns as the modulo 10 of the block round
	txnNum := int(round) % 10
	blckTxns := make([]map[string]interface{}, txnNum+1)
	for i := 0; i < txnNum; i++ {
		txnStr := fmt.Sprintf("%d-%d", round, i)
		txnSig := generateTxnHash(txnStr)
		newTxn := map[string]interface{}{
			"sig": txnSig,
			"tx": map[string]interface{}{
				"amount":     1000,
				"sender":     1234,
				"receipient": 2349,
				"type":       "txfer",
			},
		}
		blckTxns[i] = newTxn
	}

	// always add one void txn
	txnStr := fmt.Sprintf("%d-deadbeef", round)
	txnSig := generateTxnHash(txnStr)

	blckTxns[txnNum] = map[string]interface{}{
		"sig": txnSig,
		"tx": map[string]interface{}{
			"type": "void",
		},
	}

	blck["txs"] = blckTxns

	return blck, nil
}

func generateTxnHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetStatus() ([]byte, error) {
	updateLastRound()

	lastRoundResp := map[string]interface{}{
		"last-round": lastRound,
	}

	lastRoundBody, err := json.Marshal(lastRoundResp)
	if err != nil {
		return nil, err
	}
	return lastRoundBody, nil
}
