package matrix

import (
	"testing"

	"cyber.ee/pq/devkit"
)

func TestMain(m *testing.M) {
	devkit.InitSingle(128, 4294954753)
	m.Run()
}

func TestPolyMatrixSerialize(t *testing.T) {
	p := NewRandomPolyMatrix(5, 25)
	b := p.Serialize()

	n := DeserializePolyMatrix(b)

	if !n.Equals(p) {
		t.Error("PolyMatrix serialization failed")
	}
}
func TestPolyQMatrixSerialize(t *testing.T) {
	p := NewRandomPolyQMatrix(5, 25)
	b := p.Serialize()

	n := DeserializePolyQMatrix(b)

	if !n.Equals(p) {
		t.Error("PolyQMatrix serialization failed")
	}
}
