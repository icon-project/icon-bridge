package icon

import (
	"fmt"
	"testing"

	"github.com/icon-project/icon-bridge/bmr/common/crypto"
	"github.com/stretchr/testify/require"
)

func TestNextValidatorHash(t *testing.T) {
	raw := HexBytes("0xf86e950038f35eff5e5516b48a713fe3c8031c94124191f09500f526cc053c33a7c3a48b70111834cf3a71609f0c950014d4c29c4bd2bb2cc79f1284d7b6a403ad6a677a950024791b621e1f25bbac71e2bab8294ff38294a2c69500ed5f818ba1486f996b92cf02db32e4920bfc095f")
	data, err := raw.Value()
	require.NoError(t, err, "failed to decode raw")

	rawh := HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b")
	hash, err := rawh.Value()
	require.NoError(t, err, "failed to decode rawh")

	h := crypto.SHA3Sum256(data)
	require.Equal(t, hash, h, "hash should match")

	fmt.Println(NewHexBytes(h))
}
