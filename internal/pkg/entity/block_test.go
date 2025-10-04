package entity

import (
	"testing"
)

func TestNewBlock(t *testing.T) {
	block := NewBlock(1, 2)
	
	if block.X != 32 {
		t.Errorf("Expected X=32, got %f", block.X)
	}
	
	if block.Y != 64 {
		t.Errorf("Expected Y=64, got %f", block.Y)
	}
	
	if block.Type != StoneBlock {
		t.Errorf("Expected Type=StoneBlock, got %v", block.Type)
	}
}

func TestBlockGetPosition(t *testing.T) {
	block := NewBlock(3, 4)
	
	x, y := block.GetPosition()
	
	if x != 96 {
		t.Errorf("Expected x=96, got %f", x)
	}
	
	if y != 128 {
		t.Errorf("Expected y=128, got %f", y)
	}
}

func TestBlockGetGridPosition(t *testing.T) {
	block := NewBlock(5, 6)
	
	x, y := block.GetGridPosition()
	
	if x != 5 {
		t.Errorf("Expected grid x=5, got %d", x)
	}
	
	if y != 6 {
		t.Errorf("Expected grid y=6, got %d", y)
	}
}

func TestNewBlockWithType(t *testing.T) {
	block := NewBlockWithType(1, 2, DirtBlock)
	
	if block.X != 32 {
		t.Errorf("Expected X=32, got %f", block.X)
	}
	
	if block.Y != 64 {
		t.Errorf("Expected Y=64, got %f", block.Y)
	}
	
	if block.Type != DirtBlock {
		t.Errorf("Expected Type=DirtBlock, got %v", block.Type)
	}
}

func TestBlockGetType(t *testing.T) {
	block := NewBlockWithType(0, 0, GrassBlock)
	
	blockType := block.GetType()
	
	if blockType != GrassBlock {
		t.Errorf("Expected Type=GrassBlock, got %v", blockType)
	}
}