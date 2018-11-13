package libuecc

// The order of the prime field
//
// The order is 2^252 + 27742317777372353535851937790883648493.
var gfOrder = Int256{
	0xed, 0xd3, 0xf5, 0x5c, 0x1a, 0x63, 0x12, 0x58,
	0xd6, 0x9c, 0xf7, 0xa2, 0xde, 0xf9, 0xde, 0x14,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10,
}

func gfSelect(r, s Int256, b uint32) (out Int256) {
	bminus1 := uint8(b - 1)
	for j := 0; j < 32; j++ {
		t := bminus1 & (r[j] ^ s[j])
		out[j] = s[j] ^ t
	}
	return
}

// GfIsZero checks if an integer is equal to zero (after reduction)
func (in *Int256) GfIsZero() bool {
	r := in.GfReduce()
	var bits uint32
	for i := 0; i < 32; i++ {
		bits |= uint32(r[i])
	}
	return ((bits-1)>>8)&1 == 1
}

// GfAdd adds two integers as Galois field elements
func (in *Int256) GfAdd(o *Int256) *Int256 {
	out := &Int256{}
	nq := 1 - (in[31] >> 4) - (o[31] >> 4)

	u := uint32(0)
	for j := 0; j < 32; j++ {
		u += uint32(in[j]) + uint32(o[j]) +
			uint32(nq)*uint32(gfOrder[j])
		out[j] = uint8(u)
		u >>= 8
	}
	return out
}

// GfSub subtracts two integers as Galois field elements
func (in *Int256) GfSub(o *Int256) *Int256 {
	out := &Int256{}
	nq := 8 - (in[31] >> 4) + (o[31] >> 4)

	u := uint32(0)
	for j := 0; j < 32; j++ {
		u += uint32(in[j] - o[j] + nq*gfOrder[j])
		out[j] = uint8(u)
		u >>= 8
	}
	return out
}

// GfReduce reduces an integer to a unique representation in the range [0,q-1]
func (in *Int256) GfReduce() *Int256 {
	out := reduce(*in)
	o := Int256(out)
	return &o
}

// Reduces an integer to a unique representation in the range [0,q-1]
func reduce(a [32]uint8) [32]uint8 {
	nq := a[31] >> 4

	var u1, u2 uint32
	var out1, out2 Int256

	for j := 0; j < 32; j++ {
		u1 += uint32(a[j]) - uint32(nq)*uint32(gfOrder[j])
		u2 += uint32(a[j]) - uint32((nq-1))*uint32(gfOrder[j])
		out1[j] = uint8(u1)
		out2[j] = uint8(u2)
		if j != 31 {
			u1 >>= 8
			u2 >>= 8
		}
	}

	msb := (u1 >> 31) & 1
	sel := gfSelect(out1, out2, msb)
	return sel
}

func (in *Int256) reduce() {
	*in = reduce(*in)
}

// Montgomery modular multiplication algorithm
func montgomery(a, b Int256) (out Int256) {
	for i := 0; i < 32; i++ {
		u := uint32(out[0]) + uint32(a[i])*uint32(b[0])
		nq := (u * 27) & 255
		u += nq * uint32(gfOrder[0])

		for j := 1; j < 32; j++ {
			u += ((uint32(out[j]) + uint32(a[i])*uint32(b[j])) +
				uint32(nq)*uint32(gfOrder[j])) << 8
			u >>= 8
			out[j-1] = uint8(u)
		}
		out[31] = uint8(u >> 8)
	}
	return out
}

var _2_512modq = Int256{ // 2^512 mod q
	0x01, 0x0f, 0x9c, 0x44, 0xe3, 0x11, 0x06, 0xa4,
	0x47, 0x93, 0x85, 0x68, 0xa7, 0x1b, 0x0e, 0xd0,
	0x65, 0xbe, 0xf5, 0x17, 0xd2, 0x73, 0xec, 0xce,
	0x3d, 0x9a, 0x30, 0x7c, 0x1b, 0x41, 0x99, 0x03,
}

// GfMult multiplies two integers as Galois field elements
func (in *Int256) GfMult(o *Int256) *Int256 {
	b := o
	b.reduce()

	r := montgomery(*in, *b)
	out := montgomery(r, _2_512modq)
	return &out
}

// GfRecip computes the reciprocal of a Galois field element
func (in *Int256) GfRecip() *Int256 {
	var b, r1, r2 Int256

	r1[0] = 1
	a := *in
	a.reduce()

	useR2 := false
	for i := 0; i < 32; i++ {
		var c uint8

		if i == 0 {
			c = gfOrder[0] - 2
		} else {
			c = gfOrder[i]
		}

		for j := uint8(0); j < 8; j += 2 {
			if c&(1<<j) != 0 {
				if useR2 {
					r1 = montgomery(r2, a)
				} else {
					r2 = montgomery(r1, a)
				}
				useR2 = !useR2
			}
			b = montgomery(a, a)

			if c&(2<<j) != 0 {
				if useR2 {
					r1 = montgomery(r2, b)
				} else {
					r2 = montgomery(r1, b)
				}
				useR2 = !useR2
			}
			a = montgomery(b, b)
		}
	}

	out := montgomery(r2, Int256{1})
	return &out
}

// SanitizeSecret Ensures some properties of a Galois field element to
// make it fit for use as a secret key
//
// This sets the 255th bit and clears the 256th and the bottom three
// bits (so the key will be a multiple of 8). See Daniel J. Bernsteins
// paper "Curve25519: new Diffie-Hellman speed records." for the
// rationale of this.
func (in *Int256) SanitizeSecret() *Int256 {
	out := &Int256{}
	for i := 0; i < 32; i++ {
		out[i] = in[i]
	}
	out[0] &= 0xf8
	out[31] &= 0x7f
	out[31] |= 0x40
	return out
}
