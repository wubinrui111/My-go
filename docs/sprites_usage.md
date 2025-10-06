# 精灵表索引系统使用指南

## 概述

精灵表索引系统为游戏中的所有精灵提供了一个集中管理的方案。通过使用常量和辅助函数，开发者可以更容易地引用和管理精灵表中的各个精灵。

## 精灵索引常量

所有精灵索引都定义为常量，以确保类型安全和易于维护：

```go
const (
    // 玩家精灵
    PlayerSprite = iota
    
    // 方块精灵
    StoneBlockSprite
    DirtBlockSprite
    GrassBlockSprite
    WoodBlockSprite
    LeavesBlockSprite
    
    // 物品精灵（与方块精灵相同）
    StoneItemSprite = StoneBlockSprite
    DirtItemSprite  = DirtBlockSprite
    GrassItemSprite = GrassBlockSprite
    WoodItemSprite  = WoodBlockSprite
    LeavesItemSprite = LeavesBlockSprite
)
```

## 使用方法

### 1. 直接使用常量

```go
import "mygo/internal/pkg/entity"

// 在绘制函数中直接使用常量
func drawPlayer(screen *ebiten.Image, x, y float64) {
    drawSprite(screen, x, y, entity.PlayerSprite)
}

func drawStoneBlock(screen *ebiten.Image, x, y float64) {
    drawSprite(screen, x, y, entity.StoneBlockSprite)
}
```

### 2. 使用辅助函数

```go
// 根据方块类型获取精灵索引
blockType := entity.StoneBlock
spriteIndex := entity.GetBlockSpriteIndex(blockType)

// 根据物品类型获取精灵索引
itemType := entity.Wood
spriteIndex := entity.GetItemSpriteIndex(itemType)

// 根据名称获取精灵索引
spriteIndex := entity.GetSpriteIndex("player")

// 根据索引获取精灵名称
spriteName := entity.GetSpriteName(entity.GrassBlockSprite)
```

### 3. 使用精灵信息映射

```go
// 获取精灵信息
spriteInfo := entity.SpriteMap["stone_block"]
index := spriteInfo.Index  // 1
name := spriteInfo.Name    // "Stone Block"
```

## 扩展精灵表

要添加新的精灵，需要执行以下步骤：

1. 在常量定义中添加新的精灵索引
2. 在SpriteMap中添加新的精灵信息
3. 如有必要，添加专用的获取函数

```go
const (
    // ... existing sprites ...
    NewSprite
)

var SpriteMap = map[string]SpriteInfo{
    // ... existing sprites ...
    "new_sprite": {NewSprite, "New Sprite"},
}

func GetNewSpriteIndex() int {
    return NewSprite
}
```

## 最佳实践

1. **使用常量而非硬编码数字**：提高代码可读性和维护性
2. **使用辅助函数**：减少重复代码，提高一致性
3. **保持映射同步**：确保常量和映射中的信息一致
4. **添加注释**：为复杂的精灵关系添加注释说明