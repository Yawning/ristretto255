package edwards25519

import (
	"github.com/gtank/ristretto255/internal/radix51"
	"github.com/gtank/ristretto255/internal/scalar"
	"testing"
	"testing/quick"
)

// quickCheckConfig will make each quickcheck test run (1024 * -quickchecks)
// times. The default value of -quickchecks is 100.
var (
	quickCheckConfig = &quick.Config{MaxCountScale: 1 << 10}
)

func TestScalarMulSmallScalars(t *testing.T) {
	var z scalar.Scalar
	var p, check ProjP3
	p.ScalarMul(&z, &B)
	check.Zero()
	if check.Equal(&p) != 1 {
		t.Error("0*B != 0")
	}

	z = scalar.Scalar([32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	p.ScalarMul(&z, &B)
	check.Set(&B)
	if check.Equal(&p) != 1 {
		t.Error("1*B != 1")
	}
}

func TestScalarMulVsDalek(t *testing.T) {
	expected := ProjP3{
		X: radix51.FieldElement([5]uint64{778774234987948, 1589187156384239, 1213330452914652, 186161118421127, 2186284806803213}),
		Y: radix51.FieldElement([5]uint64{1241255309069369, 1115278942994853, 1016511918109334, 1303231926552315, 1801448517689873}),
		Z: radix51.FieldElement([5]uint64{353337085654440, 1327844406437681, 2207296012811921, 707394926933424, 917408459573183}),
		T: radix51.FieldElement([5]uint64{585487439439725, 1792815221887900, 946062846079052, 1954901232609667, 1418300670001780}),
	}
	z := scalar.Scalar([32]byte{219, 106, 114, 9, 174, 249, 155, 89, 69, 203, 201, 93, 92, 116, 234, 187, 78, 115, 103, 172, 182, 98, 62, 103, 187, 136, 13, 100, 248, 110, 12, 4})

	var p ProjP3
	p.ScalarMul(&z, &B)
	if expected.Equal(&p) != 1 {
		t.Error("Scalar mul does not match dalek")
	}
}

func TestScalarMulDistributesOverAdd(t *testing.T) {
	scalarMulDistributesOverAdd := func(x, y scalar.Scalar) bool {
		// The quickcheck generation strategy chooses a random
		// 32-byte array, but we require that the high bit is
		// unset.  FIXME: make Scalar opaque.  Until then,
		// mask the high bits:
		x[31] &= 127
		y[31] &= 127
		var z scalar.Scalar
		z.Add(&x, &y)
		var p, q, r, check ProjP3
		p.ScalarMul(&x, &B)
		q.ScalarMul(&y, &B)
		r.ScalarMul(&z, &B)
		check.Add(&p, &q)
		return check.Equal(&r) == 1
	}

	if err := quick.Check(scalarMulDistributesOverAdd, quickCheckConfig); err != nil {
		t.Error(err)
	}
}
