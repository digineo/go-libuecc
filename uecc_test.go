package libuecc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"testing"
)

func bitIdentical(a, b unpacked) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func bitIdenticalWork(a, b *Point) bool {
	return bitIdentical(a.X, b.X) && bitIdentical(a.Y, b.Y) && bitIdentical(a.Z, b.Z) && bitIdentical(a.T, b.T)
}

func loadUnpacked(fname string) unpacked {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	if len(buf) < 128 {
		panic("incomplete data")
	}

	var val unpacked
	for i := 0; i < 32; i++ {
		val[i] = binary.LittleEndian.Uint32(buf[i*4:])
	}
	return val
}

func loadInt256File(fname string) *Int256 {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	if len(data) < 32 {
		panic("incomplete data")
	}
	return NewInt256(data)
}

func loadInt256Hex(str string) *Int256 {
	data, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	if len(data) < 32 {
		panic("incomplete data")
	}
	return NewInt256(data)
}

func loadPoint(fname string) *Point {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	if len(data) < 512 {
		panic("incomplete data")
	}

	p := &Point{}
	for i := 0; i < 32; i++ {
		p.X[i] = binary.LittleEndian.Uint32(data[(0+i)*4:])
		p.Y[i] = binary.LittleEndian.Uint32(data[(32+i)*4:])
		p.Z[i] = binary.LittleEndian.Uint32(data[(64+i)*4:])
		p.T[i] = binary.LittleEndian.Uint32(data[(96+i)*4:])
	}
	return p
}

func TestPointDouble(t *testing.T) {
	expected := loadPoint("testdata/cases/ecc_point_double")

	actual := PointBaseLegacy().Double()

	if !bitIdenticalWork(expected, actual) {
		t.Errorf(errmsg, expected, actual)
	}
}

func TestPointAdd(t *testing.T) {
	expected := loadPoint("testdata/cases/ecc_point_add")

	actual := pointIdentity.Add(PointBaseLegacy())

	if !bitIdenticalWork(expected, actual) {
		t.Errorf(errmsg, expected, actual)
	}
}

func TestSimpleArithmetics(t *testing.T) {
	z := zero.add(zero)
	if !bitIdentical(z, zero) {
		t.Errorf("0+0 is not 0, got %v", z)
	}

	z = zero.add(one)
	if !bitIdentical(z, one) {
		t.Errorf("0+1 it not 1, got %v", z)
	}

	z = zero.add(one).add(zero)
	if !bitIdentical(z, one) {
		t.Errorf("0+1+0 it not 0, got %v", z)
	}
}

func TestParity(t *testing.T) {
	tt := map[string]struct {
		expected uint32
		eval     func() uint32
	}{
		"parity_zero":    {0, func() uint32 { return zero.squeeze().parity() }},
		"parity_one":     {1, func() uint32 { return one.squeeze().parity() }},
		"parity_minusp":  {0, func() uint32 { return minusp.squeeze().parity() }},
		"parity_sub_0_1": {0, func() uint32 { return zero.sub(one).squeeze().parity() }},
	}
	for name, tc := range tt {
		expected, actual := tc.expected, tc.eval()
		t.Run(name, func(t *testing.T) {
			if actual != expected {
				t.Errorf("expected %d, got %d\n", expected, actual)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	tt := map[string]struct {
		expected bool
		eval     func() bool
	}{
		"equal_0_0":      {true, func() bool { return zero.equals(zero) }},
		"equal_0_1":      {false, func() bool { return zero.equals(one) }},
		"equal_1_0":      {false, func() bool { return one.equals(zero) }},
		"equal_1_1":      {true, func() bool { return one.equals(one) }},
		"equal_0_minusp": {false, func() bool { return zero.equals(minusp) }},
	}
	for name, tc := range tt {
		expected, actual := tc.expected, tc.eval()
		t.Run(name, func(t *testing.T) {
			if actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tt := map[string]struct {
		expected bool
		eval     func() bool
	}{
		"is_zero_0":      {true, func() bool { return zero.isZero() }},
		"is_zerp_1":      {false, func() bool { return one.isZero() }},
		"is_zero_p":      {true, func() bool { return p.isZero() }},
		"is_zero_minusp": {false, func() bool { return minusp.isZero() }},
	}
	for name, tc := range tt {
		expected, actual := tc.expected, tc.eval()
		t.Run(name, func(t *testing.T) {
			if actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		})
	}
}

func TestSelectPoint(t *testing.T) {
	a, b := pointBaseEd25519, pointIdentity

	if w := selectPoint(&a, &b, 0); !bitIdenticalWork(&w, &a) {
		t.Error("selectPoint(a,b,0) did not return a")
	}
	if w := selectPoint(&a, &b, 1); !bitIdenticalWork(&w, &b) {
		t.Error("selectPoint(a,b,1) did not return b")
	}
}

func TestKeyLoading(t *testing.T) {
	testKeys := []string{
		"83369beddca777585167520fb54a7fb059102bf4e0a46dd5fb1c633d83db77a2",
		"b4dbdb0c05dd28204534fa27c5afca4dcda5397d833e3064f7a7281b249dc7c7",
		"346a11a8bd8fcedfcde2e19c996b6e4497d0dafc3f5af7096c915bd0f9fe4fe9",
		"3bac2ada2fbfa1ea75b2cb214490d5d718f1bbe5b226184488c07cf1a551e8d9",
	}

	for i, keyStr := range testKeys {
		t.Run(fmt.Sprintf("key%d", i), func(t *testing.T) {
			k1 := loadInt256Hex(keyStr)

			// sanity check: did the C generator work with the same key?
			k2 := loadInt256File(fmt.Sprintf("testdata/cases/ecc_key_%d", i))
			if bytes.Compare(k1.Bytes(), k2.Bytes()) != 0 {
				t.Fatal("unexpected key material\n")
			}

			actual := k1.LoadPackedLegacy()
			expected := loadPoint(fmt.Sprintf("testdata/cases/ecc_key_unpacked_%d", i))

			if !bitIdenticalWork(actual, expected) {
				t.Errorf(errmsg, expected, actual)
			}
		})

		t.Run(fmt.Sprintf("derive%d", i), func(t *testing.T) {
			k1 := loadInt256Hex(keyStr)

			work := PointBaseLegacy().ScalarMult(k1)
			actual := work.StorePackedLegacy()

			expected := loadInt256File(
				fmt.Sprintf("testdata/cases/ecc_key_derived_public_%d", i))

			if bytes.Compare(actual.Bytes(), expected.Bytes()) != 0 {
				t.Errorf(errmsg, expected, actual)
			}
		})
	}
}

func TestGeneratedData(t *testing.T) {
	tt := map[string]func() unpacked{
		"zero":   func() unpacked { return zero },
		"one":    func() unpacked { return one },
		"minusp": func() unpacked { return minusp },

		"add_0_0":      func() unpacked { return zero.add(zero) },
		"add_0_1":      func() unpacked { return zero.add(one) },
		"add_1_0":      func() unpacked { return one.add(zero) },
		"add_1_1":      func() unpacked { return one.add(one) },
		"add_0_minusp": func() unpacked { return zero.add(minusp) },
		"add_1_minusp": func() unpacked { return one.add(minusp) },

		"sub_0_0":      func() unpacked { return zero.sub(zero) },
		"sub_0_1":      func() unpacked { return zero.sub(one) },
		"sub_1_0":      func() unpacked { return one.sub(zero) },
		"sub_1_1":      func() unpacked { return one.sub(one) },
		"sub_0_minusp": func() unpacked { return zero.sub(minusp) },
		"sub_1_minusp": func() unpacked { return one.sub(minusp) },

		"sub_add_0_0_0": func() unpacked { return zero.add(zero).sub(zero) },
		"sub_add_0_0_1": func() unpacked { return zero.add(zero).sub(one) },
		"sub_add_0_1_0": func() unpacked { return zero.add(one).sub(zero) },
		"sub_add_0_1_1": func() unpacked { return zero.add(one).sub(one) },
		"sub_add_1_0_0": func() unpacked { return one.add(zero).sub(zero) },
		"sub_add_1_0_1": func() unpacked { return one.add(zero).sub(one) },
		"sub_add_1_1_0": func() unpacked { return one.add(one).sub(zero) },
		"sub_add_1_1_1": func() unpacked { return one.add(one).sub(one) },

		"add_sub_0_0_0": func() unpacked { return zero.sub(zero).add(zero) },
		"add_sub_0_0_1": func() unpacked { return zero.sub(zero).add(one) },
		"add_sub_0_1_0": func() unpacked { return zero.sub(one).add(zero) },
		"add_sub_0_1_1": func() unpacked { return zero.sub(one).add(one) },
		"add_sub_1_0_0": func() unpacked { return one.sub(zero).add(zero) },
		"add_sub_1_0_1": func() unpacked { return one.sub(zero).add(one) },
		"add_sub_1_1_0": func() unpacked { return one.sub(one).add(zero) },
		"add_sub_1_1_1": func() unpacked { return one.sub(one).add(one) },

		"squeeze_zero":    func() unpacked { return zero.squeeze() },
		"squeeze_one":     func() unpacked { return one.squeeze() },
		"squeeze_sub_0_1": func() unpacked { return zero.sub(one).squeeze() },

		"freeze_zero":    func() unpacked { return zero.freeze() },
		"freeze_one":     func() unpacked { return one.freeze() },
		"freeze_sub_0_1": func() unpacked { return zero.sub(one).freeze() },

		"mult_0_0":           func() unpacked { return zero.mult(zero) },
		"mult_0_1":           func() unpacked { return zero.mult(one) },
		"mult_1_0":           func() unpacked { return one.mult(zero) },
		"mult_1_1":           func() unpacked { return one.mult(one) },
		"mult_minusp_minusp": func() unpacked { return minusp.mult(minusp) },

		"mult_int_0_0":   func() unpacked { return zero.multInt(0) },
		"mult_int_1_0":   func() unpacked { return one.multInt(0) },
		"mult_int_0_1":   func() unpacked { return zero.multInt(1) },
		"mult_int_1_1":   func() unpacked { return one.multInt(1) },
		"mult_int_0_max": func() unpacked { return zero.multInt(math.MaxUint32) },
		"mult_int_1_max": func() unpacked { return one.multInt(math.MaxUint32) },

		"square_0":      func() unpacked { return zero.square() },
		"square_1":      func() unpacked { return one.square() },
		"square_minusp": func() unpacked { return minusp.square() },

		"select_0_1_0":      func() unpacked { return selectUnpacked(zero, one, 0) },
		"select_0_1_1":      func() unpacked { return selectUnpacked(zero, one, 1) },
		"select_0_minusp_0": func() unpacked { return selectUnpacked(zero, minusp, 0) },
		"select_0_minusp_1": func() unpacked { return selectUnpacked(zero, minusp, 1) },
		"select_1_minusp_0": func() unpacked { return selectUnpacked(one, minusp, 0) },
		"select_1_minusp_1": func() unpacked { return selectUnpacked(one, minusp, 1) },
	}

	for fname, eval := range tt {
		expected := loadUnpacked("testdata/cases/" + fname)
		actual := eval()

		t.Run(fname, func(t *testing.T) {
			if !bitIdentical(actual, expected) {
				t.Errorf(errmsg, expected, actual)
			}
		})
	}
}

const errmsg = `
expected hex: %04[1]x
got      hex: %04[2]x

expected dec: %[1]v
got      dec: %[2]v
`
