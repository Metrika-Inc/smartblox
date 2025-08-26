package smartblox

import (
	"testing"

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

	expBlockBody := []byte(`{"round":10,"txs":[{"sig":"f1af154530f05a3843a14dc1240efc639268fc96f594248237a9ee3b9502b1a5","tx":{"type":"void"}}]}`)
	require.Equal(t, string(expBlockBody), string(gotBlockBody))

	gotBlockBody, err = GetBlock(11)
	require.Nil(t, err)
	require.NotNil(t, gotBlockBody)

	expBlockBody = []byte(`{"round":11,"txs":[{"sig":"3e0e09c84da7fc381b2d858a8f3c5a20d6360eab7e331a045da01dd6c7957dd1","tx":{"amount":1000,"receipient":2349,"sender":1234,"type":"txfer"}},{"sig":"d7b0517a76c7f8a1265e26e9adfb72f852788e82f9eb82241ea1b3a6986e3bb5","tx":{"type":"void"}}]}`)
	require.Equal(t, string(expBlockBody), string(gotBlockBody))
}

func TestGetBlockError(t *testing.T) {
	gotBlockBody, err := GetBlock(-10)
	require.NotNil(t, err)
	require.Nil(t, gotBlockBody)
}

func TestGetStatus(t *testing.T) {
	gotStatusBody, err := GetStatus()
	expStatusBody := []byte(`{"last-round":200}`)
	require.Nil(t, err)
	require.Equal(t, string(expStatusBody), string(gotStatusBody))
}
