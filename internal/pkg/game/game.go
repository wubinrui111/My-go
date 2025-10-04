package game

import (
	"image/color"
	"fmt"
	"image"
	_ "image/png"

	"mygo/internal/pkg/entity"
	"mygo/internal/pkg/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	player *entity.Player
	camera *entity.Camera
	world  *world.World
	// 连续放置/破坏方块相关变量
	lastPlacePos    [2]int // 记录上次放置方块的网格位置
	lastDestroyPos  [2]int // 记录上次破坏方块的网格位置
	spriteSheet     *ebiten.Image // 精灵表
}

func NewGame() *Game {
	// 加载精灵表
	spriteSheet, _, err := ebitenutil.NewImageFromFile("image/test.png")
	if err != nil {
		panic(err)
	}
	
	// 创建世界
	w := world.NewWorld()
	
	// 添加一些测试方块
	for x := -10; x < 10; x++ {
		w.AddBlock(x, 5) // 在y=5处添加一行方块作为地面
	}
	
	// 添加一些垂直方块用于测试碰撞
	for y := 0; y < 5; y++ {
		w.AddBlock(-5, y)  // 左侧墙壁
		w.AddBlock(5, y)   // 右侧墙壁
	}
	
	// 将玩家放置在地面上方
	w.Player.SetPosition(0, -64)
	
	// 创建摄像机
	camera := entity.NewCamera(0, 0)
	camera.SetScreenSize(800, 600)
	
	return &Game{
		player: w.Player,
		camera: camera,
		world:  w,
		lastPlacePos:    [2]int{-1, -1}, // 初始化为无效位置
		lastDestroyPos:  [2]int{-1, -1}, // 初始化为无效位置
		spriteSheet: spriteSheet,
	}
}

func (g *Game) Update() error {
	// 更新玩家输入
	g.handleInput()
	
	// 更新玩家状态
	g.player.Update()
	
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
		g.drawSprite(screen, screenX, screenY, int(block.GetType())+1) // +1因为0是玩家
	}
	
	// 绘制冲刺残影
	dashTrails := g.player.GetDashTrails()
	for _, trail := range dashTrails {
		screenX, screenY := g.camera.WorldToScreen(trail.X, trail.Y)
		// 使用半透明红色方块表示残影
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenX-16, screenY-16)
		op.ColorM.Scale(1, 1, 1, trail.Alpha*0.5) // 设置透明度
		// 绘制玩家精灵（下标0）作为残影
		g.drawSpriteWithOp(screen, op, 0)
	}
	
	// 绘制玩家
	playerX, playerY := g.player.GetPosition()
	screenX, screenY := g.camera.WorldToScreen(playerX, playerY)
	// 绘制玩家精灵（下标0）
	g.drawSprite(screen, screenX-16, screenY-16, 0)
	
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
	switch itemType {
	case entity.Stone:
		return color.RGBA{128, 128, 128, 255} // 灰色
	case entity.Dirt:
		return color.RGBA{150, 100, 50, 255}  // 棕色
	case entity.Grass:
		return color.RGBA{50, 180, 50, 255}   // 绿色
	case entity.Wood:
		return color.RGBA{150, 100, 50, 255}  // 棕色
	case entity.Leaves:
		return color.RGBA{30, 120, 30, 255}   // 深绿色
	default:
		return color.RGBA{255, 0, 255, 255}   // 品红色（默认）
	}
}

// getItemToSpriteIndex 将物品类型转换为精灵表索引
func getItemToSpriteIndex(itemType entity.ItemType) int {
	switch itemType {
	case entity.Stone:
		return 1 // 石头精灵在索引1
	case entity.Dirt:
		return 2 // 泥土精灵在索引2
	case entity.Grass:
		return 3 // 草精灵在索引3
	case entity.Wood:
		return 4 // 木头精灵在索引4
	case entity.Leaves:
		return 5 // 树叶精灵在索引5
	default:
		return 1 // 默认为石头
	}
}

// getBlockColor 根据方块类型获取颜色（备用方案）
func getBlockColor(blockType entity.BlockType) color.RGBA {
	switch blockType {
	case entity.StoneBlock:
		return color.RGBA{128, 128, 128, 255} // 灰色
	case entity.DirtBlock:
		return color.RGBA{150, 100, 50, 255}  // 棕色
	case entity.GrassBlock:
		return color.RGBA{50, 180, 50, 255}   // 绿色
	case entity.WoodBlock:
		return color.RGBA{150, 100, 50, 255}  // 棕色
	case entity.LeavesBlock:
		return color.RGBA{30, 120, 30, 255}   // 深绿色
	default:
		return color.RGBA{139, 69, 19, 255}   // 棕色（默认）
	}
}

// getItemToBlockType 将物品类型转换为方块类型
func getItemToBlockType(itemType entity.ItemType) entity.BlockType {
	switch itemType {
	case entity.Stone:
		return entity.StoneBlock
	case entity.Dirt:
		return entity.DirtBlock
	case entity.Grass:
		return entity.GrassBlock
	case entity.Wood:
		return entity.WoodBlock
	case entity.Leaves:
		return entity.LeavesBlock
	default:
		return entity.StoneBlock // 默认为石头
	}
}

// placeBlock 放置方块
func (g *Game) placeBlock() {
	// 获取鼠标位置
	mx, my := ebiten.CursorPosition()
	
	// 转换为世界坐标
	worldX, worldY := g.camera.ScreenToWorld(float64(mx), float64(my))
	// 转换为网格坐标
	gridX, gridY := int(worldX/32), int(worldY/32)
	
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
	// 转换为网格坐标
	gridX, gridY := int(worldX/32), int(worldY/32)
	
	// 如果网格位置没有变化，则不重复操作
	if gridX == g.lastDestroyPos[0] && gridY == g.lastDestroyPos[1] {
		return
	}
	
	// 更新上次破坏位置
	g.lastDestroyPos[0], g.lastDestroyPos[1] = gridX, gridY
	
	// 移除方块
	g.world.RemoveBlock(gridX, gridY)
}