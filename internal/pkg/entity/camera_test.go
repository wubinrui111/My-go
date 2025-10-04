package entity

import (
	"testing"
)

func TestNewCamera(t *testing.T) {
	camera := NewCamera(10, 20)
	
	if camera.X != 10 {
		t.Errorf("Expected X=10, got %f", camera.X)
	}
	
	if camera.Y != 20 {
		t.Errorf("Expected Y=20, got %f", camera.Y)
	}
	
	if camera.FollowSpeed != 0.1 {
		t.Errorf("Expected FollowSpeed=0.1, got %f", camera.FollowSpeed)
	}
}

func TestCameraSetTarget(t *testing.T) {
	camera := NewCamera(0, 0)
	
	camera.SetTarget(100, 200)
	
	if camera.TargetX != 100 {
		t.Errorf("Expected TargetX=100, got %f", camera.TargetX)
	}
	
	if camera.TargetY != 200 {
		t.Errorf("Expected TargetY=200, got %f", camera.TargetY)
	}
}

func TestCameraSetScreenSize(t *testing.T) {
	camera := NewCamera(0, 0)
	
	camera.SetScreenSize(800, 600)
	
	if camera.OffsetX != 400 {
		t.Errorf("Expected OffsetX=400, got %f", camera.OffsetX)
	}
	
	if camera.OffsetY != 300 {
		t.Errorf("Expected OffsetY=300, got %f", camera.OffsetY)
	}
}

func TestCameraWorldToScreen(t *testing.T) {
	camera := NewCamera(0, 0)
	camera.SetScreenSize(800, 600)
	
	// 测试原点
	screenX, screenY := camera.WorldToScreen(0, 0)
	if screenX != 400 {
		t.Errorf("Expected screenX=400 for worldX=0, got %f", screenX)
	}
	
	if screenY != 300 {
		t.Errorf("Expected screenY=300 for worldY=0, got %f", screenY)
	}
	
	// 测试偏移位置
	screenX, screenY = camera.WorldToScreen(32, 32)
	if screenX != 432 {
		t.Errorf("Expected screenX=432 for worldX=32, got %f", screenX)
	}
	
	if screenY != 332 {
		t.Errorf("Expected screenY=332 for worldY=32, got %f", screenY)
	}
	
	// 测试负坐标
	screenX, screenY = camera.WorldToScreen(-32, -32)
	if screenX != 368 {
		t.Errorf("Expected screenX=368 for worldX=-32, got %f", screenX)
	}
	
	if screenY != 268 {
		t.Errorf("Expected screenY=268 for worldY=-32, got %f", screenY)
	}
	
	// 测试摄像机偏移
	camera.X = 100
	camera.Y = 50
	screenX, screenY = camera.WorldToScreen(0, 0)
	if screenX != 300 {
		t.Errorf("Expected screenX=300 for worldX=0 with cameraX=100, got %f", screenX)
	}
	
	if screenY != 250 {
		t.Errorf("Expected screenY=250 for worldY=0 with cameraY=50, got %f", screenY)
	}
}

func TestCameraScreenToWorld(t *testing.T) {
	camera := NewCamera(0, 0)
	camera.SetScreenSize(800, 600)
	
	// 测试屏幕中心
	worldX, worldY := camera.ScreenToWorld(400, 300)
	if worldX != 0 {
		t.Errorf("Expected worldX=0 for screenX=400, got %f", worldX)
	}
	
	if worldY != 0 {
		t.Errorf("Expected worldY=0 for screenY=300, got %f", worldY)
	}
	
	// 测试偏移位置
	worldX, worldY = camera.ScreenToWorld(432, 332)
	if worldX != 32 {
		t.Errorf("Expected worldX=32 for screenX=432, got %f", worldX)
	}
	
	if worldY != 32 {
		t.Errorf("Expected worldY=32 for screenY=332, got %f", worldY)
	}
	
	// 测试负坐标映射
	worldX, worldY = camera.ScreenToWorld(368, 268)
	if worldX != -32 {
		t.Errorf("Expected worldX=-32 for screenX=368, got %f", worldX)
	}
	
	if worldY != -32 {
		t.Errorf("Expected worldY=-32 for screenY=268, got %f", worldY)
	}
	
	// 测试摄像机偏移
	camera.X = 100
	camera.Y = 50
	worldX, worldY = camera.ScreenToWorld(400, 300)
	if worldX != 100 {
		t.Errorf("Expected worldX=100 for screenX=400 with cameraX=100, got %f", worldX)
	}
	
	if worldY != 50 {
		t.Errorf("Expected worldY=50 for screenY=300 with cameraY=50, got %f", worldY)
	}
}

func TestCameraUpdate(t *testing.T) {
	camera := NewCamera(0, 0)
	camera.SetTarget(100, 50)
	
	// 初始位置
	if camera.X != 0 {
		t.Errorf("Expected initial X=0, got %f", camera.X)
	}
	
	if camera.Y != 0 {
		t.Errorf("Expected initial Y=0, got %f", camera.Y)
	}
	
	// 更新一次
	camera.Update()
	
	// 检查是否向目标移动
	if camera.X <= 0 {
		t.Errorf("Expected X to increase, got %f", camera.X)
	}
	
	if camera.Y <= 0 {
		t.Errorf("Expected Y to increase, got %f", camera.Y)
	}
	
	if camera.X >= 100 {
		t.Errorf("Expected X to be less than target, got %f", camera.X)
	}
	
	if camera.Y >= 50 {
		t.Errorf("Expected Y to be less than target, got %f", camera.Y)
	}
}

func TestCameraGetPosition(t *testing.T) {
	camera := NewCamera(123.45, 678.90)
	
	x, y := camera.GetPosition()
	
	if x != 123.45 {
		t.Errorf("Expected X=123.45, got %f", x)
	}
	
	if y != 678.90 {
		t.Errorf("Expected Y=678.90, got %f", y)
	}
}