package entity

// 精灵表索引定义
// 精灵表是640x640像素，每个精灵是32x32像素
// 每行可以放20个精灵 (640/32 = 20)

const (
	// 玩家精灵
	PlayerSprite = iota

	// 方块精灵
	StoneBlockSprite
	DirtBlockSprite // 合并了草方块和泥土方块
	WoodBlockSprite
	LeavesBlockSprite

	// 物品精灵（与方块精灵相同）
	StoneItemSprite  = StoneBlockSprite
	DirtItemSprite   = DirtBlockSprite
	GrassItemSprite  = DirtBlockSprite // 草物品使用泥土精灵
	WoodItemSprite   = WoodBlockSprite
	LeavesItemSprite = LeavesBlockSprite

	// TODO: 添加更多精灵索引，如特效、UI元素等
)

// SpriteInfo 精灵信息结构体
type SpriteInfo struct {
	Index int
	Name  string
}

// 所有精灵信息映射
var SpriteMap = map[string]SpriteInfo{
	"player":       {PlayerSprite, "Player"},
	"stone_block":  {StoneBlockSprite, "Stone Block"},
	"dirt_block":   {DirtBlockSprite, "Dirt/Grass Block"}, // 合并了草方块和泥土方块
	"wood_block":   {WoodBlockSprite, "Wood Block"},
	"leaves_block": {LeavesBlockSprite, "Leaves Block"},
	"stone_item":   {StoneItemSprite, "Stone Item"},
	"dirt_item":    {DirtItemSprite, "Dirt Item"},
	"grass_item":   {GrassItemSprite, "Grass Item"},
	"wood_item":    {WoodItemSprite, "Wood Item"},
	"leaves_item":  {LeavesItemSprite, "Leaves Item"},
}

// GetSpriteIndex 根据名称获取精灵索引
func GetSpriteIndex(name string) int {
	if sprite, exists := SpriteMap[name]; exists {
		return sprite.Index
	}
	return 0 // 默认返回玩家精灵索引
}

// GetSpriteName 根据索引获取精灵名称
func GetSpriteName(index int) string {
	// 对于特定索引返回明确的名称
	switch index {
	case PlayerSprite:
		return "Player"
	case StoneBlockSprite:
		return "Stone Block"
	case DirtBlockSprite:
		return "Dirt/Grass Block"
	case WoodBlockSprite:
		return "Wood Block"
	case LeavesBlockSprite:
		return "Leaves Block"
	}
	return "Unknown"
}

// GetBlockSpriteIndex 根据方块类型获取精灵索引
func GetBlockSpriteIndex(blockType BlockType) int {
	// 由于GrassBlock已合并到DirtBlock中，需要特殊处理
	if blockType == StoneBlock {
		return StoneBlockSprite
	} else if blockType == DirtBlock || blockType == GrassBlock {
		return DirtBlockSprite // 草方块使用泥土精灵
	} else if blockType == WoodBlock {
		return WoodBlockSprite
	} else if blockType == LeavesBlock {
		return LeavesBlockSprite
	}
	return StoneBlockSprite
}

// GetItemSpriteIndex 根据物品类型获取精灵索引
func GetItemSpriteIndex(itemType ItemType) int {
	// 由于Grass已合并到Dirt中，需要特殊处理
	if itemType == Stone {
		return StoneItemSprite
	} else if itemType == Dirt || itemType == Grass {
		return GrassItemSprite // 草物品使用泥土精灵
	} else if itemType == Wood {
		return WoodItemSprite
	} else if itemType == Leaves {
		return LeavesItemSprite
	}
	return StoneItemSprite
}
