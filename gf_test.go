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

func TestInt256_GfSub(t *testing.T) {
	tt := []struct {
		inpA string
		inpB string
		out  string
	}{
		{
			"dacb7c3120e76dd9b3df60b0f20ba1a66d6c906b43757fe332c3234c98b5a50f",
			"c79f728e3a4a80318a7c5853d10580bb6d6c906b43757fe332c3234c98b5a51f",
			"689faee7d21893c0b2e6bc17f5cef7a600000000000000000000000000000080",
		},
		{
			"dacb7c3120e76dd9b3df60b0f20ba1a66d6c906b43757fe332c3234c98b5a50f",
			"44ecfc71fc962fb9730c845b512d037b0ae6f7648b5dd22bcb25643edceb5d73",
			"794ae731af1e5249cf035fe1ac82ae6464869806b817adb7679dbf0dbcc9478c",
		},
		{
			"dacb7c3120e76dd9b3df60b0f20ba1a66d6c906b43757fe332c3234c98b5a50f",
			"844af397faf02aaab7a3dd49f9bab3ca246960ad165addaa2365200a0c96546f",
			"4c18fbae96614400b5cf0d5026fb1e004a0330be2c1ba2380f5e03428c1f5180",
		},
		{
			"dacb7c3120e76dd9b3df60b0f20ba1a66d6c906b43757fe332c3234c98b5a50f",
			"e2b6a1f2f7ad7f4ffdcfd0fd7e80023705801f1f28c24ec90fff6f2b1b60ec2e",
			"3a5c75e02f18a6fa15303c10264e544069ec704c1bb3301a23c4b3207d55b980",
		},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			assert := assert.New(t)

			a, err := hex.DecodeString(tc.inpA)
			require.NoError(t, err)

			b, err := hex.DecodeString(tc.inpB)
			require.NoError(t, err)

			n := NewInt256(a).GfSub(NewInt256(b))

			assert.Equal(tc.out, hex.EncodeToString(n.Bytes()))
		})
	}
}
