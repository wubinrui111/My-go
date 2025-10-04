package entity

// Camera represents a camera that follows the player
type Camera struct {
	X, Y   float64 // Camera position
	OffsetX, OffsetY float64 // Offset to center the player on screen
	TargetX, TargetY float64 // Target position (player position)
	FollowSpeed float64 // Speed at which camera follows the player
}

// NewCamera creates a new camera
func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
		FollowSpeed: 0.1,
	}
}

// Update updates the camera position to follow the target smoothly
func (c *Camera) Update() {
	// Smoothly follow the target
	c.X += (c.TargetX - c.X) * c.FollowSpeed
	c.Y += (c.TargetY - c.Y) * c.FollowSpeed
}

// SetTarget sets the target position for the camera to follow
func (c *Camera) SetTarget(x, y float64) {
	c.TargetX = x
	c.TargetY = y
}

// SetScreenSize sets the screen size to calculate proper offsets
func (c *Camera) SetScreenSize(width, height int) {
	c.OffsetX = float64(width) / 2
	c.OffsetY = float64(height) / 2
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldX, worldY float64) (float64, float64) {
	screenX := worldX - c.X + c.OffsetX
	screenY := worldY - c.Y + c.OffsetY
	return screenX, screenY
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenX, screenY float64) (float64, float64) {
	worldX := screenX + c.X - c.OffsetX
	worldY := screenY + c.Y - c.OffsetY
	return worldX, worldY
}

// GetPosition returns the camera position
func (c *Camera) GetPosition() (float64, float64) {
	return c.X, c.Y
}