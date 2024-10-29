package vector

import (
	"testing"

	"cyber.ee/pq/devkit"
	"cyber.ee/pq/devkit/poly"
)

func TestMain(m *testing.M) {
	devkit.InitSingle(128, 4294954753)
	m.Run()
}

func TestPolyVectorSerialize(t *testing.T) {
	p := NewRandomPolyVector(5)
	b := p.Serialize()

	n := DeserializePolyVector(b)

	if !n.Equals(p) {
		t.Error("PolyVector serialization failed")
	}
}
func TestPolyQVectorSerialize(t *testing.T) {
	p := NewRandomPolyQVector(devkit.DefaultUniformSampler, 5)
	b := p.Serialize()

	n := DeserializePolyQVector(b)

	if !n.Equals(p) {
		t.Error("PolyQVector serialization failed")
	}
}

func TestPolyVectorDotProduct(t *testing.T) {
	base := PolyVector{
		poly.NewPolyFromCoeffs(1, 2, 3),
		poly.NewPolyFromCoeffs(4, 5, 6),
	}
	other := PolyVector{
		poly.NewPolyFromCoeffs(7, 8, 9),
		poly.NewPolyFromCoeffs(10, 11, 12),
	}

	expected := poly.NewPolyFromCoeffs(47, 116, 209, 168, 99)

	result := base.DotProduct(other).(poly.Poly)

	if !expected.Equals(result) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestPolyQVectorDotProduct(t *testing.T) {
	base := PolyQVector{
		poly.NewPolyQFromCoeffs(1, 2, 3),
		poly.NewPolyQFromCoeffs(4, 5, 6),
	}
	other := PolyQVector{
		poly.NewPolyQFromCoeffs(7, 8, 9),
		poly.NewPolyQFromCoeffs(10, 11, 12),
	}

	expected := poly.NewPolyQFromCoeffs(47, 116, 209, 168, 99)

	result := base.DotProduct(other)

	if !expected.Equals(result) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
