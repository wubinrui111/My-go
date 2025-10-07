package entity

import (
	"testing"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestItemEntityDraw(t *testing.T) {
	// 创建测试用的精灵表（简单颜色块）
	spriteSheet := ebiten.NewImage(640, 640)
	spriteSheet.Fill(color.RGBA{255, 255, 255, 255})
	
	// 创建一个掉落物
	item := NewItemEntityFromBlock(100, 50, StoneBlock, 1)
	
	// 创建测试用的屏幕
	screen := ebiten.NewImage(800, 600)
	
	// 测试绘制方法不会崩溃
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw方法发生崩溃: %v", r)
		}
	}()
	
	item.Draw(screen, spriteSheet)
}

func TestItemEntityDrawWithCamera(t *testing.T) {
	// 创建测试用的精灵表（简单颜色块）
	spriteSheet := ebiten.NewImage(640, 640)
	spriteSheet.Fill(color.RGBA{255, 255, 255, 255})
	
	// 创建相机
	camera := NewCamera(0, 0)
	camera.SetScreenSize(800, 600)
	
	// 创建一个掉落物
	item := NewItemEntityFromBlock(100, 50, StoneBlock, 1)
	
	// 创建测试用的屏幕
	screen := ebiten.NewImage(800, 600)
	
	// 测试带相机的绘制方法不会崩溃
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DrawWithCamera方法发生崩溃: %v", r)
		}
	}()
	
	item.DrawWithCamera(screen, spriteSheet, camera)
}

func TestItemEntityDrawWithCount(t *testing.T) {
	// 创建测试用的精灵表（简单颜色块）
	spriteSheet := ebiten.NewImage(640, 640)
	spriteSheet.Fill(color.RGBA{255, 255, 255, 255})
	
	// 创建一个数量大于1的掉落物
	item := NewItemEntityFromBlock(100, 50, StoneBlock, 5)
	
	// 创建测试用的屏幕
	screen := ebiten.NewImage(800, 600)
	
	// 测试绘制方法不会崩溃
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw方法发生崩溃: %v", r)
		}
	}()
	
	item.Draw(screen, spriteSheet)
}

func TestItemEntityDrawWithCameraAndCount(t *testing.T) {
	// 创建测试用的精灵表（简单颜色块）
	spriteSheet := ebiten.NewImage(640, 640)
	spriteSheet.Fill(color.RGBA{255, 255, 255, 255})
	
	// 创建相机
	camera := NewCamera(0, 0)
	camera.SetScreenSize(800, 600)
	
	// 创建一个数量大于1的掉落物
	item := NewItemEntityFromBlock(100, 50, StoneBlock, 5)
	
	// 创建测试用的屏幕
	screen := ebiten.NewImage(800, 600)
	
	// 测试带相机的绘制方法不会崩溃
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DrawWithCamera方法发生崩溃: %v", r)
		}
	}()
	
	item.DrawWithCamera(screen, spriteSheet, camera)
}