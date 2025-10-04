package entity

const (
	BlockSize = 32
)

// BlockType 方块类型
type BlockType int

const (
	StoneBlock BlockType = iota
	DirtBlock
	GrassBlock
	WoodBlock
	LeavesBlock
)

// Block represents a block in the game world
type Block struct {
	X, Y   float64
	Type   BlockType  // 方块类型
}

// NewBlock creates a new block at the specified grid position
func NewBlock(x, y int) *Block {
	return &Block{
		X: float64(x * BlockSize),
		Y: float64(y * BlockSize),
		Type: StoneBlock, // 默认类型为石头
	}
}

// NewBlockWithType 创建指定类型的方块
func NewBlockWithType(x, y int, blockType BlockType) *Block {
	return &Block{
		X: float64(x * BlockSize),
		Y: float64(y * BlockSize),
		Type: blockType,
	}
}

// GetPosition returns the world position of the block
func (b *Block) GetPosition() (float64, float64) {
	return b.X, b.Y
}

// GetGridPosition returns the grid position of the block
func (b *Block) GetGridPosition() (int, int) {
	return int(b.X / BlockSize), int(b.Y / BlockSize)
}

// GetType 获取方块类型
func (b *Block) GetType() BlockType {
	return b.Type
}