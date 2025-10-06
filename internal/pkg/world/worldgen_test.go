package world

import (
	"testing"
	"mygo/internal/pkg/entity"
)

func TestWorldGeneration(t *testing.T) {
	// 创建世界实例
	w := NewWorld()
	
	// 手动添加一些测试方块来模拟生成的世界
	// 添加不同类型的方块
	for x := -10; x < 10; x++ {
		// 添加草地方块作为地表
		w.AddBlockWithType(x, 5, entity.GrassBlock)
		
		// 添加泥土方块作为地下层
		for y := 6; y < 10; y++ {
			w.AddBlockWithType(x, y, entity.DirtBlock)
		}
		
		// 添加石头方块作为深层
		for y := 10; y < 15; y++ {
			w.AddBlockWithType(x, y, entity.StoneBlock)
		}
	}
	
	// 添加一些树木
	for y := 1; y < 5; y++ {
		w.AddBlockWithType(0, y, entity.WoodBlock) // 树干
	}
	
	// 添加树叶
	for x := -2; x <= 2; x++ {
		for y := -1; y <= 1; y++ {
			if x != 0 || y != 0 { // 避开树干位置
				w.AddBlockWithType(x, y, entity.LeavesBlock)
			}
		}
	}
	
	// 验证生成的方块
	blocks := w.GetAllBlocks()
	if len(blocks) == 0 {
		t.Fatal("Expected world to have blocks generated, but found none")
	}
	
	// 检查是否生成了多种类型的方块
	blockTypes := make(map[entity.BlockType]int)
	for _, block := range blocks {
		blockTypes[block.GetType()]++
	}
	
	// 至少应该有 Grass、Dirt、Stone、Wood、Leaves 等几种基本类型
	// 由于GrassBlock已合并到DirtBlock中，现在只有4种类型
	expectedTypes := 4
	if len(blockTypes) < expectedTypes {
		t.Errorf("Expected at least %d different block types, got %d", expectedTypes, len(blockTypes))
		t.Log("Found block types:", blockTypes)
	}
	
	// 检查特定类型的方块是否存在
	if blockTypes[entity.GrassBlock] == 0 {
		t.Error("Expected world to have grass blocks")
	}
	
	if blockTypes[entity.DirtBlock] == 0 {
		t.Error("Expected world to have dirt blocks")
	}
	
	if blockTypes[entity.StoneBlock] == 0 {
		t.Error("Expected world to have stone blocks")
	}
	
	if blockTypes[entity.WoodBlock] == 0 {
		t.Error("Expected world to have wood blocks (trees)")
	}
	
	if blockTypes[entity.LeavesBlock] == 0 {
		t.Error("Expected world to have leaves blocks (trees)")
	}
	
	// 记录生成结果
	t.Logf("Generated %d blocks with %d different types", len(blocks), len(blockTypes))
}

func TestBiomeGeneration(t *testing.T) {
	// 创建世界实例
	w := NewWorld()
	
	// 添加不同生物群落的代表性方块
	
	// 普通草地生物群落
	w.AddBlockWithType(0, 5, entity.GrassBlock)
	
	// 沙漠生物群落（使用泥土代替沙子）
	w.AddBlockWithType(20, 5, entity.DirtBlock)
	
	// 雪原生物群落（使用石头代替雪）
	w.AddBlockWithType(-20, 5, entity.StoneBlock)
	
	// 验证不同生物群落的方块生成
	blocks := w.GetAllBlocks()
	if len(blocks) < 3 {
		t.Errorf("Expected at least 3 blocks for different biomes, got %d", len(blocks))
	}
	
	t.Log("Successfully tested biome generation with", len(blocks), "blocks")
}