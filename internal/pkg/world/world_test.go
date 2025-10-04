package world

import (
	"testing"
	"mygo/internal/pkg/entity"
)

func TestNewWorld(t *testing.T) {
	world := NewWorld()
	
	if world.Player == nil {
		t.Error("Expected world to have a player")
	}
	
	if world.Blocks == nil {
		t.Error("Expected world to have a blocks map")
	}
	
	if len(world.Blocks) != 0 {
		t.Errorf("Expected empty blocks map, got %d blocks", len(world.Blocks))
	}
}

func TestWorldAddBlock(t *testing.T) {
	world := NewWorld()
	
	// 添加方块
	world.AddBlock(1, 2)
	
	// 检查方块是否存在
	if !world.IsBlockAt(1, 2) {
		t.Error("Expected block at position (1, 2)")
	}
	
	// 检查方块数量
	if len(world.Blocks) != 1 {
		t.Errorf("Expected 1 block, got %d blocks", len(world.Blocks))
	}
	
	// 尝试添加重复方块
	world.AddBlock(1, 2)
	
	// 检查方块数量是否仍为1
	if len(world.Blocks) != 1 {
		t.Errorf("Expected 1 block after adding duplicate, got %d blocks", len(world.Blocks))
	}
}

func TestWorldAddBlockWithType(t *testing.T) {
	world := NewWorld()
	
	// 添加指定类型的方块
	world.AddBlockWithType(1, 2, entity.DirtBlock)
	
	// 检查方块是否存在
	block, exists := world.GetBlock(1, 2)
	if !exists {
		t.Error("Expected block at position (1, 2)")
	}
	
	// 检查方块类型
	if block.GetType() != entity.DirtBlock {
		t.Errorf("Expected DirtBlock, got %v", block.GetType())
	}
}

func TestWorldRemoveBlock(t *testing.T) {
	world := NewWorld()
	
	// 添加方块
	world.AddBlock(1, 2)
	
	// 确保方块存在
	if !world.IsBlockAt(1, 2) {
		t.Error("Expected block at position (1, 2)")
	}
	
	// 移除方块
	world.RemoveBlock(1, 2)
	
	// 检查方块是否已移除
	if world.IsBlockAt(1, 2) {
		t.Error("Expected no block at position (1, 2) after removal")
	}
	
	// 检查方块数量
	if len(world.Blocks) != 0 {
		t.Errorf("Expected 0 blocks after removal, got %d blocks", len(world.Blocks))
	}
}

func TestWorldGetBlock(t *testing.T) {
	world := NewWorld()
	
	// 添加方块
	world.AddBlock(3, 4)
	
	// 获取方块
	block, exists := world.GetBlock(3, 4)
	
	if !exists {
		t.Error("Expected block to exist at position (3, 4)")
	}
	
	if block == nil {
		t.Error("Expected block to be non-nil")
	}
	
	// 尝试获取不存在的方块
	block, exists = world.GetBlock(5, 6)
	
	if exists {
		t.Error("Expected no block at position (5, 6)")
	}
	
	if block != nil {
		t.Error("Expected block to be nil for non-existent block")
	}
}

func TestWorldGetAllBlocks(t *testing.T) {
	world := NewWorld()
	
	// 添加几个方块
	world.AddBlock(1, 1)
	world.AddBlock(2, 2)
	world.AddBlock(3, 3)
	
	// 获取所有方块
	blocks := world.GetAllBlocks()
	
	if len(blocks) != 3 {
		t.Errorf("Expected 3 blocks, got %d blocks", len(blocks))
	}
	
	// 检查所有方块都不为nil
	for i, block := range blocks {
		if block == nil {
			t.Errorf("Expected block %d to be non-nil", i)
		}
	}
}

func TestWorldIsBlockAt(t *testing.T) {
	world := NewWorld()
	
	// 检查不存在的方块
	if world.IsBlockAt(1, 1) {
		t.Error("Expected no block at position (1, 1)")
	}
	
	// 添加方块
	world.AddBlock(1, 1)
	
	// 检查存在的方块
	if !world.IsBlockAt(1, 1) {
		t.Error("Expected block at position (1, 1)")
	}
}

func TestWorldBlockOperationsAtNegativePositions(t *testing.T) {
	world := NewWorld()
	
	// 在负坐标位置添加方块
	world.AddBlock(-1, -2)
	
	// 检查方块是否存在
	if !world.IsBlockAt(-1, -2) {
		t.Error("Expected block at position (-1, -2)")
	}
	
	// 获取方块
	block, exists := world.GetBlock(-1, -2)
	if !exists {
		t.Error("Expected block to exist at position (-1, -2)")
	}
	
	if block == nil {
		t.Error("Expected block to be non-nil")
	}
	
	// 移除方块
	world.RemoveBlock(-1, -2)
	
	// 检查方块是否已移除
	if world.IsBlockAt(-1, -2) {
		t.Error("Expected no block at position (-1, -2) after removal")
	}
}