package entity

import (
	"math"
)

const (
	PlayerSize     = 32
	PlayerSpeed    = 4.0
	JumpPower      = 12.0
	Gravity        = 0.5
	AirResistance  = 0.1
	DoubleJumpMax  = 2
	DashDistance   = 15 // 提升冲刺距离
	DashDuration   = 25  // 延长冲刺持续时间
)

// DashTrail 表示冲刺残影
type DashTrail struct {
	X, Y       float64
	Alpha      float64
	Timer      int
	Duration   int
}

type Player struct {
	X, Y          float64
	VX, VY        float64
	OnGround      bool
	DoubleJump    int
	Dashing       bool
	DashTimer     int
	DashDirection int // 0=right, 1=left, 2=up, 3=down
	World         World // 添加对世界的引用，用于碰撞检测
	DashTrails    []DashTrail // 存储冲刺残影
	Inventory     *Inventory  // 玩家物品栏
}

// World 接口定义，用于玩家与世界交互
type World interface {
	IsBlockAt(x, y int) bool
}

func NewPlayer(x, y float64) *Player { 
	inventory := NewInventory()
	
	return &Player{
		X: x, 
		Y: y,
		OnGround: true,
		DoubleJump: DoubleJumpMax,
		DashTrails: make([]DashTrail, 0),
		Inventory: inventory,
	}
}

// SetWorld 设置玩家所在的世界
func (p *Player) SetWorld(world World) {
	p.World = world
}

// GetInventory 获取玩家物品栏
func (p *Player) GetInventory() *Inventory {
	return p.Inventory
}

// Update 更新玩家状态
func (p *Player) Update() {
	// 处理冲刺状态
	if p.Dashing {
		p.updateDash()
	}
	
	// 更新残影
	p.updateDashTrails()

	// 应用重力
	if !p.OnGround && !p.Dashing {
		p.VY += Gravity
	}

	// 应用空气阻力
	if p.VX > 0 {
		p.VX = max(0, p.VX-AirResistance)
	} else if p.VX < 0 {
		p.VX = min(0, p.VX+AirResistance)
	}

	// 更新位置并处理碰撞
	p.updatePosition()
}

// updatePosition 更新位置并处理碰撞
func (p *Player) updatePosition() {
	// 保存原始位置
	oldX, _ := p.X, p.Y
	
	// 先处理水平移动
	p.X += p.VX
	
	// 检查水平碰撞
	if p.World != nil && p.checkHorizontalCollision() {
		p.X = oldX
		p.VX = 0
	}
	
	// 处理垂直移动
	p.Y += p.VY
	
	// 检查垂直碰撞
	onGround := false
	if p.World != nil {
		onGround = p.checkVerticalCollision()
		if onGround {
			p.VY = 0
		}
	}
	
	// 更新OnGround状态
	p.SetOnGround(onGround)
}

// checkHorizontalCollision 检查水平碰撞
func (p *Player) checkHorizontalCollision() bool {
	if p.World == nil {
		return false
	}
	
	// 计算玩家的边界
	left := int(math.Floor((p.X - PlayerSize/2) / BlockSize))
	right := int(math.Floor((p.X + PlayerSize/2 - 1) / BlockSize))
	top := int(math.Floor((p.Y - PlayerSize/2) / BlockSize))
	bottom := int(math.Floor((p.Y + PlayerSize/2 - 1) / BlockSize))
	
	// 检查左边和右边是否有方块
	for y := top; y <= bottom; y++ {
		if p.VX < 0 && p.World.IsBlockAt(left, y) {
			return true
		}
		if p.VX > 0 && p.World.IsBlockAt(right, y) {
			return true
		}
	}
	
	return false
}

// checkVerticalCollision 检查垂直碰撞
func (p *Player) checkVerticalCollision() bool {
	if p.World == nil {
		return false
	}
	
	// 计算玩家的边界
	left := int(math.Floor((p.X - PlayerSize/2) / BlockSize))
	right := int(math.Floor((p.X + PlayerSize/2 - 1) / BlockSize))
	top := int(math.Floor((p.Y - PlayerSize/2) / BlockSize))
	bottom := int(math.Floor((p.Y + PlayerSize/2 - 1) / BlockSize))
	
	// 检查下方是否有方块（着陆）
	if p.VY >= 0 {
		for x := left; x <= right; x++ {
			if p.World.IsBlockAt(x, bottom) {
				// 调整玩家位置到方块上方
				p.Y = float64(bottom*BlockSize - PlayerSize/2)
				return true
			}
		}
	}
	
	// 检查上方是否有方块（撞头）
	if p.VY < 0 {
		for x := left; x <= right; x++ {
			if p.World.IsBlockAt(x, top) {
				// 调整玩家位置到方块下方
				p.Y = float64((top+1)*BlockSize + PlayerSize/2)
				return false // 仍然不着地
			}
		}
	}
	
	return false
}

// checkDashCollision 检查冲刺时的碰撞
func (p *Player) checkDashCollision(oldX, oldY float64) bool {
	if p.World == nil {
		return false
	}
	
	// 计算玩家的边界
	left := int(math.Floor((p.X - PlayerSize/2) / BlockSize))
	right := int(math.Floor((p.X + PlayerSize/2 - 1) / BlockSize))
	top := int(math.Floor((p.Y - PlayerSize/2) / BlockSize))
	bottom := int(math.Floor((p.Y + PlayerSize/2 - 1) / BlockSize))
	
	// 检查所有方向上是否有方块
	for x := left; x <= right; x++ {
		for y := top; y <= bottom; y++ {
			if p.World.IsBlockAt(x, y) {
				return true
			}
		}
	}
	
	return false
}

// MoveHorizontal 控制水平移动
func (p *Player) MoveHorizontal(direction int) {
	// direction: -1=left, 1=right
	p.VX = float64(direction) * PlayerSpeed
}

// Jump 跳跃
func (p *Player) Jump() {
	if p.OnGround {
		p.VY = -JumpPower
		p.OnGround = false
		p.DoubleJump = DoubleJumpMax - 1
	} else if p.DoubleJump > 0 {
		p.VY = -JumpPower
		p.DoubleJump--
	}
}

// Dash 冲刺，根据玩家状态和鼠标位置决定方向
func (p *Player) Dash(mouseX, mouseY float64) {
	if !p.Dashing {
		p.Dashing = true
		p.DashTimer = DashDuration
		
		// 如果玩家正在移动，则朝移动方向冲刺
		if math.Abs(p.VX) > 0.1 {
			if p.VX > 0 {
				p.DashDirection = 0 // 右
			} else {
				p.DashDirection = 1 // 左
			}
		} else {
			// 玩家静止时，朝鼠标方向冲刺
			playerCenterX := p.X
			playerCenterY := p.Y
			
			// 计算从玩家中心到鼠标的向量
			dx := mouseX - playerCenterX
			dy := mouseY - playerCenterY
			
			// 确定主要方向
			if math.Abs(dx) > math.Abs(dy) {
				// 水平方向为主
				if dx > 0 {
					p.DashDirection = 0 // 右
				} else {
					p.DashDirection = 1 // 左
				}
			} else {
				// 垂直方向为主
				if dy > 0 {
					p.DashDirection = 3 // 下
				} else {
					p.DashDirection = 2 // 上
				}
			}
		}
		
		// 冲刺时忽略重力
		p.VY = 0
		
		// 添加初始残影
		p.addDashTrail()
	}
}

// updateDash 更新冲刺状态
func (p *Player) updateDash() {
	p.DashTimer--
	
	// 添加更多残影
	if p.DashTimer%2 == 0 {
		p.addDashTrail()
	}
	
	// 保存当前位置
	oldX, oldY := p.X, p.Y
	
	// 根据方向移动
	switch p.DashDirection {
	case 0: // Right
		p.X += DashDistance
	case 1: // Left
		p.X -= DashDistance
	case 2: // Up
		p.Y -= DashDistance
	case 3: // Down
		p.Y += DashDistance
	}
	
	// 检查碰撞，如果碰撞则停止冲刺
	if p.checkDashCollision(oldX, oldY) {
		p.X, p.Y = oldX, oldY
		p.Dashing = false
		return
	}
	
	if p.DashTimer <= 0 {
		p.Dashing = false
	}
}

// addDashTrail 添加冲刺残影
func (p *Player) addDashTrail() {
	trail := DashTrail{
		X: p.X,
		Y: p.Y,
		Alpha: 1.0,
		Timer: 0,
		Duration: 20,
	}
	p.DashTrails = append(p.DashTrails, trail)
}

// updateDashTrails 更新残影
func (p *Player) updateDashTrails() {
	for i := len(p.DashTrails) - 1; i >= 0; i-- {
		trail := &p.DashTrails[i]
		trail.Timer++
		
		// 计算透明度
		trail.Alpha = 1.0 - float64(trail.Timer)/float64(trail.Duration)
		
		// 移除过期的残影
		if trail.Timer >= trail.Duration {
			// 从切片中移除元素
			p.DashTrails = append(p.DashTrails[:i], p.DashTrails[i+1:]...)
		}
	}
}

// GetDashTrails 获取冲刺残影用于渲染
func (p *Player) GetDashTrails() []DashTrail {
	return p.DashTrails
}

// SetPosition 设置玩家位置
func (p *Player) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
}

// GetPosition 获取玩家位置
func (p *Player) GetPosition() (float64, float64) {
	return p.X, p.Y
}

// IsOnGround 检查玩家是否在地面上
func (p *Player) IsOnGround() bool {
	return p.OnGround
}

// SetOnGround 设置玩家地面状态
func (p *Player) SetOnGround(onGround bool) {
	// 当从空中落到地面时，重置双跳次数
	if !p.OnGround && onGround {
		p.DoubleJump = DoubleJumpMax
	}
	p.OnGround = onGround
	if onGround {
		p.VY = 0
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}