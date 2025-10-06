package entity

import (
	"math"
	"math/rand"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ItemSize = 16
	ItemPickupRange = 32.0
	ItemGravity = 0.2
	ItemMaxFallSpeed = 8.0
)

// ItemEntity 表示世界中的掉落物
type ItemEntity struct {
	X, Y        float64     // 位置
	VX, VY      float64     // 速度
	ItemType    ItemType    // 物品类型
	Count       int         // 数量
	Lifetime    int         // 存活时间
	World       interface{  // 世界引用（使用接口避免循环依赖）
		IsBlockAt(x, y int) bool
	}
}

// NewItemEntity 创建新的掉落物
func NewItemEntity(x, y float64, itemType ItemType, count int) *ItemEntity {
	// 给掉落物一个随机的初始速度
	angle := rand.Float64() * 2 * math.Pi
	speed := 2.0 + rand.Float64() * 2.0
	vx := math.Cos(angle) * speed
	vy := math.Sin(angle) * speed - 2 // 向上弹起

	return &ItemEntity{
		X:        x,
		Y:        y,
		VX:       vx,
		VY:       vy,
		ItemType: itemType,
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
	
	// 更新位置
	item.X += item.VX
	item.Y += item.VY
	
	// 检查与方块的碰撞
	if item.World != nil {
		item.checkCollision()
	}
}

// checkCollision 检查与方块的碰撞
func (item *ItemEntity) checkCollision() {
	// 计算掉落物的边界
	left := int(math.Floor((item.X - ItemSize/2) / BlockSize))
	right := int(math.Floor((item.X + ItemSize/2 - 1) / BlockSize))
	top := int(math.Floor((item.Y - ItemSize/2) / BlockSize))
	bottom := int(math.Floor((item.Y + ItemSize/2 - 1) / BlockSize))
	
	// 检查下方是否有方块
	onGround := false
	if item.VY >= 0 {
		for x := left; x <= right; x++ {
			if item.World.IsBlockAt(x, bottom) {
				// 调整掉落物位置到方块上方
				item.Y = float64(bottom*BlockSize - ItemSize/2)
				onGround = true
				break
			}
		}
	}
	
	// 如果在地面上，应用摩擦力
	if onGround {
		item.VY = 0
		item.VX *= 0.8 // 摩擦力
		if math.Abs(item.VX) < 0.1 {
			item.VX = 0
		}
	}
	
	// 检查上方是否有方块
	if item.VY < 0 {
		for x := left; x <= right; x++ {
			if item.World.IsBlockAt(x, top) {
				// 调整掉落物位置到方块下方
				item.Y = float64((top+1)*BlockSize + ItemSize/2)
				item.VY = 0
				break
			}
		}
	}
	
	// 检查水平碰撞
	if item.VX != 0 {
		facingLeft := item.VX < 0
		checkX := left
		if !facingLeft {
			checkX = right
		}
		
		collision := false
		for y := top; y <= bottom; y++ {
			if item.World.IsBlockAt(checkX, y) {
				collision = true
				break
			}
		}
		
		if collision {
			if facingLeft {
				item.X = float64((checkX+1)*BlockSize + ItemSize/2)
			} else {
				item.X = float64(checkX*BlockSize - ItemSize/2)
			}
			item.VX = 0
		}
	}
}

// TryPickup 尝试拾取掉落物
func (item *ItemEntity) TryPickup(playerX, playerY float64) bool {
	// 计算玩家与掉落物之间的距离
	dx := item.X - playerX
	dy := item.Y - playerY
	distance := math.Sqrt(dx*dx + dy*dy)
	
	// 如果在拾取范围内，拾取成功
	if distance <= ItemPickupRange {
		return true
	}
	return false
}

// GetPosition 获取掉落物位置
func (item *ItemEntity) GetPosition() (float64, float64) {
	return item.X, item.Y
}

// GetItemType 获取物品类型
func (item *ItemEntity) GetItemType() ItemType {
	return item.ItemType
}

// GetCount 获取物品数量
func (item *ItemEntity) GetCount() int {
	return item.Count
}

// IsExpired 检查是否已过期
func (item *ItemEntity) IsExpired() bool {
	return item.Lifetime <= 0
}

// SetWorld 设置世界引用
func (item *ItemEntity) SetWorld(world interface {
	IsBlockAt(x, y int) bool
}) {
	item.World = world
}

// Draw 绘制掉落物
func (item *ItemEntity) Draw(screen *ebiten.Image, cameraX, cameraY, offsetX, offsetY float64) {
	// 计算屏幕位置
	screenX := item.X - cameraX + offsetX
	screenY := item.Y - cameraY + offsetY
	
	// 绘制物品（简单的彩色方块）
	itemColor := getItemColor(item.ItemType)
	ebitenutil.DrawRect(screen, screenX-ItemSize/2, screenY-ItemSize/2, ItemSize, ItemSize, itemColor)
}

// getItemColor 根据物品类型获取颜色
func getItemColor(itemType ItemType) color.RGBA {
	// 由于Grass已合并到Dirt中，需要特殊处理
	if itemType == Stone {
		return color.RGBA{128, 128, 128, 255} // 灰色
	} else if itemType == Dirt || itemType == Grass {
		return color.RGBA{150, 100, 50, 255}  // 棕色
	} else if itemType == Wood {
		return color.RGBA{150, 100, 50, 255}  // 棕色
	} else if itemType == Leaves {
		return color.RGBA{30, 120, 30, 255}   // 深绿色
	}
	return color.RGBA{255, 0, 255, 255}   // 品红色（默认）
}