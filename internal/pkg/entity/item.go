package entity

import (
	"fmt"
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ItemSize         = 16
	ItemPickupRange  = 32.0
	ItemGravity      = 0.2
	ItemMaxFallSpeed = 8.0
	ItemFriction     = 0.85 // 很大的摩擦力
	ItemBounce       = 0.2  // 弹跳系数
)

// ItemEntity 表示世界中的掉落物
type ItemEntity struct {
	X, Y       float64     // 位置
	VX, VY     float64     // 速度
	ItemType   ItemType    // 物品类型
	BlockType  BlockType   // 方块类型（用于显示方块的缩影）
	Count      int         // 数量
	Lifetime   int         // 存活时间
	World      interface { // 世界引用（使用接口避免循环依赖）
		IsBlockAt(x, y int) bool
	}
}

// NewItemEntity 创建新的掉落物
func NewItemEntity(x, y float64, itemType ItemType, count int) *ItemEntity {
	// 给掉落物一个随机的初始速度
	angle := rand.Float64() * 2 * math.Pi
	speed := 2.0 + rand.Float64()*2.0
	vx := math.Cos(angle) * speed
	vy := math.Sin(angle)*speed - 2 // 向上弹起

	return &ItemEntity{
		X:        x,
		Y:        y,
		VX:       vx,
		VY:       vy,
		ItemType: itemType,
		BlockType: StoneBlock, // 默认为石头方块
		Count:    count,
		Lifetime: 600, // 10秒 (60fps * 10)
	}
}

// NewItemEntityFromBlock 创建来自方块的掉落物
func NewItemEntityFromBlock(x, y float64, blockType BlockType, count int) *ItemEntity {
	// 给掉落物一个随机的初始速度
	angle := rand.Float64() * 2 * math.Pi
	speed := 2.0 + rand.Float64()*2.0
	vx := math.Cos(angle) * speed
	vy := math.Sin(angle)*speed - 2 // 向上弹起

	// 将方块类型转换为物品类型
	itemType := getBlockToItem(blockType)

	return &ItemEntity{
		X:        x,
		Y:        y,
		VX:       vx,
		VY:       vy,
		ItemType: itemType,
		BlockType: blockType,
		Count:    count,
		Lifetime: 600, // 10秒 (60fps * 10)
	}
}

// Update 更新掉落物状态
func (item *ItemEntity) Update() {
	// 减少存活时间
	item.Lifetime--

	// 应用重力
	item.VY += ItemGravity
	if item.VY > ItemMaxFallSpeed {
		item.VY = ItemMaxFallSpeed
	}

	// 应用水平摩擦力（很大的摩擦力）
	item.VX *= ItemFriction

	// 更新位置
	newX := item.X + item.VX
	newY := item.Y + item.VY

	// 碰撞检测和响应
	if item.World != nil {
		// 检查水平碰撞
		if item.checkHorizontalCollision(newX, item.Y) {
			// 水平碰撞，反弹并减速
			item.VX = -item.VX * ItemBounce
			newX = item.X // 回退到碰撞前的位置
		}

		// 检查垂直碰撞
		if item.checkVerticalCollision(item.X, newY) {
			// 垂直碰撞，反弹并减速
			item.VY = -item.VY * ItemBounce
			newY = item.Y // 回退到碰撞前的位置
			
			// 如果速度很小，完全停止
			if math.Abs(item.VY) < 0.5 {
				item.VY = 0
			}
		}
	}

	// 应用新位置
	item.X = newX
	item.Y = newY
}

// checkHorizontalCollision 检查水平碰撞
func (item *ItemEntity) checkHorizontalCollision(x, y float64) bool {
	// 计算物品的边界
	left := int(math.Floor((x - ItemSize/2) / 32))
	right := int(math.Floor((x + ItemSize/2 - 1) / 32))
	top := int(math.Floor((y - ItemSize/2) / 32))
	bottom := int(math.Floor((y + ItemSize/2 - 1) / 32))

	// 检查左边和右边是否有方块
	for gridY := top; gridY <= bottom; gridY++ {
		if item.VX < 0 && item.World.IsBlockAt(left, gridY) {
			return true
		}
		if item.VX > 0 && item.World.IsBlockAt(right, gridY) {
			return true
		}
	}

	return false
}

// checkVerticalCollision 检查垂直碰撞
func (item *ItemEntity) checkVerticalCollision(x, y float64) bool {
	// 计算物品的边界
	left := int(math.Floor((x - ItemSize/2) / 32))
	right := int(math.Floor((x + ItemSize/2 - 1) / 32))
	top := int(math.Floor((y - ItemSize/2) / 32))
	bottom := int(math.Floor((y + ItemSize/2 - 1) / 32))

	// 检查顶部和底部是否有方块
	for gridX := left; gridX <= right; gridX++ {
		if item.VY < 0 && item.World.IsBlockAt(gridX, top) {
			return true
		}
		if item.VY > 0 && item.World.IsBlockAt(gridX, bottom) {
			return true
		}
	}

	return false
}

// Draw 绘制掉落物
func (item *ItemEntity) Draw(screen *ebiten.Image, spriteSheet *ebiten.Image) {
	// 如果是方块类型的掉落物，绘制方块的缩影
	spriteIndex := GetBlockSpriteIndex(item.BlockType)
	
	// 计算精灵在精灵表中的位置
	spriteX := (spriteIndex % 20) * 32 // 每行20个精灵
	spriteY := (spriteIndex / 20) * 32 // 每列20个精灵
	
	// 创建绘制选项
	op := &ebiten.DrawImageOptions{}
	
	// 缩放到物品大小 (16x16)
	op.GeoM.Scale(float64(ItemSize)/32.0, float64(ItemSize)/32.0)
	
	// 移动到物品位置
	op.GeoM.Translate(item.X-float64(ItemSize)/2, item.Y-float64(ItemSize)/2)
	
	// 从精灵表中裁剪出对应的精灵
	sprite := spriteSheet.SubImage(image.Rect(spriteX, spriteY, spriteX+32, spriteY+32)).(*ebiten.Image)
	
	// 绘制精灵
	screen.DrawImage(sprite, op)
	
	// 如果数量大于1，绘制数量
	if item.Count > 1 {
		countStr := fmt.Sprintf("%d", item.Count)
		ebitenutil.DebugPrintAt(screen, countStr, int(item.X)+8, int(item.Y)+8)
	}
}

// DrawWithCamera 使用相机坐标绘制掉落物
func (item *ItemEntity) DrawWithCamera(screen *ebiten.Image, spriteSheet *ebiten.Image, camera *Camera) {
	// 获取屏幕坐标
	screenX, screenY := camera.WorldToScreen(item.X, item.Y)
	
	// 如果是方块类型的掉落物，绘制方块的缩影
	spriteIndex := GetBlockSpriteIndex(item.BlockType)
	
	// 计算精灵在精灵表中的位置
	spriteX := (spriteIndex % 20) * 32 // 每行20个精灵
	spriteY := (spriteIndex / 20) * 32 // 每列20个精灵
	
	// 创建绘制选项
	op := &ebiten.DrawImageOptions{}
	
	// 缩放到物品大小 (16x16)
	op.GeoM.Scale(float64(ItemSize)/32.0, float64(ItemSize)/32.0)
	
	// 移动到物品位置（相对于屏幕）
	op.GeoM.Translate(screenX-float64(ItemSize)/2, screenY-float64(ItemSize)/2)
	
	// 从精灵表中裁剪出对应的精灵
	sprite := spriteSheet.SubImage(image.Rect(spriteX, spriteY, spriteX+32, spriteY+32)).(*ebiten.Image)
	
	// 绘制精灵
	screen.DrawImage(sprite, op)
	
	// 如果数量大于1，绘制数量
	if item.Count > 1 {
		countStr := fmt.Sprintf("%d", item.Count)
		ebitenutil.DebugPrintAt(screen, countStr, int(screenX)+8, int(screenY)+8)
	}
}

// GetPosition 返回掉落物的位置
func (item *ItemEntity) GetPosition() (float64, float64) {
	return item.X, item.Y
}

// GetItemType 返回掉落物的物品类型
func (item *ItemEntity) GetItemType() ItemType {
	return item.ItemType
}

// GetCount 返回掉落物的数量
func (item *ItemEntity) GetCount() int {
	return item.Count
}

// IsExpired 返回掉落物是否过期
func (item *ItemEntity) IsExpired() bool {
	return item.Lifetime <= 0
}

// TryPickup 尝试拾取掉落物
func (item *ItemEntity) TryPickup(playerX, playerY float64) bool {
	// 计算玩家与掉落物之间的距离
	dx := item.X - playerX
	dy := item.Y - playerY
	distance := math.Sqrt(dx*dx + dy*dy)
	
	// 如果距离小于拾取范围，则可以拾取
	return distance <= ItemPickupRange
}

// SetWorld 设置世界引用
func (item *ItemEntity) SetWorld(world interface {
	IsBlockAt(x, y int) bool
}) {
	item.World = world
}
