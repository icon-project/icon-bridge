package intconv

import (
	"encoding/json"
	"math"
	"math/big"
)

type BigInt struct {
	big.Int
}

func (i BigInt) MarshalJSON() ([]byte, error) {
	if i.Int.Cmp(big.NewInt(math.MaxInt64)) > 0 {
		return json.Marshal(i.String())
	}
	return json.Marshal(i.Int64())
}

func (i *BigInt) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &i.Int)
	if err != nil {
		var str string
		err = json.Unmarshal(b, &str)
		if err != nil {
			return err
		}
		i.Int.SetString(str, 10)
	}
	return nil
}
