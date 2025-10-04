package entity

import (
	"testing"
)

func TestNewItemEntity(t *testing.T) {
	item := NewItemEntity(100, 50, Stone, 5)
	
	if item.X != 100 {
		t.Errorf("Expected X=100, got %f", item.X)
	}
	
	if item.Y != 50 {
		t.Errorf("Expected Y=50, got %f", item.Y)
	}
	
	if item.ItemType != Stone {
		t.Errorf("Expected ItemType=Stone, got %v", item.ItemType)
	}
	
	if item.Count != 5 {
		t.Errorf("Expected Count=5, got %d", item.Count)
	}
	
	if item.Lifetime != 600 {
		t.Errorf("Expected Lifetime=600, got %d", item.Lifetime)
	}
}

func TestItemEntityGetPosition(t *testing.T) {
	item := NewItemEntity(100, 50, Stone, 5)
	
	x, y := item.GetPosition()
	
	if x != 100 {
		t.Errorf("Expected x=100, got %f", x)
	}
	
	if y != 50 {
		t.Errorf("Expected y=50, got %f", y)
	}
}

func TestItemEntityGetItemType(t *testing.T) {
	item := NewItemEntity(100, 50, Dirt, 3)
	
	itemType := item.GetItemType()
	
	if itemType != Dirt {
		t.Errorf("Expected ItemType=Dirt, got %v", itemType)
	}
}

func TestItemEntityGetCount(t *testing.T) {
	item := NewItemEntity(100, 50, Grass, 7)
	
	count := item.GetCount()
	
	if count != 7 {
		t.Errorf("Expected Count=7, got %d", count)
	}
}

func TestItemEntityIsExpired(t *testing.T) {
	item := NewItemEntity(100, 50, Wood, 1)
	
	// 初始状态下不应该过期
	if item.IsExpired() {
		t.Error("Expected item not to be expired initially")
	}
	
	// 将生命周期减少到0
	item.Lifetime = 0
	
	// 现在应该过期
	if !item.IsExpired() {
		t.Error("Expected item to be expired after lifetime reaches 0")
	}
}

func TestItemEntityTryPickup(t *testing.T) {
	item := NewItemEntity(100, 50, Leaves, 2)
	
	// 测试在拾取范围内的位置
	if !item.TryPickup(100, 50) { // 距离0
		t.Error("Expected pickup to succeed when player is at the same position")
	}
	
	if !item.TryPickup(90, 50) { // 距离10
		t.Error("Expected pickup to succeed when player is within range")
	}
	
	if !item.TryPickup(100, 40) { // 距离10
		t.Error("Expected pickup to succeed when player is within range")
	}
	
	// 测试在拾取范围外的位置
	if item.TryPickup(100, 90) { // 距离40
		t.Error("Expected pickup to fail when player is outside range")
	}
}