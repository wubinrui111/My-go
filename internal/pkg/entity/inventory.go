package entity

const (
	// HotbarSlotCount 底部快捷栏槽数量
	HotbarSlotCount = 9
	// InventorySlotCount 完整物品栏槽数量
	InventorySlotCount = 27
	// TotalSlotCount 总槽数量
	TotalSlotCount = HotbarSlotCount + InventorySlotCount
)

// ItemType 物品类型
type ItemType int

const (
	Air ItemType = iota
	Stone
	Dirt
	Grass
	Wood
	Leaves
)

// ItemStack 物品堆叠
type ItemStack struct {
	Type  ItemType
	Count int
}

// Inventory 物品栏
type Inventory struct {
	Slots       []ItemStack
	SelectedSlot int // 当前选中的快捷栏槽位 (0-8)
	Open        bool // 物品栏是否展开
}

// NewInventory 创建新的物品栏
func NewInventory() *Inventory {
	slots := make([]ItemStack, TotalSlotCount)
	for i := range slots {
		slots[i] = ItemStack{Type: Air, Count: 0}
	}
	
	// 初始化一些测试物品
	slots[0] = ItemStack{Type: Stone, Count: 64}
	slots[1] = ItemStack{Type: Dirt, Count: 64}
	slots[2] = ItemStack{Type: Grass, Count: 64}
	slots[3] = ItemStack{Type: Wood, Count: 64}
	slots[4] = ItemStack{Type: Leaves, Count: 64}
	
	return &Inventory{
		Slots:       slots,
		SelectedSlot: 0,
		Open:        false,
	}
}

// GetHotbarSlot 获取快捷栏指定槽位的物品
func (inv *Inventory) GetHotbarSlot(slot int) ItemStack {
	if slot < 0 || slot >= HotbarSlotCount {
		return ItemStack{Type: Air, Count: 0}
	}
	return inv.Slots[slot]
}

// GetSlot 获取指定槽位的物品
func (inv *Inventory) GetSlot(slot int) ItemStack {
	if slot < 0 || slot >= TotalSlotCount {
		return ItemStack{Type: Air, Count: 0}
	}
	return inv.Slots[slot]
}

// SetSlot 设置指定槽位的物品
func (inv *Inventory) SetSlot(slot int, item ItemStack) {
	if slot < 0 || slot >= TotalSlotCount {
		return
	}
	inv.Slots[slot] = item
}

// GetSelectedSlot 获取当前选中的槽位索引
func (inv *Inventory) GetSelectedSlot() int {
	return inv.SelectedSlot
}

// SetSelectedSlot 设置当前选中的槽位索引
func (inv *Inventory) SetSelectedSlot(slot int) {
	if slot >= 0 && slot < HotbarSlotCount {
		inv.SelectedSlot = slot
	}
}

// GetSelectedItem 获取当前选中的物品
func (inv *Inventory) GetSelectedItem() ItemStack {
	return inv.GetHotbarSlot(inv.SelectedSlot)
}

// ConsumeSelectedItem 消耗选中的物品（放置方块时调用）
func (inv *Inventory) ConsumeSelectedItem() {
	slot := inv.SelectedSlot
	if inv.Slots[slot].Count > 0 {
		inv.Slots[slot].Count--
		// 如果数量为0，则清空该槽位
		if inv.Slots[slot].Count == 0 {
			inv.Slots[slot] = ItemStack{Type: Air, Count: 0}
		}
	}
}

// IsOpen 返回物品栏是否展开
func (inv *Inventory) IsOpen() bool {
	return inv.Open
}

// ToggleOpen 切换物品栏展开状态
func (inv *Inventory) ToggleOpen() {
	inv.Open = !inv.Open
}

// OpenInventory 展开物品栏
func (inv *Inventory) OpenInventory() {
	inv.Open = true
}

// CloseInventory 关闭物品栏
func (inv *Inventory) CloseInventory() {
	inv.Open = false
}

// GetHotbar(获取底部快捷栏物品
func (inv *Inventory) GetHotbar() []ItemStack {
	hotbar := make([]ItemStack, HotbarSlotCount)
	for i := 0; i < HotbarSlotCount; i++ {
		hotbar[i] = inv.Slots[i]
	}
	return hotbar
}

// GetInventory(获取完整物品栏物品（不包括快捷栏）
func (inv *Inventory) GetInventory() []ItemStack {
	inventory := make([]ItemStack, InventorySlotCount)
	for i := 0; i < InventorySlotCount; i++ {
		inventory[i] = inv.Slots[HotbarSlotCount+i]
	}
	return inventory
}