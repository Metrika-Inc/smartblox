package blckmock

import (
	"encoding/json"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
)

var mockTimeNow = func() int64 {
	return genesis + 1000
}

func TestMain(t *testing.M) {
	timeNowWas := timeNow
	timeNow = mockTimeNow
	defer func() { timeNow = timeNowWas }()

	t.Run()
}

func TestGetBlock(t *testing.T) {
	gotBlockBody, err := GetBlock(10)
	require.Nil(t, err)
	require.NotNil(t, gotBlockBody)

	var expBlock map[string]interface{}
	copier.Copy(&expBlock, defaultBlck)
	expBlock["round"] = 10
	expBlockBody, _ := json.Marshal(expBlock)
	require.Equal(t, string(expBlockBody), string(gotBlockBody))
}

func TestGetBlockError(t *testing.T) {
	gotBlockBody, err := GetBlock(-10)
	require.NotNil(t, err)
	require.Nil(t, gotBlockBody)
}

func TestGetStatus(t *testing.T) {
	gotStatus, err := GetStatus()
	require.Nil(t, err)
	require.Equal(t, int64(200), gotStatus)
}
