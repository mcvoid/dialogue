package lexeme

import "testing"

func TestCompareLexeme(t *testing.T) {
	a := Item{Type: Number, Val: "5"}
	b := Item{Type: Symbol, Val: "abc"}
	c := Item{Type: Symbol, Val: "def"}
	d := Item{Type: Number, Val: "5"}

	if !a.CompareItem(d) {
		t.Errorf("expected CompareItem to be true, got false")
	}

	if a.CompareItem(b) {
		t.Errorf("expected CompareItem to be false, got true")
	}

	if b.CompareItem(c) {
		t.Errorf("expected CompareItem to be false, got true")
	}
}
