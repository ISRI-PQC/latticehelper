package poly

import (
	"testing"
)

func TestPolyQSerialize(t *testing.T) {
	p := NewRandomPolyQ()
	b := p.Serialize()

	n := DeserializePolyQ(b)

	if !n.Equals(p) {
		t.Error("Poly serialization failed")
	}
}

func TestPolyQNeg(t *testing.T) {
	result := NewPolyQFromCoeffs(1, 2, 3, 4).Neg()
	expected := NewPolyQFromCoeffs(-1, -2, -3, -4)
	if !result.Equals(expected) {
		t.Error("Poly negation failed")
	}
}

func TestPolyQAdd(t *testing.T) {
	result := NewPolyQFromCoeffs(1, 2).Add(NewPolyQFromCoeffs(3, 4))
	expected := NewPolyQFromCoeffs(4, 6)
	if !result.Equals(expected) {
		t.Error("Poly addition failed")
	}
}

func TestPolyQSub(t *testing.T) {
	result := NewPolyQFromCoeffs(1, 2).Sub(NewPolyQFromCoeffs(3, 4))
	expected := NewPolyQFromCoeffs(-2, -2)
	if !result.Equals(expected) {
		t.Error("Poly subtraction failed")
	}
}

func TestPolyQMul(t *testing.T) {
	result := NewPolyQFromCoeffs(1, 2, 3, 4).Mul(NewPolyQFromCoeffs(5, 6, 7, 8))
	expected := NewPolyQFromCoeffs(5, 16, 34, 60, 61, 52, 32)
	if !result.Equals(expected) {
		t.Error("Poly multiplication failed")
	}
}

func TestPolyQPow(t *testing.T) {
	result := NewPolyQFromCoeffs(1, 2, 3, 4).Pow(3)
	expected := NewPolyQFromCoeffs(1, 6, 21, 56, 111, 174, 219, 204, 144, 64)
	if !result.Equals(expected) {
		t.Error("Poly power failed")
	}
}

func TestPolyQScale(t *testing.T) {
	result := NewPolyQFromCoeffs(2, 3).ScaleByInt(5)
	expected := NewPolyQFromCoeffs(10, 15)
	if !result.Equals(expected) {
		t.Error("Poly scale failed")
	}
}
