package entity

import (
	"testing"
)

func TestNewCamera(t *testing.T) {
	camera := NewCamera(100, 50)
	
	if camera.X != 100 {
		t.Errorf("Expected X=100, got %f", camera.X)
	}
	
	if camera.Y != 50 {
		t.Errorf("Expected Y=50, got %f", camera.Y)
	}
	
	if camera.FollowSpeed != 0.1 {
		t.Errorf("Expected FollowSpeed=0.1, got %f", camera.FollowSpeed)
	}
}

func TestCameraSetTarget(t *testing.T) {
	camera := NewCamera(0, 0)
	
	camera.SetTarget(100, 50)
	
	if camera.TargetX != 100 {
		t.Errorf("Expected TargetX=100, got %f", camera.TargetX)
	}
	
	if camera.TargetY != 50 {
		t.Errorf("Expected TargetY=50, got %f", camera.TargetY)
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
	camera := NewCamera(100, 100)
	camera.SetScreenSize(800, 600)
	
	screenX, screenY := camera.WorldToScreen(150, 150)
	
	// Expected: (150 - 100 + 400, 150 - 100 + 300) = (450, 350)
	if screenX != 450 {
		t.Errorf("Expected screenX=450, got %f", screenX)
	}
	
	if screenY != 350 {
		t.Errorf("Expected screenY=350, got %f", screenY)
	}
}

func TestCameraScreenToWorld(t *testing.T) {
	camera := NewCamera(100, 100)
	camera.SetScreenSize(800, 600)
	
	worldX, worldY := camera.ScreenToWorld(450, 350)
	
	// Expected: (450 + 100 - 400, 350 + 100 - 300) = (150, 150)
	if worldX != 150 {
		t.Errorf("Expected worldX=150, got %f", worldX)
	}
	
	if worldY != 150 {
		t.Errorf("Expected worldY=150, got %f", worldY)
	}
}

func TestCameraUpdate(t *testing.T) {
	camera := NewCamera(0, 0)
	camera.SetTarget(100, 100)
	
	// Run several updates to approach the target
	for i := 0; i < 100; i++ {
		camera.Update()
	}
	
	// After many updates, camera should be very close to target
	if abs(camera.X-100) > 0.1 {
		t.Errorf("Expected X to be close to 100, got %f", camera.X)
	}
	
	if abs(camera.Y-100) > 0.1 {
		t.Errorf("Expected Y to be close to 100, got %f", camera.Y)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}