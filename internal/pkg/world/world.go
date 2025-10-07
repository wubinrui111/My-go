package world

import (
	"fmt"
	"mygo/internal/pkg/entity"
)

// World represents the game world
type World struct {
	Player *entity.Player
	Blocks map[string]*entity.Block
	Items  []*entity.ItemEntity // 掉落物列表
}

// NewWorld creates a new world
func NewWorld() *World { 
	world := &World{
		Blocks: make(map[string]*entity.Block),
		Items:  make([]*entity.ItemEntity, 0),
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

// RemoveBlock removes a block from the world and creates a drop item
func (w *World) RemoveBlock(x, y int) {
	key := blockKey(x, y)
	block, exists := w.Blocks[key]
	if !exists {
		return
	}
	
	// 创建掉落物（方块的缩影）
	blockX, blockY := block.GetPosition()
	// 在方块的中心位置生成掉落物
	itemX := blockX + float64(entity.BlockSize)/2
	itemY := blockY + float64(entity.BlockSize)/2
	blockType := block.GetType()
	item := entity.NewItemEntityFromBlock(itemX, itemY, blockType, 1)
	item.SetWorld(w)
	w.Items = append(w.Items, item)
	
	// 移除方块
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

// AddItem 添加掉落物到世界
func (w *World) AddItem(item *entity.ItemEntity) {
	item.SetWorld(w)
	w.Items = append(w.Items, item)
}

// GetAllItems 获取所有掉落物
func (w *World) GetAllItems() []*entity.ItemEntity {
	return w.Items
}

// Update 更新世界状态
func (w *World) Update() {
	// 更新所有掉落物
	for i := len(w.Items) - 1; i >= 0; i-- {
		item := w.Items[i]
		item.Update()
		
		// 检查是否可以被玩家拾取
		playerX, playerY := w.Player.GetPosition()
		if item.TryPickup(playerX, playerY) {
			// 拾取物品
			w.Player.GetInventory().AddItem(item.GetItemType(), item.GetCount())
			
			// 从世界中移除掉落物
			w.Items = append(w.Items[:i], w.Items[i+1:]...)
			continue
		}
		
		// 检查是否过期
		if item.IsExpired() {
			// 从世界中移除掉落物
			w.Items = append(w.Items[:i], w.Items[i+1:]...)
		}
	}
}

// blockKey generates a unique key for a block position
func blockKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// getBlockDropItemType 根据方块类型获取掉落物品类型
func getBlockDropItemType(blockType entity.BlockType) entity.ItemType {
	// 由于GrassBlock已合并到DirtBlock中，需要特殊处理
	if blockType == entity.StoneBlock {
		return entity.Stone
	} else if blockType == entity.DirtBlock || blockType == entity.GrassBlock {
		return entity.Dirt // 泥土方块和草方块都掉落泥土物品
	} else if blockType == entity.WoodBlock {
		return entity.Wood
	} else if blockType == entity.LeavesBlock {
		return entity.Leaves
	}
	return entity.Stone // 默认为石头
}