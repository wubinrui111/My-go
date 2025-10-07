package entity

import (
	"testing"
)

func TestItemEntityLifetime(t *testing.T) {
	// 创建一个掉落物
	initialLifetime := 100
	item := NewItemEntityFromBlock(100, 50, StoneBlock, 1)
	
	// 手动设置生命周期用于测试
	item.Lifetime = initialLifetime
	
	// 检查初始生命周期
	if item.Lifetime != initialLifetime {
		t.Errorf("初始生命周期不正确，期望: %d, 实际: %d", initialLifetime, item.Lifetime)
	}
	
	// 更新一次
	item.Update()
	
	// 检查生命周期是否减少
	if item.Lifetime != initialLifetime-1 {
		t.Errorf("生命周期未正确减少，期望: %d, 实际: %d", initialLifetime-1, item.Lifetime)
	}
}

func TestItemEntityExpiration(t *testing.T) {
	// 创建一个即将过期的掉落物
	item := NewItemEntityFromBlock(100, 50, StoneBlock, 1)
	item.Lifetime = 1
	
	// 检查未过期状态
	if item.IsExpired() {
		t.Error("物品不应在生命周期大于0时过期")
	}
	
	// 更新一次使其过期
	item.Update()
	
	// 检查过期状态
	if !item.IsExpired() {
		t.Error("物品应在生命周期为0或负数时过期")
	}
}

func TestItemEntityPickupRange(t *testing.T) {
	// 创建一个掉落物
	item := NewItemEntityFromBlock(100, 50, StoneBlock, 1)
	
	// 测试在拾取范围内的点
	inRangeX, inRangeY := 105.0, 55.0 // 距离约7.07，小于拾取范围32
	if !item.TryPickup(inRangeX, inRangeY) {
		t.Error("应在拾取范围内可以拾取物品")
	}
	
	// 测试在拾取范围外的点
	outRangeX, outRangeY := 200.0, 150.0 // 距离约141.42，大于拾取范围32
	if item.TryPickup(outRangeX, outRangeY) {
		t.Error("应在拾取范围外无法拾取物品")
	}
}

func TestItemEntityFromBlockCreation(t *testing.T) {
	// 测试不同类型方块创建的掉落物
	blockTypes := []BlockType{StoneBlock, DirtBlock, WoodBlock, LeavesBlock}
	
	for _, blockType := range blockTypes {
		item := NewItemEntityFromBlock(0, 0, blockType, 1)
		
		// 检查方块类型是否正确设置
		if item.BlockType != blockType {
			t.Errorf("方块类型未正确设置，期望: %v, 实际: %v", blockType, item.BlockType)
		}
		
		// 检查物品类型是否正确转换
		expectedItemType := getBlockToItem(blockType)
		if item.ItemType != expectedItemType {
			t.Errorf("物品类型未正确转换，方块: %v, 期望物品: %v, 实际物品: %v", 
				blockType, expectedItemType, item.ItemType)
		}
	}
}

func TestItemEntityCreation(t *testing.T) {
	// 测试直接创建物品实体
	itemTypes := []ItemType{Stone, Dirt, Wood, Leaves}
	
	for _, itemType := range itemTypes {
		item := NewItemEntity(0, 0, itemType, 1)
		
		// 检查物品类型是否正确设置
		if item.ItemType != itemType {
			t.Errorf("物品类型未正确设置，期望: %v, 实际: %v", itemType, item.ItemType)
		}
		
		// 检查数量是否正确设置
		if item.Count != 1 {
			t.Errorf("物品数量未正确设置，期望: %d, 实际: %d", 1, item.Count)
		}
	}
}