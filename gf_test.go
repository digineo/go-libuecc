package libuecc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt256_GfReduce(t *testing.T) {
	tt := []struct {
		inp string
		out string
	}{
		{
			"b7f1ee9373416a49835747455ec4d287bcccc5a4bf8c38156483d46b35ce4dbd",
			"88d65e9551ff9f804d9aa344cd073ea2bbccc5a4bf8c38156483d46b35ce4d0d",
		},
		{
			"f45151f5253c62de69c95935f083b5649876fdb661412d4f32065a7b018bf68b",
			"8cb2a20d5323cf1db7e29c1dfbb4bdbd9776fdb661412d4f32065a7b018bf60b",
		},
		{
			"77f04111cf23a2831ad5ce51903577bff91b281780e445264368d1c78fab157f",
			"fc248986166e211b3e8b09dd79605e2df91b281780e445264368d1c78fab150f",
		},
		{
			"82ce01315f33fac08cf774a8feb1054d933a94dc8aea9f96724ca553557b39a5",
			"4087678f575442502dd7c84a4cef4f7c923a94dc8aea9f96724ca553557b3905",
		},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			assert := assert.New(t)

			b, err := hex.DecodeString(tc.inp)
			require.NoError(t, err)

			n := NewInt256(b).GfReduce()

			assert.Equal(tc.out, hex.EncodeToString(n.Bytes()))
		})
	}
}
