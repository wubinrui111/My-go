package entity

import (
	"testing"
)

func TestSpriteConstants(t *testing.T) {
	// 测试精灵索引常量是否正确分配
	if PlayerSprite != 0 {
		t.Errorf("PlayerSprite should be 0, got %d", PlayerSprite)
	}
	
	if StoneBlockSprite != 1 {
		t.Errorf("StoneBlockSprite should be 1, got %d", StoneBlockSprite)
	}
	
	if DirtBlockSprite != 2 {
		t.Errorf("DirtBlockSprite should be 2, got %d", DirtBlockSprite)
	}
	
	if WoodBlockSprite != 3 {
		t.Errorf("WoodBlockSprite should be 3, got %d", WoodBlockSprite)
	}
	
	if LeavesBlockSprite != 4 {
		t.Errorf("LeavesBlockSprite should be 4, got %d", LeavesBlockSprite)
	}
	
	// 测试物品精灵索引是否与方块精灵索引相同
	if StoneItemSprite != StoneBlockSprite {
		t.Error("StoneItemSprite should equal StoneBlockSprite")
	}
	
	if DirtItemSprite != DirtBlockSprite {
		t.Error("DirtItemSprite should equal DirtBlockSprite")
	}
	
	if GrassItemSprite != DirtBlockSprite {
		t.Error("GrassItemSprite should equal DirtBlockSprite")
	}
	
	if WoodItemSprite != WoodBlockSprite {
		t.Error("WoodItemSprite should equal WoodBlockSprite")
	}
	
	if LeavesItemSprite != LeavesBlockSprite {
		t.Error("LeavesItemSprite should equal LeavesBlockSprite")
	}
}

func TestSpriteMap(t *testing.T) {
	// 测试SpriteMap是否包含所有必要的精灵
	expectedSprites := []string{
		"player", "stone_block", "dirt_block", "wood_block", "leaves_block",
		"stone_item", "dirt_item", "grass_item", "wood_item", "leaves_item",
	}
	
	for _, spriteName := range expectedSprites {
		if _, exists := SpriteMap[spriteName]; !exists {
			t.Errorf("SpriteMap should contain %s", spriteName)
		}
	}
	
	// 测试SpriteMap中的值是否正确
	if SpriteMap["player"].Index != PlayerSprite {
		t.Error("Player sprite index in SpriteMap is incorrect")
	}
	
	if SpriteMap["stone_block"].Index != StoneBlockSprite {
		t.Error("Stone block sprite index in SpriteMap is incorrect")
	}
	
	if SpriteMap["dirt_block"].Name != "Dirt/Grass Block" {
		t.Error("Dirt block sprite name in SpriteMap is incorrect")
	}
}

func TestGetSpriteIndex(t *testing.T) {
	// 测试通过名称获取精灵索引
	if GetSpriteIndex("player") != PlayerSprite {
		t.Error("GetSpriteIndex failed for player")
	}
	
	if GetSpriteIndex("stone_block") != StoneBlockSprite {
		t.Error("GetSpriteIndex failed for stone_block")
	}
	
	if GetSpriteIndex("unknown") != 0 {
		t.Error("GetSpriteIndex should return 0 for unknown sprites")
	}
}

func TestGetSpriteName(t *testing.T) {
	// 测试通过索引获取精灵名称
	if GetSpriteName(PlayerSprite) != "Player" {
		t.Errorf("GetSpriteName failed for PlayerSprite, got: %s", GetSpriteName(PlayerSprite))
	}
	
	if GetSpriteName(StoneBlockSprite) != "Stone Block" {
		t.Errorf("GetSpriteName failed for StoneBlockSprite, got: %s", GetSpriteName(StoneBlockSprite))
	}
	
	if GetSpriteName(DirtBlockSprite) != "Dirt/Grass Block" {
		t.Errorf("GetSpriteName failed for DirtBlockSprite, got: %s", GetSpriteName(DirtBlockSprite))
	}
	
	if GetSpriteName(999) != "Unknown" {
		t.Errorf("GetSpriteName should return 'Unknown' for unknown indices, got: %s", GetSpriteName(999))
	}
}

func TestGetBlockSpriteIndex(t *testing.T) {
	// 测试通过方块类型获取精灵索引
	if GetBlockSpriteIndex(StoneBlock) != StoneBlockSprite {
		t.Errorf("GetBlockSpriteIndex failed for StoneBlock, got: %d, expected: %d", GetBlockSpriteIndex(StoneBlock), StoneBlockSprite)
	}
	
	if GetBlockSpriteIndex(DirtBlock) != DirtBlockSprite {
		t.Errorf("GetBlockSpriteIndex failed for DirtBlock, got: %d, expected: %d", GetBlockSpriteIndex(DirtBlock), DirtBlockSprite)
	}
	
	// 草方块现在使用泥土精灵
	if GetBlockSpriteIndex(GrassBlock) != DirtBlockSprite {
		t.Errorf("GetBlockSpriteIndex failed for GrassBlock - should use DirtBlockSprite, got: %d, expected: %d", GetBlockSpriteIndex(GrassBlock), DirtBlockSprite)
	}
	
	if GetBlockSpriteIndex(WoodBlock) != WoodBlockSprite {
		t.Errorf("GetBlockSpriteIndex failed for WoodBlock, got: %d, expected: %d", GetBlockSpriteIndex(WoodBlock), WoodBlockSprite)
	}
	
	if GetBlockSpriteIndex(LeavesBlock) != LeavesBlockSprite {
		t.Errorf("GetBlockSpriteIndex failed for LeavesBlock, got: %d, expected: %d", GetBlockSpriteIndex(LeavesBlock), LeavesBlockSprite)
	}
}

func TestGetItemSpriteIndex(t *testing.T) {
	// 测试通过物品类型获取精灵索引
	if GetItemSpriteIndex(Stone) != StoneItemSprite {
		t.Errorf("GetItemSpriteIndex failed for Stone, got: %d, expected: %d", GetItemSpriteIndex(Stone), StoneItemSprite)
	}
	
	if GetItemSpriteIndex(Dirt) != DirtItemSprite {
		t.Errorf("GetItemSpriteIndex failed for Dirt, got: %d, expected: %d", GetItemSpriteIndex(Dirt), DirtItemSprite)
	}
	
	// 草物品现在使用泥土精灵
	if GetItemSpriteIndex(Grass) != GrassItemSprite {
		t.Errorf("GetItemSpriteIndex failed for Grass - should use GrassItemSprite, got: %d, expected: %d", GetItemSpriteIndex(Grass), GrassItemSprite)
	}
	
	if GetItemSpriteIndex(Wood) != WoodItemSprite {
		t.Errorf("GetItemSpriteIndex failed for Wood, got: %d, expected: %d", GetItemSpriteIndex(Wood), WoodItemSprite)
	}
	
	if GetItemSpriteIndex(Leaves) != LeavesItemSprite {
		t.Errorf("GetItemSpriteIndex failed for Leaves, got: %d, expected: %d", GetItemSpriteIndex(Leaves), LeavesItemSprite)
	}
}