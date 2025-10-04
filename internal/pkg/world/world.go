package world

import (
	"fmt"
	"mygo/internal/pkg/entity"
)

// World represents the game world
type World struct {
	Player *entity.Player
	Blocks map[string]*entity.Block
}

// NewWorld creates a new world
func NewWorld() *World { 
	world := &World{
		Blocks: make(map[string]*entity.Block),
	}
	
	// 创建玩家并设置世界引用
	world.Player = entity.NewPlayer(0, 0)
	world.Player.SetWorld(world)
	
	return world
}

// AddBlock adds a block to the world
func (w *World) AddBlock(x, y int) {
	w.AddBlockWithType(x, y, entity.StoneBlock)
}

// AddBlockWithType 添加指定类型的方块到世界
func (w *World) AddBlockWithType(x, y int, blockType entity.BlockType) {
	key := blockKey(x, y)
	if _, exists := w.Blocks[key]; !exists {
		w.Blocks[key] = entity.NewBlockWithType(x, y, blockType)
	}
}

// RemoveBlock removes a block from the world
func (w *World) RemoveBlock(x, y int) {
	key := blockKey(x, y)
	delete(w.Blocks, key)
}

// GetBlock returns a block at the specified position
func (w *World) GetBlock(x, y int) (*entity.Block, bool) {
	key := blockKey(x, y)
	block, exists := w.Blocks[key]
	return block, exists
}

// GetAllBlocks returns all blocks in the world
func (w *World) GetAllBlocks() []*entity.Block {
	blocks := make([]*entity.Block, 0, len(w.Blocks))
	for _, block := range w.Blocks {
		blocks = append(blocks, block)
	}
	return blocks
}

// IsBlockAt checks if there is a block at the specified grid position
func (w *World) IsBlockAt(x, y int) bool {
	_, exists := w.Blocks[blockKey(x, y)]
	return exists
}

// blockKey generates a unique key for a block position
func blockKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}