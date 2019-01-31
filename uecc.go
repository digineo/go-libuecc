package libuecc

type unpacked [32]uint32

// Int256 represents a 256 bit integer.
//
// All functions of libuecc treat Int256 as unsigned little-endian.
type Int256 [32]uint8

// NewInt256 constructs a new Int256 from the given point
func NewInt256(p []uint8) *Int256 {
	if len(p) < 32 {
		return nil
	}
	out := Int256{}
	copy(out[:], p)
	return &out
}

// Bytes serializes an 256 bit integer
func (in *Int256) Bytes() []uint8 {
	out := in
	return out[:]
}

// Point is a point on the curve unpacked for efficient calculation.
//
// The internal representation of an unpacked Point isn't unique, so for
// serialization it should always be packed.
type Point struct {
	X, Y, Z, T unpacked
}

// pointIdentity is the identity element
var pointIdentity = Point{
	X: unpacked{0},
	Y: unpacked{1},
	Z: unpacked{1},
	T: unpacked{0},
}

// pointBaseEd25519 is the generator point used by Ed25519. It is the
// inverse of pointBaseLegacy.
//
// The order of the base point is 2^252 + 27742317777372353535851937790883648493.
var pointBaseEd25519 = Point{
	X: unpacked{
		0x1a, 0xd5, 0x25, 0x8f, 0x60, 0x2d, 0x56, 0xc9,
		0xb2, 0xa7, 0x25, 0x95, 0x60, 0xc7, 0x2c, 0x69,
		0x5c, 0xdc, 0xd6, 0xfd, 0x31, 0xe2, 0xa4, 0xc0,
		0xfe, 0x53, 0x6e, 0xcd, 0xd3, 0x36, 0x69, 0x21,
	},
	Y: unpacked{
		0x58, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
	},
	Z: unpacked{1},
	T: unpacked{
		0xa3, 0xdd, 0xb7, 0xa5, 0xb3, 0x8a, 0xde, 0x6d,
		0xf5, 0x52, 0x51, 0x77, 0x80, 0x9f, 0xf0, 0x20,
		0x7d, 0xe3, 0xab, 0x64, 0x8e, 0x4e, 0xea, 0x66,
		0x65, 0x76, 0x8b, 0xd7, 0x0f, 0x5f, 0x87, 0x67,
	},
}

// pointBaseLegacy is the ec25519 legacy generator point. It is the
// inverse of pointBaseEd25519.
//
// In new software that doesn't need compablity with older libuecc
// releases, you should consider pointBaseEd25519 instead.
//
// The order of the base point is 2^252 + 27742317777372353535851937790883648493.
var pointBaseLegacy = Point{
	X: unpacked{
		0x1a, 0xd5, 0x25, 0x8f, 0x60, 0x2d, 0x56, 0xc9,
		0xb2, 0xa7, 0x25, 0x95, 0x60, 0xc7, 0x2c, 0x69,
		0x5c, 0xdc, 0xd6, 0xfd, 0x31, 0xe2, 0xa4, 0xc0,
		0xfe, 0x53, 0x6e, 0xcd, 0xd3, 0x36, 0x69, 0x21,
	},
	Y: unpacked{
		0x58, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
	},
	Z: unpacked{1},
	T: unpacked{
		0xa3, 0xdd, 0xb7, 0xa5, 0xb3, 0x8a, 0xde, 0x6d,
		0xf5, 0x52, 0x51, 0x77, 0x80, 0x9f, 0xf0, 0x20,
		0x7d, 0xe3, 0xab, 0x64, 0x8e, 0x4e, 0xea, 0x66,
		0x65, 0x76, 0x8b, 0xd7, 0x0f, 0x5f, 0x87, 0x67,
	},
}

// PointBaseEd25519 returns a pointer to a copy of the generator point
// used by Ed25519. It is the inverse of PointBaseLegacy()
//
// The order of the base point is 2^252 + 27742317777372353535851937790883648493.
func PointBaseEd25519() *Point {
	p := pointBaseEd25519
	return &p
}

// PointBaseLegacy returns a pointer to a copy of the ec25519 legacy
// generator point. It is the inverse of PointBaseEd25519()
//
// In new software that doesn't need compablity with older libuecc
// releases, you should consider pointBaseEd25519 instead.
//
// The order of the base point is 2^252 + 27742317777372353535851937790883648493.
func PointBaseLegacy() *Point {
	p := pointBaseLegacy
	return &p
}

var zero = unpacked{0}
var one = unpacked{1}

var minusp = unpacked{
	19, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 128,
}

// Factor to multiply the X coordinate with to convert from the legacy
// to the Ed25519 curve
var legacyToEd25519 = unpacked{
	0xe7, 0x81, 0xba, 0x00, 0x55, 0xfb, 0x91, 0x33,
	0x7d, 0xe5, 0x82, 0xb4, 0x2e, 0x2c, 0x5e, 0x3a,
	0x81, 0xb0, 0x03, 0xfc, 0x23, 0xf7, 0x84, 0x2d,
	0x44, 0xf9, 0x5f, 0x9f, 0x0b, 0x12, 0xd9, 0x70,
}

// Factor to multiply the X coordinate with to convert from the Ed25519
// to the legacy curve
var ed25519ToLegacy = unpacked{
	0xe9, 0x68, 0x42, 0xdb, 0xaf, 0x04, 0xb4, 0x40,
	0xa1, 0xd5, 0x43, 0xf2, 0xf9, 0x38, 0x31, 0x28,
	0x01, 0x17, 0x05, 0x67, 0x9b, 0x81, 0x61, 0xf8,
	0xa9, 0x5b, 0x3e, 0x6a, 0x20, 0x67, 0x4b, 0x24,
}

// Adds two unpacked integers (modulo p)
func (a unpacked) add(b unpacked) (out unpacked) {
	var u uint32

	for j := 0; j < 31; j++ {
		u += a[j] + b[j]
		out[j] = u & 255
		u >>= 8
	}

	u += a[31] + b[31]
	out[31] = u
	return
}

// Subtracts two unpacked integers (modulo p)
//
// b must be squeezed.
func (a unpacked) sub(b unpacked) (out unpacked) {
	u := uint32(218)

	for j := 0; j < 31; j++ {
		u += a[j] + uint32(0xff00) - b[j]
		out[j] = u & 0xff
		u >>= 8
	}

	u += a[31] - b[31]
	out[31] = u
	return
}

// Performs carry and reduce on an unpacked integer
//
// The result is not always fully reduced, but it will be significantly
// smaller than 2p.
func (a unpacked) squeeze() (out unpacked) {
	u := uint32(0)

	for j := 0; j < 31; j++ {
		u += a[j]
		out[j] = u & 0xff
		u >>= 8
	}

	u += a[31]
	out[31] = u & 127
	u = 19 * (u >> 7)

	for j := 0; j < 31; j++ {
		u += out[j]
		out[j] = u & 0xff
		u >>= 8
	}

	u += out[31]
	out[31] = u
	return
}

// Ensures that the output of a previous squeeze is fully reduced
//
// After a freeze, only the lower byte of each integer part holds a
// meaningful value.
func (a unpacked) freeze() (out unpacked) {
	aorig := a

	out = a.add(minusp)
	negative := -((out[31] >> 7) & 1)

	for j := 0; j < 32; j++ {
		out[j] ^= negative & (aorig[j] ^ out[j])
	}
	return
}

// Returns the parity (lowest bit of the fully reduced value) of a
//
// The input must be squeezed.
func (a unpacked) parity() uint32 {
	b := a.add(minusp)
	return (a[0] ^ (b[31] >> 7) ^ 1) & 1
}

// Multiplies two unpacked integers (modulo p)
//
// The result will be squeezed.
func (a unpacked) mult(b unpacked) (out unpacked) {
	for i := 0; i < 32; i++ {
		u := uint32(0)

		for j := 0; j <= i; j++ {
			u += a[j] * b[i-j]
		}
		for j := i + 1; j < 32; j++ {
			u += 38 * a[j] * b[i+32-j]
		}
		out[i] = u
	}
	return out.squeeze()
}

// Multiplies an unpacked integer with a small integer (modulo p)
//
// The result will be squeezed.
func (a unpacked) multInt(n uint32) (out unpacked) {
	u := uint32(0)

	for j := 0; j < 31; j++ {
		u += n * a[j]
		out[j] = u & 255
		u >>= 8
	}

	u += n * a[31]
	out[31] = u & 127
	u = 19 * (u >> 7)

	for j := 0; j < 31; j++ {
		u += out[j]
		out[j] = u & 255
		u >>= 8
	}

	u += out[31]
	out[31] = u
	return
}

// Squares an unpacked integer
//
// The result will be sqeezed.
func (a unpacked) square() (out unpacked) {
	for i := 0; i < 32; i++ {
		u := uint32(0)

		for j := 0; j < i-j; j++ {
			u += a[j] * a[i-j]
		}
		for j := i + 1; j < i+32-j; j++ {
			u += 38 * a[j] * a[i+32-j]
		}

		u *= 2
		if (i & 1) == 0 {
			u += a[i/2] * a[i/2]
			u += 38 * a[i/2+16] * a[i/2+16]
		}

		out[i] = u
	}
	return out.squeeze()
}

// Checks for the equality of two unpacked integers
func (a unpacked) equals(b unpacked) bool {
	differentbits := uint32(0)

	for i := 0; i < 32; i++ {
		differentbits |= ((a[i] ^ b[i]) & 0xffff)
		differentbits |= ((a[i] ^ b[i]) >> 16)
	}

	return (1 & ((differentbits - 1) >> 16)) == 1
}

var p = unpacked{
	0xed, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
}

func (a unpacked) isZero() bool {
	return a.equals(zero) || a.equals(p)
}

// Copies r to out when b == 0, s when b == 1
func selectPoint(r, s *Point, b uint32) (out Point) {
	bminus1 := b - 1
	for j := 0; j < 32; j++ {
		t := bminus1 & (r.X[j] ^ s.X[j])
		out.X[j] = s.X[j] ^ t

		t = bminus1 & (r.Y[j] ^ s.Y[j])
		out.Y[j] = s.Y[j] ^ t

		t = bminus1 & (r.Z[j] ^ s.Z[j])
		out.Z[j] = s.Z[j] ^ t

		t = bminus1 & (r.T[j] ^ s.T[j])
		out.T[j] = s.T[j] ^ t
	}
	return
}

// Copies r to out when b == 0, s when b == 1
func selectUnpacked(r, s unpacked, b uint32) (out unpacked) {
	bminus1 := b - 1
	for j := 0; j < 32; j++ {
		t := bminus1 & (r[j] ^ s[j])
		out[j] = s[j] ^ t
	}
	return
}

var minus1 = unpacked{
	0xec, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
}

var ρS = unpacked{ // rho_s
	0xb0, 0xa0, 0x0e, 0x4a, 0x27, 0x1b, 0xee, 0xc4,
	0x78, 0xe4, 0x2f, 0xad, 0x06, 0x18, 0x43, 0x2f,
	0xa7, 0xd7, 0xfb, 0x3d, 0x99, 0x00, 0x4d, 0x2b,
	0x0b, 0xdf, 0xc1, 0x4f, 0x80, 0x24, 0x83, 0x2b,
}

// Computes the square root of an unpacked integer (in the prime field modulo p)
//
// If the given integer has no square root, hasRoot=false is returned
func (a unpacked) sqrt() (out unpacked, hasRoot bool) {
	// raise z to power (2^252-2), check if power (2^253-5) equals -1
	a2 := a.square()      // 2
	t1 := a2.square()     // 4
	t0 := t1.square()     // 8
	a9 := t0.mult(a)      // 9
	a11 := a9.mult(a2)    // 11
	t0 = a11.square()     // 22
	a2_5_0 := t0.mult(a9) // 2^5 - 2^0 = 31

	t0 = a2_5_0.square()       // 2^6 - 2^1
	t1 = t0.square()           // 2^7 - 2^2
	t0 = t1.square()           // 2^8 - 2^3
	t1 = t0.square()           // 2^9 - 2^4
	t0 = t1.square()           // 2^10 - 2^5
	a2_10_0 := t0.mult(a2_5_0) // 2^10 - 2^0

	t0 = a2_10_0.square()        // 2^11 - 2^1
	t1 = t0.square()             // 2^12 - 2^2
	for i := 2; i < 10; i += 2 { // 2^20 - 2^10
		t0 = t1.square()
		t1 = t0.square()
	}
	a2_20_0 := t1.mult(a2_10_0) // 2^20 - 2^0

	t0 = a2_20_0.square()        // 2^21 - 2^1
	t1 = t0.square()             // 2^22 - 2^2
	for i := 2; i < 20; i += 2 { // 2^40 - 2^20
		t0 = t1.square()
		t1 = t0.square()
	}
	t0 = t1.mult(a2_20_0) // 2^40 - 2^0

	t1 = t0.square()             // 2^41 - 2^1
	t0 = t1.square()             // 2^42 - 2^2
	for i := 2; i < 10; i += 2 { // 2^50 - 2^10
		t1 = t0.square()
		t0 = t1.square()
	}
	a2_50_0 := t0.mult(a2_10_0) // 2^50 - 2^0

	t0 = a2_50_0.square()        // 2^51 - 2^1
	t1 = t0.square()             // 2^52 - 2^2
	for i := 2; i < 50; i += 2 { // 2^100 - 2^50
		t0 = t1.square()
		t1 = t0.square()
	}
	a2_100_0 := t1.mult(a2_50_0) // 2^100 - 2^0

	t1 = a2_100_0.square()        // 2^101 - 2^1
	t0 = t1.square()              // 2^102 - 2^2
	for i := 2; i < 100; i += 2 { // 2^200 - 2^100
		t1 = t0.square()
		t0 = t1.square()
	}
	t1 = t0.mult(a2_100_0) // 2^200 - 2^0

	t0 = t1.square()             // 2^201 - 2^1
	t1 = t0.square()             // 2^202 - 2^2
	for i := 2; i < 50; i += 2 { // 2^250 - 2^50
		t0 = t1.square()
		t1 = t0.square()
	}
	t0 = t1.mult(a2_50_0) // 2^250 - 2^0

	t1 = t0.square()        // 2^251 - 2^1
	t0 = t1.square()        // 2^252 - 2^2
	a2_252_1 := t0.mult(a2) // 2^252 - 2^1

	t1 = t0.square() // 2^253 - 2^3
	t0 = t1.mult(a2) // 2^253 - 6
	t1 = t0.mult(a)  // 2^253 - 5
	a2_252_1ρS := a2_252_1.mult(ρS)

	if t1.equals(minus1) {
		out = selectUnpacked(a2_252_1, a2_252_1ρS, 1)
	} else {
		out = selectUnpacked(a2_252_1, a2_252_1ρS, 0)
	}

	// Check the root
	hasRoot = out.square().equals(a)
	return
}

// Computes the reciprocal of an unpacked integer (in the prime field modulo p)
func (a unpacked) recip() (out unpacked) {
	a2 := a.square()      // 2
	t1 := a2.square()     // 4
	t0 := t1.square()     // 8
	a9 := t0.mult(a)      // 9
	a11 := a9.mult(a2)    // 11
	t0 = a11.square()     // 22
	a2_5_0 := t0.mult(a9) // 2^5 - 2^0 = 31

	t0 = a2_5_0.square()       // 2^6 - 2^1
	t1 = t0.square()           // 2^7 - 2^2
	t0 = t1.square()           // 2^8 - 2^3
	t1 = t0.square()           // 2^9 - 2^4
	t0 = t1.square()           // 2^10 - 2^5
	a2_10_0 := t0.mult(a2_5_0) // 2^10 - 2^0

	t0 = a2_10_0.square()        // 2^11 - 2^1
	t1 = t0.square()             // 2^12 - 2^2
	for i := 2; i < 10; i += 2 { // 2^20 - 2^10
		t0 = t1.square()
		t1 = t0.square()
	}
	a2_20_0 := t1.mult(a2_10_0) // 2^20 - 2^0

	t0 = a2_20_0.square()        // 2^21 - 2^1
	t1 = t0.square()             // 2^22 - 2^2
	for i := 2; i < 20; i += 2 { // 2^40 - 2^20
		t0 = t1.square()
		t1 = t0.square()
	}
	t0 = t1.mult(a2_20_0) // 2^40 - 2^0

	t1 = t0.square()             // 2^41 - 2^1
	t0 = t1.square()             // 2^42 - 2^2
	for i := 2; i < 10; i += 2 { // 2^50 - 2^10
		t1 = t0.square()
		t0 = t1.square()
	}
	a2_50_0 := t0.mult(a2_10_0) // 2^50 - 2^0

	t0 = a2_50_0.square()        // 2^51 - 2^1
	t1 = t0.square()             // 2^52 - 2^2
	for i := 2; i < 50; i += 2 { // 2^100 - 2^50
		t0 = t1.square()
		t1 = t0.square()
	}
	a2_100_0 := t1.mult(a2_50_0) // 2^100 - 2^0

	t1 = a2_100_0.square()        // 2^101 - 2^1
	t0 = t1.square()              // 2^102 - 2^2
	for i := 2; i < 100; i += 2 { // 2^200 - 2^100
		t1 = t0.square()
		t0 = t1.square()
	}
	t1 = t0.mult(a2_100_0) // 2^200 - 2^0

	t0 = t1.square()             // 2^201 - 2^1
	t1 = t0.square()             // 2^202 - 2^2
	for i := 2; i < 50; i += 2 { // 2^250 - 2^50
		t0 = t1.square()
		t1 = t0.square()
	}
	t0 = t1.mult(a2_50_0) // 2^250 - 2^0

	t1 = t0.square() // 2^251 - 2^1
	t0 = t1.square() // 2^252 - 2^2
	t1 = t0.square() // 2^253 - 2^3
	t0 = t1.square() // 2^254 - 2^4
	t1 = t0.square() // 2^255 - 2^5

	return t1.mult(a11) // 2^255 - 21
}

// Checks if the X and Y coordinates of a work structure represent a valid point of the curve
//
// Also fills out the T coordinate.
func (w *Point) checkLoadXY() (ok bool) {
	/* Check validity */
	x2 := w.X.square()
	y2 := w.Y.square()
	ax2 := x2.multInt(486664)
	dx2 := x2.multInt(486660)
	ax2y2 := ax2.add(y2)
	dx2y2 := dx2.mult(y2)
	r := ax2y2.sub(one.add(dx2y2)).squeeze()

	if ok = r.isZero(); ok {
		w.T = w.X.mult(w.Y)
	}
	return
}

// LoadXYEd25519 loads a point of the Ed25519 curve with given
// coordinates into its unpacked representation
func LoadXYEd25519(x, y *Int256) (out *Point, ok bool) {
	var tmp unpacked
	for i := 0; i < 32; i++ {
		tmp[i] = uint32(x[i])
		out.Y[i] = uint32(y[i])
		if i == 0 {
			out.Z[i] = 1
		}
	}

	out.X = tmp.mult(ed25519ToLegacy)
	if out.checkLoadXY() {
		return out, true
	}
	return nil, false
}

// LoadXYLegacy loads a point of the legacy curve with given coordinates
// into its unpacked representation
func LoadXYLegacy(x, y *Int256) (out *Point, ok bool) {
	var tmp unpacked

	for i := 0; i < 32; i++ {
		tmp[i] = uint32(x[i])
		out.Y[i] = uint32(y[i])
		if i == 0 {
			out.Z[i] = 1
		}
	}

	out.X = tmp.mult(legacyToEd25519)

	if out.checkLoadXY() {
		return out, true
	}
	return nil, false
}

// StoreXYEd25519 stores the x and y coordinates of a point of the Ed25519 curve
func (w *Point) StoreXYEd25519() (x, y *Int256) {
	x, y = &Int256{}, &Int256{}

	Z := w.Z.recip()
	tmp := Z.mult(w.X)
	X := tmp.mult(legacyToEd25519).freeze()
	for i := 0; i < 32; i++ {
		x[i] = uint8(X[i])
	}

	Y := Z.mult(w.Y).freeze()
	for i := 0; i < 32; i++ {
		y[i] = uint8(Y[i])
	}
	return
}

// StoreXYLegacy stores the x and y coordinates of a point of the legacy curve
func (w *Point) StoreXYLegacy() (x, y *Int256) {
	x, y = &Int256{}, &Int256{}

	Z := w.Z.recip()
	X := Z.mult(w.X).mult(ed25519ToLegacy).freeze()
	for i := 0; i < 32; i++ {
		x[i] = uint8(X[i])
	}

	Y := Z.mult(w.Y).freeze()
	for i := 0; i < 32; i++ {
		y[i] = uint8(Y[i])
	}
	return
}

// LoadPackedEd25519 loads a packed point of the Ed25519 curve into its
// unpacked representation
//
// The packed format is different from the legacy one: the legacy format
// contains that X coordinate and the parity of the Y coordinate,
// Ed25519 uses the Y coordinate and the parity of the X coordinate.
func (in *Int256) LoadPackedEd25519() *Point {
	a := unpacked{486664}

	out := &Point{Z: unpacked{1}}
	for i := 0; i < 32; i++ {
		out.Y[i] = uint32(in[i])
	}
	out.Y[31] &= 0x7f

	Y2 := out.Y.square()
	dY2 := Y2.multInt(486660)
	X2 := one.sub(Y2).mult(a.sub(dY2).recip())

	X, ok := X2.sqrt()
	if !ok {
		return nil
	}

	parity := X.mult(legacyToEd25519).parity()

	// No squeeze is necessary after subtractions from zero if the
	// subtrahend is squeezed
	Xt := zero.sub(X)

	out.X = selectUnpacked(X, Xt, uint32((in[31]>>7))^parity)
	out.T = out.X.mult(out.Y)

	return out
}

// LoadPackedLegacy loads a packed point of the legacy curve into its
// unpacked representation
//
// The packed format is different from the Ed25519 one: the legacy
// format contains that X coordinate and the parity of the Y coordinate,
// Ed25519 uses the Y coordinate and the parity of the X coordinate.
func (in *Int256) LoadPackedLegacy() *Point {
	out := &Point{Z: unpacked{1}}
	var xLegacy unpacked
	for i := 0; i < 32; i++ {
		xLegacy[i] = uint32(in[i])
	}
	xLegacy[31] &= 0x7f

	X2 := xLegacy.square()
	aX2 := X2.multInt(486664)
	dX2 := X2.multInt(486660)
	Y2 := one.sub(aX2).mult(one.sub(dX2).recip())

	Y, ok := Y2.sqrt()
	if !ok {
		return nil
	}

	// No squeeze is necessary after subtractions from zero if the
	// subtrahend is squeezed
	Yt := zero.sub(Y)

	out.Y = selectUnpacked(Y, Yt, uint32(in[31]>>7)^Y.parity())
	out.X = xLegacy.mult(legacyToEd25519)
	out.T = out.X.mult(out.Y)

	return out
}

// StorePackedEd25519 stores a point of the Ed25519 curve into its
// packed representation
//
// The packed format is different from the Ed25519 one: the legacy
// format contains that X coordinate and the parity of the Y coordinate,
// Ed25519 uses the Y coordinate and the parity of the X coordinate.
func (w *Point) StorePackedEd25519() *Int256 {
	x, out := w.StoreXYEd25519()
	out[31] |= (x[0] << 7)
	return out
}

// StorePackedLegacy stores a point of the legacy curve into its packed
// representation
//
// The packed format is different from the Ed25519 one: the legacy
// format contains that X coordinate and the parity of the Y coordinate,
// Ed25519 uses the Y coordinate and the parity of the X coordinate.
func (w *Point) StorePackedLegacy() *Int256 {
	out, y := w.StoreXYLegacy()
	out[31] |= (y[0] << 7)
	return out
}

// IsIdentity checks if a point is the identity element of the Elliptic
// Curve group
func (w *Point) IsIdentity() bool {
	yz := w.Y.sub(w.Z).squeeze()
	return w.X.isZero() && yz.isZero()
}

// Negate negates a point of the Elliptic Curve
func (w *Point) Negate() *Point {
	out := &Point{}
	for i := 0; i < 32; i++ {
		out.Y[i] = w.Y[i]
		out.Z[i] = w.Z[i]
	}

	// No squeeze is necessary after subtractions from zero if the
	// subtrahend is squeezed
	out.X = zero.sub(w.X)
	out.T = zero.sub(w.T)
	return out
}

// Double doubles a point of the Elliptic Curve
//
// w.Double() is equivalent to w.add(w), but faster.
func (w *Point) Double() *Point {
	A := w.X.square()

	B := w.Y.square()

	t0 := w.Z.square()
	C := t0.multInt(2)

	D := zero.sub(A)

	t0 = w.X.add(w.Y)
	t1 := t0.square()
	t0 = t1.sub(A)
	E := t0.sub(B)

	G := D.add(B)
	F := G.sub(C)
	H := D.sub(B)

	return &Point{
		X: E.mult(F),
		Y: G.mult(H),
		T: E.mult(H),
		Z: F.mult(G),
	}
}

// Add adds two points of the Elliptic Curve
func (w *Point) Add(o *Point) *Point {
	t0 := w.Y.sub(w.X)
	t1 := t0.multInt(60833)
	t0 = o.Y.sub(o.X)
	A := t0.mult(t1)

	t0 = w.Y.add(w.X)
	t1 = t0.multInt(60833)
	t0 = o.Y.add(o.X)
	B := t0.mult(t1)

	t0 = o.T.multInt(121665)
	C := w.T.mult(t0)

	t0 = o.Z.multInt(2*60833)
	D := w.Z.mult(t0)

	E := B.sub(A)
	F := D.add(C)
	G := D.sub(C)
	H := B.add(A)

	return &Point{
		X: E.mult(F),
		Y: G.mult(H),
		T: E.mult(H),
		Z: F.mult(G),
	}
}

// Sub subtracts two points of the Elliptic Curve
func (w *Point) Sub(o *Point) *Point {
	return w.Add(o.Negate())
}

// ScalarMultBits does a scalar multiplication of a point of the
// Elliptic Curve with an integer of a given bit length
//
// To speed up scalar multiplication when it is known that not the whole
// 256 bits of the scalar are used. The bit length should always be a
// constant and not computed at runtime to ensure that no timing attacks
// are possible.
func (w *Point) ScalarMultBits(n *Int256, bits int) *Point {
	if bits > 256 {
		bits = 256
	}

	cur := pointIdentity
	for pos := bits - 1; pos >= 0; pos-- {
		b := n[pos/8] >> (uint(pos) & 7)
		b &= 1

		q2 := cur.Double()
		q2p := q2.Add(w)
		cur = selectPoint(q2, q2p, uint32(b))
	}
	return &cur
}

// ScalarMult does a scalar multiplication of a point of the Elliptic
// Curve with an integer
func (w *Point) ScalarMult(n *Int256) *Point {
	return w.ScalarMultBits(n, 256)
}
