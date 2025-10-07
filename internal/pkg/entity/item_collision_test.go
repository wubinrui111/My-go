package entity

import (
	"testing"
	"math"
)

func TestItemFriction(t *testing.T) {
	// 创建一个模拟世界
	mockWorld := NewMockWorld()
	
	// 创建一个掉落物
	item := NewItemEntityFromBlock(0, 0, StoneBlock, 1)
	item.SetWorld(mockWorld)
	
	// 设置一个初始速度
	initialVX := 5.0
	item.VX = initialVX
	
	// 更新几次观察摩擦力效果
	for i := 0; i < 10; i++ {
		item.Update()
	}
	
	// 检查速度是否因摩擦力而减小
	if math.Abs(item.VX) >= initialVX {
		t.Error("摩擦力应该减小水平速度")
	}
	
	// 检查物品是否最终停止
	if math.Abs(item.VX) > 0.1 {
		// 给予一些容差，继续更新直到几乎停止
		for i := 0; i < 100; i++ {
			item.Update()
		}
		if math.Abs(item.VX) > 0.1 {
			t.Error("物品应该因摩擦力而最终停止")
		}
	}
}

func TestItemHorizontalCollision(t *testing.T) {
	// 创建一个模拟世界
	mockWorld := NewMockWorld()
	
	// 在物品右侧添加一个墙壁
	mockWorld.AddBlock(2, 1) // 墙壁在网格(2,1)
	
	// 创建一个掉落物
	item := NewItemEntityFromBlock(0, 0, StoneBlock, 1)
	item.SetWorld(mockWorld)
	
	// 将物品放置在网格(1,1)内，向右移动
	item.X = float64(1*32 + 16) // 网格1的中心
	item.Y = float64(1*32 + 16) // 网格1的中心
	item.VX = 5.0 // 向右移动
	
	// 记录初始速度
	initialVX := item.VX
	
	// 更新直到发生碰撞
	for i := 0; i < 50; i++ {
		item.Update()
		// 检查是否发生碰撞（速度方向改变或减小）
		if item.VX <= 0 || math.Abs(item.VX) < initialVX {
			break
		}
	}
	
	// 检查速度是否因碰撞而改变
	if item.VX >= initialVX {
		t.Error("物品碰到墙壁后应该改变水平速度")
	}
}

func TestItemVerticalCollision(t *testing.T) {
	// 创建一个模拟世界
	mockWorld := NewMockWorld()
	
	// 在物品下方添加一个地面
	mockWorld.AddBlock(1, 2) // 地面在网格(1,2)
	
	// 创建一个掉落物
	item := NewItemEntityFromBlock(0, 0, StoneBlock, 1)
	item.SetWorld(mockWorld)
	
	// 将物品放置在网格(1,1)内，向下移动
	item.X = float64(1*32 + 16) // 网格1的中心
	item.Y = float64(1*32 + 16) // 网格1的中心
	item.VY = 5.0 // 向下移动
	
	// 记录初始速度
	initialVY := item.VY
	
	// 更新直到发生碰撞
	for i := 0; i < 50; i++ {
		item.Update()
		// 检查是否发生碰撞（速度减小或变为0）
		if item.VY <= 0 {
			break
		}
	}
	
	// 检查速度是否因碰撞而改变
	if item.VY >= initialVY {
		t.Error("物品碰到地面后应该改变垂直速度")
	}
}

func TestItemCollisionBox(t *testing.T) {
	// 创建一个模拟世界
	mockWorld := NewMockWorld()
	
	// 在物品左侧添加一个墙壁
	mockWorld.AddBlock(0, 1) // 墙壁在网格(0,1)
	
	// 创建一个掉落物
	item := NewItemEntityFromBlock(0, 0, StoneBlock, 1)
	item.SetWorld(mockWorld)
	
	// 将物品放置在网格(1,1)内
	item.X = float64(1*32 + 16) // 网格1的中心
	item.Y = float64(1*32 + 16) // 网格1的中心
	
	// 检查碰撞箱边界计算是否正确
	left := int(math.Floor((item.X - ItemSize/2) / 32))
	right := int(math.Floor((item.X + ItemSize/2 - 1) / 32))
	top := int(math.Floor((item.Y - ItemSize/2) / 32))
	bottom := int(math.Floor((item.Y + ItemSize/2 - 1) / 32))
	
	// 对于16x16的物品在(48,48)位置，应该占据网格(1,1)
	expectedLeft := 1
	expectedRight := 1
	expectedTop := 1
	expectedBottom := 1
	
	if left != expectedLeft || right != expectedRight || top != expectedTop || bottom != expectedBottom {
		t.Errorf("碰撞箱计算错误，期望: [%d,%d,%d,%d]，实际: [%d,%d,%d,%d]",
			expectedLeft, expectedRight, expectedTop, expectedBottom,
			left, right, top, bottom)
	}
}

func TestItemStopOnGround(t *testing.T) {
	// 创建一个模拟世界
	mockWorld := NewMockWorld()
	
	// 在物品下方添加一个地面
	mockWorld.AddBlock(1, 2) // 地面在网格(1,2)
	
	// 创建一个掉落物
	item := NewItemEntityFromBlock(0, 0, StoneBlock, 1)
	item.SetWorld(mockWorld)
	
	// 将物品放置在网格(1,1)内，向下移动
	item.X = float64(1*32 + 16) // 网格1的中心
	item.Y = float64(1*32 + 16) // 网格1的中心
	item.VY = 1.0 // 向下移动
	
	// 更新多次直到物品应该停止
	for i := 0; i < 100; i++ {
		item.Update()
	}
	
	// 检查物品是否停止（速度接近0）
	if math.Abs(item.VY) > 0.1 {
		t.Error("物品在地面上应该停止移动")
	}
	
	// 检查物品是否在地面上方正确位置
	expectedY := float64(2*32 - ItemSize/2)
	if math.Abs(item.Y-expectedY) > 1 {
		t.Errorf("物品应在地面上方停止，期望Y坐标: %f, 实际: %f", expectedY, item.Y)
	}
}