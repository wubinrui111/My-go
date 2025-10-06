package game

import (
	"testing"
	"mygo/internal/pkg/entity"
	"mygo/internal/pkg/world"
)

// 创建一个简化版的Game用于测试，避免加载图像文件
func newTestGame() *Game {
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
		spriteSheet: nil, // 测试时不需要图像
	}
}

func TestNewGame(t *testing.T) {
	game := newTestGame()
	
	if game.player == nil {
		t.Error("Expected game to have a player")
	}
	
	if game.camera == nil {
		t.Error("Expected game to have a camera")
	}
	
	if game.world == nil {
		t.Error("Expected game to have a world")
	}
}

func TestGamePlaceAndDestroyBlock(t *testing.T) {
	game := newTestGame()
	
	// 检查初始状态下特定位置没有方块
	if game.world.IsBlockAt(0, 0) {
		t.Error("Expected no block at position (0, 0) initially")
	}
	
	// 模拟在(0, 0)位置放置方块
	game.world.AddBlock(0, 0)
	
	// 检查方块是否已放置
	if !game.world.IsBlockAt(0, 0) {
		t.Error("Expected block at position (0, 0) after placing")
	}
	
	// 模拟破坏(0, 0)位置的方块
	game.world.RemoveBlock(0, 0)
	
	// 检查方块是否已移除
	if game.world.IsBlockAt(0, 0) {
		t.Error("Expected no block at position (0, 0) after removal")
	}
}

func TestGamePlaceBlockWithType(t *testing.T) {
	game := newTestGame()
	
	// 检查初始状态下特定位置没有方块
	if game.world.IsBlockAt(1, 1) {
		t.Error("Expected no block at position (1, 1) initially")
	}
	
	// 模拟在(1, 1)位置放置特定类型的方块
	game.world.AddBlockWithType(1, 1, entity.DirtBlock)
	
	// 检查方块是否已放置
	block, exists := game.world.GetBlock(1, 1)
	if !exists {
		t.Error("Expected block at position (1, 1) after placing")
	}
	
	// 检查方块类型
	if block.GetType() != entity.DirtBlock {
		t.Errorf("Expected DirtBlock, got %v", block.GetType())
	}
}

func TestGameCoordinateTransformations(t *testing.T) {
	game := newTestGame()
	
	// 测试摄像机初始位置
	camX, camY := game.camera.GetPosition()
	if camX != 0 || camY != 0 {
		t.Errorf("Expected camera at (0, 0), got (%f, %f)", camX, camY)
	}
	
	// 测试屏幕到世界的坐标转换
	worldX, worldY := game.camera.ScreenToWorld(400, 300) // 屏幕中心
	if worldX != 0 || worldY != 0 {
		t.Errorf("Expected world center at (0, 0), got (%f, %f)", worldX, worldY)
	}
	
	// 测试世界到屏幕的坐标转换
	screenX, screenY := game.camera.WorldToScreen(0, 0) // 世界中心
	if screenX != 400 || screenY != 300 {
		t.Errorf("Expected screen center at (400, 300), got (%f, %f)", screenX, screenY)
	}
}

func TestGameNegativeCoordinateOperations(t *testing.T) {
	game := newTestGame()
	
	// 在负坐标位置放置方块
	game.world.AddBlock(-1, -1)
	
	// 检查方块是否存在
	if !game.world.IsBlockAt(-1, -1) {
		t.Error("Expected block at position (-1, -1)")
	}
	
	// 获取方块
	block, exists := game.world.GetBlock(-1, -1)
	if !exists {
		t.Error("Expected block to exist at position (-1, -1)")
	}
	
	if block == nil {
		t.Error("Expected block to be non-nil")
	}
	
	// 移除方块
	game.world.RemoveBlock(-1, -1)
	
	// 检查方块是否已移除
	if game.world.IsBlockAt(-1, -1) {
		t.Error("Expected no block at position (-1, -1) after removal")
	}
}

func TestWorldGeneration(t *testing.T) {
	// 创建游戏实例（简化版，避免加载图像）
	g := &Game{
		world: world.NewWorld(),
		player: entity.NewPlayer(0, 0),
	}
	
	// 生成世界地形
	g.GenerateWorldTerrain()
	
	// 检查是否生成了方块
	blocks := g.world.GetAllBlocks()
	
	if len(blocks) == 0 {
		t.Error("Expected world to have blocks generated, but found none")
	}
	
	// 检查是否生成了不同类型的方块
	blockTypes := make(map[entity.BlockType]int)
	for _, block := range blocks {
		blockTypes[block.GetType()]++
	}
	
	if len(blockTypes) < 3 {
		t.Errorf("Expected at least 3 different block types, got %d", len(blockTypes))
	}
	
	t.Logf("Generated %d blocks with %d different types", len(blocks), len(blockTypes))
}

func TestPerlinNoise(t *testing.T) {
	// 测试噪声生成器
	noise := NewPerlinNoise(12345)
	
	// 测试噪声值范围
	value := noise.Noise(0.5, 0.5)
	if value < -1.0 || value > 1.0 {
		t.Errorf("Noise value out of range: %f", value)
	}
	
	// 测试FBM值范围
	fbmValue := noise.FBM(0.5, 0.5, 1.0, 1.0, 4)
	if fbmValue < -1.0 || fbmValue > 1.0 {
		t.Errorf("FBM value out of range: %f", fbmValue)
	}
	
	t.Logf("Noise value: %f, FBM value: %f", value, fbmValue)
}