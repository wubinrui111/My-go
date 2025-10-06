package game

import (
	"image/color"
	"fmt"
	"image"
	"math"
	"math/rand"
	_ "image/png"

	"mygo/internal/pkg/entity"
	"mygo/internal/pkg/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// PerlinNoise 是一个简单的Perlin噪声生成器
type PerlinNoise struct {
	seed  int64
	p     [512]int
}

// NewPerlinNoise 创建一个新的Perlin噪声生成器
func NewPerlinNoise(seed int64) *PerlinNoise {
	p := &PerlinNoise{seed: seed}
	
	// 初始化置换表
	rand.Seed(seed)
	
	// 创建初始置换表
	var permutation [256]int
	for i := 0; i < 256; i++ {
		permutation[i] = i
	}
	
	// 随机打乱置换表
	for i := 0; i < 256; i++ {
		j := rand.Intn(256)
		permutation[i], permutation[j] = permutation[j], permutation[i]
	}
	
	// 复制置换表两次以避免边界检查
	for i := 0; i < 256; i++ {
		p.p[i] = permutation[i]
		p.p[i+256] = permutation[i]
	}
	
	return p
}

// fade 减缓插值曲线
func (p *PerlinNoise) fade(t float64) float64 {
	return t * t * t * (t * (t * 6 - 15) + 10)
}

// lerp 线性插值
func (p *PerlinNoise) lerp(t, a, b float64) float64 {
	return a + t * (b - a)
}

// grad 计算梯度
func (p *PerlinNoise) grad(hash int, x, y float64) float64 {
	h := hash & 15
	u := x
	if h > 7 {
		u = y
	}
	v := y
	if h > 3 && h != 12 && h != 14 {
		v = x
	}
	return (float64((h&1)<<1) - 1) * u + (float64(((h>>1)&1)<<1) - 1) * v
}

// Noise 生成Perlin噪声值 (-1 to 1)
func (p *PerlinNoise) Noise(x, y float64) float64 {
	// 找到单元格坐标
	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255
	
	// 找到单元格内的坐标
	x -= math.Floor(x)
	y -= math.Floor(y)
	
	// 计算淡入淡出值
	u := p.fade(x)
	v := p.fade(y)
	
	// 获取梯度索引
	A := p.p[X] + Y
	B := p.p[X+1] + Y
	
	// 计算噪声值
	return p.lerp(v, 
		p.lerp(u, p.grad(p.p[A], x, y), p.grad(p.p[B], x-1, y)),
		p.lerp(u, p.grad(p.p[A+1], x, y-1), p.grad(p.p[B+1], x-1, y-1)))
}

// FBM (Fractal Brownian Motion) 分形布朗运动
func (p *PerlinNoise) FBM(x, y, frequency, amplitude float64, octaves int) float64 {
	value := 0.0
	maxValue := 0.0
	
	for i := 0; i < octaves; i++ {
		value += p.Noise(x*frequency, y*frequency) * amplitude
		maxValue += amplitude
		frequency *= 2
		amplitude /= 2
	}
	
	return value / maxValue
}

// GenerateWorldTerrain 生成世界地形
// 使用Perlin噪声生成多样化的地形特征，包括：
// 1. 基础地形（起伏的地面）
// 2. 不同的生物群落（草地、沙漠、雪原）
// 3. 地下矿石层
// 4. 洞穴系统
// 5. 湖泊
// 6. 山脉
// 7. 不同类型的植被（树木、仙人掌等）
func (g *Game) GenerateWorldTerrain() {
	noise := NewPerlinNoise(12345)
	
	// 生成基础地形，范围从-300到300格
	for x := -300; x < 300; x++ {
		// 使用噪声函数生成基础地形高度
		// 通过调整频率参数(0.02)可以控制地形的起伏程度
		terrainNoise := noise.FBM(float64(x)*0.02, 0, 1.0, 1.0, 6)
		groundHeight := int(4 + terrainNoise*12)
		
		// 根据位置生成不同的生物群落
		// 使用低频噪声确定生物群落类型
		biomeNoise := noise.FBM(float64(x)*0.005, 300, 1.0, 1.0, 3)
		
		// 生成地面层（地表和地下几层）
		for y := groundHeight; y < groundHeight+8; y++ {
			blockType := entity.DirtBlock
			if y == groundHeight {
				// 表面是草方块
				blockType = entity.GrassBlock
				// 根据生物群落类型生成不同的表面
				if biomeNoise > 0.5 {
					// 沙漠生物群落 - 使用泥土代替草地
					blockType = entity.DirtBlock
				} else if biomeNoise < -0.5 {
					// 雪原生物群落 - 使用石头
					blockType = entity.StoneBlock
				}
			} else if y > groundHeight+3 {
				// 深层是石头
				blockType = entity.StoneBlock
			}
			g.world.AddBlockWithType(x, y, blockType)
		}
		
		// 生成地下层（石头和矿石）
		// 从地面以下8格开始，一直到y=50
		for y := groundHeight+8; y < 50; y++ {
			// 使用更高频的噪声生成洞穴系统
			caveNoise := noise.FBM(float64(x)*0.05, float64(y)*0.05, 1.0, 1.0, 5)
			
			// 添加洞穴系统 - 如果噪声值小于某个阈值，则不生成方块（形成洞穴）
			if caveNoise > -0.1 {
				// 根据深度和噪声决定方块类型
				if y > groundHeight+25 && noise.Noise(float64(x)*0.1, float64(y)*0.1) > 0.7 {
					// 在较深的地方生成矿石（这里简化为特殊石头）
					g.world.AddBlockWithType(x, y, entity.StoneBlock)
				} else {
					// 生成普通石头
					g.world.AddBlockWithType(x, y, entity.StoneBlock)
				}
			}
		}
		
		// 随机生成树木（在地面上）
		// 每12个单位生成一棵树
		if x%12 == 0 && noise.Noise(float64(x)*0.05, 10) > 0.3 {
			// 树的高度根据噪声值确定
			treeHeight := 4 + int(math.Abs(noise.Noise(float64(x), 20)*4))
			// 生成树干
			for y := groundHeight - treeHeight; y < groundHeight; y++ {
				g.world.AddBlockWithType(x, y, entity.WoodBlock)
			}
			
			// 添加树叶 - 更自然的树冠形状
			for lx := x - 3; lx <= x + 3; lx++ {
				for ly := groundHeight - treeHeight - 4; ly <= groundHeight - treeHeight + 1; ly++ {
					// 使用距离判断生成圆形树冠
					dx := math.Abs(float64(lx - x))
					dy := math.Abs(float64(ly - (groundHeight - treeHeight)))
					distance := math.Sqrt(dx*dx + dy*dy)
					
					if distance <= 3.5 {
						// 检查位置是否已有方块
						if !g.world.IsBlockAt(lx, ly) {
							// 随机决定是否生成树叶，边缘更稀疏
							if noise.Noise(float64(lx)*0.4, float64(ly)*0.4) > -0.3 {
								g.world.AddBlockWithType(lx, ly, entity.LeavesBlock)
							}
						}
					}
				}
			}
		}
		
		// 在特定生物群落生成特殊植物
		if biomeNoise > 0.5 && x%8 == 0 {
			// 沙漠仙人掌
			cactusHeight := 3 + int(math.Abs(noise.Noise(float64(x), 40)*3))
			for y := groundHeight - cactusHeight; y < groundHeight; y++ {
				if !g.world.IsBlockAt(x, y) {
					g.world.AddBlockWithType(x, y, entity.WoodBlock)
				}
			}
		} else if biomeNoise < -0.5 && x%10 == 0 {
			// 雪原云杉树
			treeHeight := 5 + int(math.Abs(noise.Noise(float64(x), 50)*5))
			for y := groundHeight - treeHeight; y < groundHeight; y++ {
				g.world.AddBlockWithType(x, y, entity.WoodBlock)
			}
			
			// 添加针叶树叶
			for ly := groundHeight - treeHeight - 3; ly <= groundHeight - treeHeight + 1; ly++ {
				for lx := x - 2; lx <= x + 2; lx++ {
					if math.Abs(float64(lx-x)) + math.Abs(float64(ly-(groundHeight-treeHeight+1))) <= 2.5 {
						if !g.world.IsBlockAt(lx, ly) {
							g.world.AddBlockWithType(lx, ly, entity.LeavesBlock)
						}
					}
				}
			}
		}
	}
	
	// 生成大型洞穴系统
	// 随机生成5个大型椭圆形洞穴
	for i := 0; i < 5; i++ {
		caveCenterX := rand.Intn(600) - 300
		caveCenterY := rand.Intn(30) + 10
		caveSize := rand.Intn(20) + 10
		
		// 生成椭圆形洞穴
		for x := caveCenterX - caveSize; x <= caveCenterX + caveSize; x++ {
			for y := caveCenterY - caveSize/2; y <= caveCenterY + caveSize/2; y++ {
				// 椭圆方程: (x-h)²/a² + (y-k)²/b² <= 1
				dx := float64(x - caveCenterX)
				dy := float64(y - caveCenterY)
				a := float64(caveSize)
				b := float64(caveSize / 2)
				
				if (dx*dx)/(a*a) + (dy*dy)/(b*b) <= 1.0 {
					key := fmt.Sprintf("%d,%d", x, y)
					delete(g.world.Blocks, key)
				}
			}
		}
	}
	
	// 生成湖泊
	// 使用噪声确定湖泊位置
	for x := -150; x < 150; x++ {
		lakeNoise := noise.FBM(float64(x)*0.04, 100, 1.0, 1.0, 3)
		if lakeNoise < -0.4 {
			// 确定湖泊深度
			lakeDepth := int(2 + math.Abs(lakeNoise)*5)
			
			// 获取地表高度
			terrainNoise := noise.FBM(float64(x)*0.02, 0, 1.0, 1.0, 6)
			groundHeight := int(4 + terrainNoise*12)
			
			// 移除湖泊区域的方块
			for y := groundHeight - lakeDepth; y <= groundHeight; y++ {
				if g.world.IsBlockAt(x, y) {
					key := fmt.Sprintf("%d,%d", x, y)
					delete(g.world.Blocks, key)
				}
			}
			
			// 在湖泊底部添加泥土
			if !g.world.IsBlockAt(x, groundHeight+1) {
				g.world.AddBlockWithType(x, groundHeight+1, entity.DirtBlock)
			}
		}
	}
	
	// 生成山脉
	// 使用低频噪声生成高山
	for x := -200; x < 200; x++ {
		mountainNoise := noise.FBM(float64(x)*0.01, 200, 1.0, 1.0, 4)
		if mountainNoise > 0.6 {
			// 山脉高度
			mountainHeight := int(mountainNoise * 25)
			
			// 获取基础地形高度
			terrainNoise := noise.FBM(float64(x)*0.02, 0, 1.0, 1.0, 6)
			baseHeight := int(4 + terrainNoise*12)
			
			// 生成山脉
			for y := baseHeight - mountainHeight; y < baseHeight; y++ {
				if !g.world.IsBlockAt(x, y) {
					blockType := entity.StoneBlock
					if y == baseHeight - mountainHeight {
						// 山顶可能是雪地或石头
						if noise.Noise(float64(x), float64(y)) > 0.5 {
							blockType = entity.StoneBlock
						} else {
							blockType = entity.StoneBlock
						}
					} else if y > baseHeight - mountainHeight + 15 {
						blockType = entity.StoneBlock
					} else {
						blockType = entity.StoneBlock
					}
					g.world.AddBlockWithType(x, y, blockType)
				}
			}
		}
	}
	
	// 确保玩家出生点附近是安全的，移除周围的方块
	for x := -5; x <= 5; x++ {
		for y := -12; y <= 6; y++ {
			if g.world.IsBlockAt(x, y) {
				// 移除玩家出生点附近的方块
				key := fmt.Sprintf("%d,%d", x, y)
				delete(g.world.Blocks, key)
			}
		}
	}
	
	// 将玩家放置在地面上方
	g.player.SetPosition(0, float64(-10 * entity.BlockSize))
}

func NewGame() *Game {
	// 加载精灵表
	spriteSheet, _, err := ebitenutil.NewImageFromFile("image/test.png")
	if err != nil {
		panic(err)
	}
	
	// 创建世界
	w := world.NewWorld()
	
	// 创建游戏实例
	g := &Game{
		player: w.Player,
		camera: entity.NewCamera(0, 0),
		world:  w,
		lastPlacePos:    [2]int{-1, -1}, // 初始化为无效位置
		lastDestroyPos:  [2]int{-1, -1}, // 初始化为无效位置
		spriteSheet: spriteSheet,
	}
	
	// 生成世界地形
	g.GenerateWorldTerrain()
	
	// 设置相机
	g.camera.SetScreenSize(800, 600)
	
	return g
}

type Game struct {
	player *entity.Player
	camera *entity.Camera
	world  *world.World
	// 连续放置/破坏方块相关变量
	lastPlacePos    [2]int // 记录上次放置方块的网格位置
	lastDestroyPos  [2]int // 记录上次破坏方块的网格位置
	spriteSheet     *ebiten.Image // 精灵表
}

func (g *Game) Update() error {
	// 更新玩家输入
	g.handleInput()
	
	// 更新玩家状态
	g.player.Update()
	
	// 更新世界状态（包括掉落物）
	g.world.Update()
	
	// 更新摄像机跟随
	g.camera.SetTarget(g.player.GetPosition())
	g.camera.Update()
	
	// 处理连续放置方块（仅当物品栏未展开时）
	if !g.player.GetInventory().IsOpen() {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			g.placeBlock()
		}
		
		// 处理连续破坏方块（仅当物品栏未展开时）
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.destroyBlock()
		}
	}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x18, 0x20, 0xff})
	
	// 绘制所有方块
	blocks := g.world.GetAllBlocks()
	for _, block := range blocks {
		blockX, blockY := block.GetPosition()
		screenX, screenY := g.camera.WorldToScreen(blockX, blockY)
		
		// 根据方块类型绘制对应的精灵
		spriteIndex := entity.GetBlockSpriteIndex(block.GetType())
		g.drawSprite(screen, screenX, screenY, spriteIndex)
	}
	
	// 绘制所有掉落物
	items := g.world.GetAllItems()
	for _, item := range items {
		itemX, itemY := item.GetPosition()
		screenX, screenY := g.camera.WorldToScreen(itemX, itemY)
		// 绘制物品（简单的彩色方块）
		itemColor := getItemColor(item.GetItemType())
		ebitenutil.DrawRect(screen, screenX-8, screenY-8, 16, 16, itemColor)
	}
	
	// 绘制冲刺残影
	dashTrails := g.player.GetDashTrails()
	for _, trail := range dashTrails {
		screenX, screenY := g.camera.WorldToScreen(trail.X, trail.Y)
		// 使用半透明红色方块表示残影
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenX-16, screenY-16)
		op.ColorM.Scale(1, 1, 1, trail.Alpha*0.5) // 设置透明度
		// 绘制玩家精灵作为残影
		g.drawSpriteWithOp(screen, op, entity.PlayerSprite)
	}
	
	// 绘制玩家
	playerX, playerY := g.player.GetPosition()
	screenX, screenY := g.camera.WorldToScreen(playerX, playerY)
	// 绘制玩家精灵
	g.drawSprite(screen, screenX-16, screenY-16, entity.PlayerSprite)
	
	// 绘制底部快捷栏
	g.drawHotbar(screen)
	
	// 如果物品栏展开，绘制完整物品栏
	if g.player.GetInventory().IsOpen() {
		g.drawInventory(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) { 
	return 800, 600 
}

// handleInput 处理用户输入
func (g *Game) handleInput() {
	// 只有当物品栏未展开时才处理移动
	if !g.player.GetInventory().IsOpen() {
		// 水平移动
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.player.MoveHorizontal(-1)
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.player.MoveHorizontal(1)
		}
		
		// 跳跃
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.player.Jump()
		}
		
		// 冲刺
		if inpututil.IsKeyJustPressed(ebiten.KeyShift) {
			// 获取鼠标位置
			mx, my := ebiten.CursorPosition()
			// 转换为世界坐标
			worldX, worldY := g.camera.ScreenToWorld(float64(mx), float64(my))
			// 向玩家传递世界坐标进行冲刺
			g.player.Dash(worldX, worldY)
		}
		
		// 单次放置方块（向后兼容）
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			g.placeBlock()
		}
		
		// 单次破坏方块（向后兼容）
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.destroyBlock()
		}
	}
	
	// 处理物品栏相关输入（任何时候都可以）
	// E键切换物品栏展开状态
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.player.GetInventory().ToggleOpen()
	}
	
	// 数字键1-9选择快捷栏槽位
	for i := 0; i < 9; i++ {
		if inpututil.IsKeyJustPressed(ebiten.Key1 + ebiten.Key(i)) {
			g.player.GetInventory().SetSelectedSlot(i)
			break
		}
	}
	
	// 鼠标滚轮切换物品
	_, wheelY := ebiten.Wheel()
	if wheelY > 0 {
		// 向上滚动，选择前一个槽位
		currentSlot := g.player.GetInventory().GetSelectedSlot()
		newSlot := (currentSlot + 9 - 1) % 9 // 加9是为了处理负数情况
		g.player.GetInventory().SetSelectedSlot(newSlot)
	} else if wheelY < 0 {
		// 向下滚动，选择后一个槽位
		currentSlot := g.player.GetInventory().GetSelectedSlot()
		newSlot := (currentSlot + 1) % 9
		g.player.GetInventory().SetSelectedSlot(newSlot)
	}
	
	// ESC键关闭物品栏
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.player.GetInventory().CloseInventory()
	}
}

// drawHotbar 绘制底部快捷栏
func (g *Game) drawHotbar(screen *ebiten.Image) {
	inventory := g.player.GetInventory()
	
	// 快捷栏背景
	hotbarWidth := 9 * 40
	hotbarHeight := 40
	hotbarX := (800 - hotbarWidth) / 2
	hotbarY := 600 - hotbarHeight - 10
	
	// 绘制快捷栏背景
	ebitenutil.DrawRect(screen, float64(hotbarX), float64(hotbarY), float64(hotbarWidth), float64(hotbarHeight), color.RGBA{0, 0, 0, 100})
	
	// 绘制每个槽位
	for i := 0; i < entity.HotbarSlotCount; i++ {
		x := hotbarX + i*40 + 2
		y := hotbarY + 2
		
		// 绘制槽位背景
		slotColor := color.RGBA{50, 50, 50, 200}
		if i == inventory.GetSelectedSlot() {
			// 选中的槽位高亮显示
			slotColor = color.RGBA{100, 100, 100, 200}
		}
		ebitenutil.DrawRect(screen, float64(x), float64(y), 36, 36, slotColor)
		
		// 绘制物品
		item := inventory.GetHotbarSlot(i)
		if item.Type != entity.Air {
			// 根据物品类型绘制对应的精灵
			spriteIndex := getItemToSpriteIndex(item.Type)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(36.0/32.0, 36.0/32.0) // 缩放到36x36
			op.GeoM.Translate(float64(x+4), float64(y+4))
			g.drawSpriteWithOp(screen, op, spriteIndex)
			
			// 绘制数量
			if item.Count > 1 {
				countText := fmt.Sprintf("%d", item.Count)
				ebitenutil.DebugPrintAt(screen, countText, x+25-len(countText)*3, y+20)
			}
		}
		
		// 绘制槽位编号
		slotNumber := fmt.Sprintf("%d", (i+1)%10)
		ebitenutil.DebugPrintAt(screen, slotNumber, x+2, y+2)
	}
}

// drawInventory 绘制完整物品栏
func (g *Game) drawInventory(screen *ebiten.Image) {
	// 绘制半透明背景覆盖整个屏幕
	ebitenutil.DrawRect(screen, 0, 0, 800, 600, color.RGBA{0, 0, 0, 150})
	
	// 物品栏网格 9x3
	inventory := g.player.GetInventory()
	gridWidth := 9 * 40
	gridHeight := 3 * 40
	gridX := (800 - gridWidth) / 2
	gridY := (600 - gridHeight) / 2
	
	// 绘制物品栏标题
	ebitenutil.DebugPrintAt(screen, "Inventory", gridX, gridY-20)
	
	// 绘制物品栏槽位
	for row := 0; row < 3; row++ {
		for col := 0; col < 9; col++ {
			slotIndex := entity.HotbarSlotCount + row*9 + col
			x := gridX + col*40 + 2
			y := gridY + row*40 + 2
			
			// 绘制槽位背景
			ebitenutil.DrawRect(screen, float64(x), float64(y), 36, 36, color.RGBA{50, 50, 50, 200})
			
			// 绘制物品
			item := inventory.GetSlot(slotIndex)
			if item.Type != entity.Air {
				// 根据物品类型绘制对应的精灵
				spriteIndex := getItemToSpriteIndex(item.Type)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(28.0/32.0, 28.0/32.0) // 缩放到28x28
				op.GeoM.Translate(float64(x+4), float64(y+4))
				g.drawSpriteWithOp(screen, op, spriteIndex)
				
				// 绘制数量
				if item.Count > 1 {
					countText := fmt.Sprintf("%d", item.Count)
					ebitenutil.DebugPrintAt(screen, countText, x+25-len(countText)*3, y+20)
				}
			}
		}
	}
	
	// 绘制说明文字
	ebitenutil.DebugPrintAt(screen, "Press ESC or E to close inventory", 280, 580)
}

// drawSprite 绘制精灵
func (g *Game) drawSprite(screen *ebiten.Image, x, y float64, index int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	g.drawSpriteWithOp(screen, op, index)
}

// drawSpriteWithOp 使用指定选项绘制精灵
func (g *Game) drawSpriteWithOp(screen *ebiten.Image, op *ebiten.DrawImageOptions, index int) {
	// 精灵表是640x640，每个精灵是32x32
	// 每行可以放20个精灵 (640/32 = 20)
	spriteRow := index / 20
	spriteCol := index % 20
	
	// 计算源矩形
	srcX := spriteCol * 32
	srcY := spriteRow * 32
	
	// 创建源图像的部分
	srcRect := image.Rect(srcX, srcY, srcX+32, srcY+32)
	
	// 绘制精灵
	screen.DrawImage(g.spriteSheet.SubImage(srcRect).(*ebiten.Image), op)
}

// getItemColor 根据物品类型获取颜色（备用方案）
func getItemColor(itemType entity.ItemType) color.RGBA {
	// 由于Grass已合并到Dirt中，需要特殊处理
	if itemType == entity.Stone {
		return color.RGBA{128, 128, 128, 255} // 灰色
	} else if itemType == entity.Dirt || itemType == entity.Grass {
		return color.RGBA{150, 100, 50, 255}  // 棕色
	} else if itemType == entity.Wood {
		return color.RGBA{150, 100, 50, 255}  // 棕色
	} else if itemType == entity.Leaves {
		return color.RGBA{30, 120, 30, 255}   // 深绿色
	}
	return color.RGBA{255, 0, 255, 255}   // 品红色（默认）
}

// getItemToSpriteIndex 将物品类型转换为精灵表索引
func getItemToSpriteIndex(itemType entity.ItemType) int {
	return entity.GetItemSpriteIndex(itemType)
}

// getBlockColor 根据方块类型获取颜色
func getBlockColor(blockType entity.BlockType) color.RGBA {
	// 由于GrassBlock已合并到DirtBlock中，需要特殊处理
	if blockType == entity.StoneBlock {
		return color.RGBA{128, 128, 128, 255} // 灰色
	} else if blockType == entity.DirtBlock || blockType == entity.GrassBlock {
		return color.RGBA{150, 100, 50, 255}  // 棕色
	} else if blockType == entity.WoodBlock {
		return color.RGBA{150, 100, 50, 255}  // 棕色
	} else if blockType == entity.LeavesBlock {
		return color.RGBA{30, 120, 30, 255}   // 深绿色
	}
	return color.RGBA{139, 69, 19, 255}   // 棕色（默认）
}

// getItemToBlockType 将物品类型转换为方块类型
func getItemToBlockType(itemType entity.ItemType) entity.BlockType {
	// 由于Grass已合并到Dirt中，需要特殊处理
	if itemType == entity.Stone {
		return entity.StoneBlock
	} else if itemType == entity.Dirt || itemType == entity.Grass {
		return entity.DirtBlock
	} else if itemType == entity.Wood {
		return entity.WoodBlock
	} else if itemType == entity.Leaves {
		return entity.LeavesBlock
	}
	return entity.StoneBlock // 默认为石头
}

// placeBlock 放置方块
func (g *Game) placeBlock() {
	// 获取鼠标位置
	mx, my := ebiten.CursorPosition()
	
	// 转换为世界坐标
	worldX, worldY := g.camera.ScreenToWorld(float64(mx), float64(my))
	// 转换为网格坐标（使用math.Floor确保负数也能正确处理）
	gridX, gridY := int(math.Floor(worldX/32)), int(math.Floor(worldY/32))
	
	// 如果网格位置没有变化，则不重复操作
	if gridX == g.lastPlacePos[0] && gridY == g.lastPlacePos[1] {
		return
	}
	
	// 更新上次放置位置
	g.lastPlacePos[0], g.lastPlacePos[1] = gridX, gridY
	
	// 检查该位置是否已经有方块
	if g.world.IsBlockAt(gridX, gridY) {
		return // 如果已经有方块，则不放置也不消耗物品
	}
	
	// 获取当前选中的物品
	selectedItem := g.player.GetInventory().GetSelectedItem()
	if selectedItem.Type == entity.Air {
		return // 空气不能放置
	}
	
	// 添加对应类型的方块
	blockType := getItemToBlockType(selectedItem.Type)
	g.world.AddBlockWithType(gridX, gridY, blockType)
	
	// 消耗选中的物品
	g.player.GetInventory().ConsumeSelectedItem()
}

// destroyBlock 破坏方块
func (g *Game) destroyBlock() {
	// 获取鼠标位置
	mx, my := ebiten.CursorPosition()
	
	// 转换为世界坐标
	worldX, worldY := g.camera.ScreenToWorld(float64(mx), float64(my))
	// 转换为网格坐标（使用math.Floor确保负数也能正确处理）
	gridX, gridY := int(math.Floor(worldX/32)), int(math.Floor(worldY/32))
	
	// 如果网格位置没有变化，则不重复操作
	if gridX == g.lastDestroyPos[0] && gridY == g.lastDestroyPos[1] {
		return
	}
	
	// 更新上次破坏位置
	g.lastDestroyPos[0], g.lastDestroyPos[1] = gridX, gridY
	
	// 移除方块
	g.world.RemoveBlock(gridX, gridY)
}