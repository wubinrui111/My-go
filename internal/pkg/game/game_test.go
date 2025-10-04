package game

import "testing"

func TestNewGame(t *testing.T) {
	g := NewGame()
	if g == nil {
		t.Fatal("NewGame returned nil")
	}
}
