package entity

import (
	"testing"
)

func TestNewInventory(t *testing.T) {
	inventory := NewInventory()
	
	if len(inventory.Slots) != TotalSlotCount {
		t.Errorf("Expected %d slots, got %d", TotalSlotCount, len(inventory.Slots))
	}
	
	if inventory.SelectedSlot != 0 {
		t.Errorf("Expected selected slot to be 0, got %d", inventory.SelectedSlot)
	}
	
	if inventory.Open {
		t.Error("Expected inventory to be closed initially")
	}
}

func TestGetHotbarSlot(t *testing.T) {
	inventory := NewInventory()
	
	// 测试有效槽位
	item := inventory.GetHotbarSlot(0)
	if item.Type != Stone {
		t.Errorf("Expected Stone in slot 0, got %v", item.Type)
	}
	
	if item.Count != 64 {
		t.Errorf("Expected count 64 in slot 0, got %d", item.Count)
	}
	
	// 测试无效槽位
	item = inventory.GetHotbarSlot(-1)
	if item.Type != Air {
		t.Errorf("Expected Air for invalid slot, got %v", item.Type)
	}
	
	item = inventory.GetHotbarSlot(HotbarSlotCount)
	if item.Type != Air {
		t.Errorf("Expected Air for invalid slot, got %v", item.Type)
	}
}

func TestGetSlot(t *testing.T) {
	inventory := NewInventory()
	
	// 测试有效槽位
	item := inventory.GetSlot(0)
	if item.Type != Stone {
		t.Errorf("Expected Stone in slot 0, got %v", item.Type)
	}
	
	// 测试无效槽位
	item = inventory.GetSlot(-1)
	if item.Type != Air {
		t.Errorf("Expected Air for invalid slot, got %v", item.Type)
	}
	
	item = inventory.GetSlot(TotalSlotCount)
	if item.Type != Air {
		t.Errorf("Expected Air for invalid slot, got %v", item.Type)
	}
}

func TestSetSlot(t *testing.T) {
	inventory := NewInventory()
	
	// 测试设置有效槽位
	newItem := ItemStack{Type: Dirt, Count: 32}
	inventory.SetSlot(5, newItem)
	
	item := inventory.GetSlot(5)
	if item.Type != Dirt {
		t.Errorf("Expected Dirt in slot 5, got %v", item.Type)
	}
	
	if item.Count != 32 {
		t.Errorf("Expected count 32 in slot 5, got %d", item.Count)
	}
	
	// 测试设置无效槽位（不应崩溃）
	inventory.SetSlot(-1, newItem)
	inventory.SetSlot(TotalSlotCount, newItem)
}

func TestGetSelectedItem(t *testing.T) {
	inventory := NewInventory()
	
	// 测试默认选中项
	item := inventory.GetSelectedItem()
	if item.Type != Stone {
		t.Errorf("Expected Stone as selected item, got %v", item.Type)
	}
	
	// 测试切换选中项后
	inventory.SetSelectedSlot(1)
	item = inventory.GetSelectedItem()
	if item.Type != Dirt {
		t.Errorf("Expected Dirt as selected item, got %v", item.Type)
	}
}

func TestInventoryToggle(t *testing.T) {
	inventory := NewInventory()
	
	// 初始状态应为关闭
	if inventory.IsOpen() {
		t.Error("Expected inventory to be closed initially")
	}
	
	// 切换状态
	inventory.ToggleOpen()
	if !inventory.IsOpen() {
		t.Error("Expected inventory to be open after toggle")
	}
	
	// 再次切换状态
	inventory.ToggleOpen()
	if inventory.IsOpen() {
		t.Error("Expected inventory to be closed after second toggle")
	}
}

func TestSetSelectedSlot(t *testing.T) {
	inventory := NewInventory()
	
	// 测试设置有效槽位
	inventory.SetSelectedSlot(5)
	if inventory.SelectedSlot != 5 {
		t.Errorf("Expected selected slot to be 5, got %d", inventory.SelectedSlot)
	}
	
	// 测试设置无效槽位（应保持原值）
	inventory.SetSelectedSlot(-1)
	if inventory.SelectedSlot != 5 {
		t.Errorf("Expected selected slot to remain 5, got %d", inventory.SelectedSlot)
	}
	
	inventory.SetSelectedSlot(HotbarSlotCount)
	if inventory.SelectedSlot != 5 {
		t.Errorf("Expected selected slot to remain 5, got %d", inventory.SelectedSlot)
	}
}

func TestGetHotbar(t *testing.T) {
	inventory := NewInventory()
	hotbar := inventory.GetHotbar()
	
	if len(hotbar) != HotbarSlotCount {
		t.Errorf("Expected hotbar length %d, got %d", HotbarSlotCount, len(hotbar))
	}
	
	// 检查前几个槽位
	if hotbar[0].Type != Stone {
		t.Errorf("Expected Stone in hotbar slot 0, got %v", hotbar[0].Type)
	}
	
	if hotbar[1].Type != Dirt {
		t.Errorf("Expected Dirt in hotbar slot 1, got %v", hotbar[1].Type)
	}
}

func TestGetInventory(t *testing.T) {
	inventory := NewInventory()
	inventorySlots := inventory.GetInventory()
	
	if len(inventorySlots) != InventorySlotCount {
		t.Errorf("Expected inventory length %d, got %d", InventorySlotCount, len(inventorySlots))
	}
	
	// 检查所有槽位初始应为空气
	for i, item := range inventorySlots {
		if item.Type != Air {
			t.Errorf("Expected Air in inventory slot %d, got %v", i, item.Type)
		}
	}
}

func TestConsumeSelectedItem(t *testing.T) {
	inventory := NewInventory()
	
	// 初始状态下选中物品应该是Stone，数量为64
	item := inventory.GetSelectedItem()
	if item.Type != Stone || item.Count != 64 {
		t.Errorf("Expected Stone with count 64, got %v with count %d", item.Type, item.Count)
	}
	
	// 消耗一个物品
	inventory.ConsumeSelectedItem()
	
	// 检查数量是否减少
	item = inventory.GetSelectedItem()
	if item.Count != 63 {
		t.Errorf("Expected count 63 after consuming, got %d", item.Count)
	}
	
	// 设置选中槽位数量为1并消耗
	inventory.SetSlot(inventory.GetSelectedSlot(), ItemStack{Type: Stone, Count: 1})
	inventory.ConsumeSelectedItem()
	
	// 检查物品是否被清空
	item = inventory.GetSelectedItem()
	if item.Type != Air || item.Count != 0 {
		t.Errorf("Expected Air with count 0 after consuming last item, got %v with count %d", item.Type, item.Count)
	}
}