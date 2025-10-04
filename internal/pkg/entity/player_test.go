package entity

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	player := NewPlayer(10, 20)
	
	if player.X != 10 {
		t.Errorf("Expected X=10, got %f", player.X)
	}
	
	if player.Y != 20 {
		t.Errorf("Expected Y=20, got %f", player.Y)
	}
	
	if !player.OnGround {
		t.Error("Expected player to be on ground initially")
	}
	
	if player.DoubleJump != DoubleJumpMax {
		t.Errorf("Expected DoubleJump=%d, got %d", DoubleJumpMax, player.DoubleJump)
	}
}

func TestPlayerMoveHorizontal(t *testing.T) {
	player := NewPlayer(0, 0)
	
	// 测试向右移动
	player.MoveHorizontal(1)
	if player.VX != PlayerSpeed {
		t.Errorf("Expected VX=%f when moving right, got %f", PlayerSpeed, player.VX)
	}
	
	// 测试向左移动
	player.MoveHorizontal(-1)
	if player.VX != -PlayerSpeed {
		t.Errorf("Expected VX=%f when moving left, got %f", -PlayerSpeed, player.VX)
	}
}

func TestPlayerJump(t *testing.T) {
	player := NewPlayer(0, 0)
	
	// 测试在地面上跳跃
	player.Jump()
	if player.VY != -JumpPower {
		t.Errorf("Expected VY=%f when jumping, got %f", -JumpPower, player.VY)
	}
	
	if player.OnGround {
		t.Error("Expected player not to be on ground after jumping")
	}
	
	if player.DoubleJump != DoubleJumpMax-1 {
		t.Errorf("Expected DoubleJump=%d after first jump, got %d", DoubleJumpMax-1, player.DoubleJump)
	}
	
	// 测试二段跳
	player.Jump()
	if player.VY != -JumpPower {
		t.Errorf("Expected VY=%f when double jumping, got %f", -JumpPower, player.VY)
	}
	
	// 测试三段跳（应该无效）
	initialVY := player.VY
	player.Jump()
	if player.VY != initialVY {
		t.Error("Expected no change in VY when trying to jump more than double jump max")
	}
}

func TestPlayerDash(t *testing.T) {
	player := NewPlayer(0, 0)
	
	// 测试冲刺
	player.Dash(100, 0) // 鼠标在右侧
	
	if !player.Dashing {
		t.Error("Expected player to be dashing")
	}
	
	if player.DashTimer != DashDuration {
		t.Errorf("Expected DashTimer=%d, got %d", DashDuration, player.DashTimer)
	}
	
	if player.VY != 0 {
		t.Error("Expected VY=0 when dashing")
	}
}

func TestPlayerUpdate(t *testing.T) {
	player := NewPlayer(0, 0)
	
	// 测试重力
	player.SetOnGround(false)
	initialY := player.Y
	player.Update()
	
	if player.Y <= initialY {
		t.Error("Expected player Y to increase due to gravity")
	}
	
	// 测试在地面上不应用重力
	player.SetOnGround(true)
	player.VY = 0
	initialY = player.Y
	player.Update()
	
	if player.Y != initialY {
		t.Error("Expected player Y to remain the same when on ground")
	}
	
	// 测试从空中落地时重置双跳
	player.DoubleJump = 0 // 设置为0
	player.SetOnGround(false) // 先确保在空中
	player.SetOnGround(true) // 然后落地
	if player.DoubleJump != DoubleJumpMax {
		t.Errorf("Expected DoubleJump=%d when landing, got %d", DoubleJumpMax, player.DoubleJump)
	}
}