package poly

import (
	"testing"

	"cyber.ee/muzosh/pq/devkit"
)

func TestMain(m *testing.M) {
	devkit.InitSingle(128, 4294954753)
	m.Run()
}

func TestPolySerialize(t *testing.T) {
	p := NewRandomPoly()
	b := p.Serialize()

	n := DeserializePoly(b)

	if !n.Equals(p) {
		t.Error("Poly serialization failed")
	}
}

func TestPolyNeg(t *testing.T) {
	result := NewPolyFromCoeffs(1, 2, 3, 4).Neg()
	expected := NewPolyFromCoeffs(-1, -2, -3, -4)
	if !result.Equals(expected) {
		t.Error("Poly negation failed")
	}
}

func TestPolyAdd(t *testing.T) {
	result := NewPolyFromCoeffs(1, 2).Add(NewPolyFromCoeffs(3, 4))
	expected := NewPolyFromCoeffs(4, 6)
	if !result.Equals(expected) {
		t.Error("Poly addition failed")
	}
}

func TestPolySub(t *testing.T) {
	result := NewPolyFromCoeffs(1, 2).Sub(NewPolyFromCoeffs(3, 4))
	expected := NewPolyFromCoeffs(-2, -2)
	if !result.Equals(expected) {
		t.Error("Poly subtraction failed")
	}
}

func TestPolyMul(t *testing.T) {
	result := NewPolyFromCoeffs(1, 2, 3, 4).Mul(NewPolyFromCoeffs(5, 6, 7, 8))
	expected := NewPolyFromCoeffs(5, 16, 34, 60, 61, 52, 32)
	if !result.Equals(expected) {
		t.Error("Poly multiplication failed")
	}
}

func TestPolyPow(t *testing.T) {
	result := NewPolyFromCoeffs(1, 2, 3, 4).Pow(3)
	expected := NewPolyFromCoeffs(1, 6, 21, 56, 111, 174, 219, 204, 144, 64)
	if !result.Equals(expected) {
		t.Error("Poly power failed")
	}
}

func TestPolyScale(t *testing.T) {
	result := NewPolyFromCoeffs(2, 3).ScaleByInt(5)
	expected := NewPolyFromCoeffs(10, 15)
	if !result.Equals(expected) {
		t.Error("Poly scale failed")
	}
}
