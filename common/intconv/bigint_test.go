package intconv

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigIntLessThanMaxInt64MarshalJSON(t *testing.T) {
	i := BigInt{}
	i.SetInt64(math.MaxInt64)
	b, err := json.Marshal(i)
	require.NoError(t, err)
	require.Equal(t, "9223372036854775807", string(b))
}

func TestBigIntGreaterThanMaxInt64MarshalJSON(t *testing.T) {
	i := BigInt{}
	i.SetUint64(math.MaxUint64)
	b, err := json.Marshal(i)
	require.NoError(t, err)
	require.Equal(t, `"18446744073709551615"`, string(b))
}
