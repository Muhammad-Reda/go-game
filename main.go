package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenHeight       = 800
	screenWidth        = 600
	windowTitle        = "Go Game (Fullscreen/Maximized)"
	maxNPCs            = 1000
	collisionSoundFile = "hell-naw-dog.mp3"
)

type Player struct {
	Texture  rl.Texture2D
	Position rl.Vector2
	Velocity rl.Vector2
	Width    int32
	Height   int32
	Speed    float32
}

type NPC struct {
	Texture  rl.Texture2D
	Position rl.Vector2
	Velocity rl.Vector2
	Width    int32
	Height   int32
	Speed    float32
	MaxSpeed float32
}

func (p *Player) Spawn() {
	rl.DrawTexture(p.Texture, int32(p.Position.X), int32(p.Position.Y), rl.White)
}

func (p *Player) Move() {
	if rl.IsKeyDown(rl.KeyW) && p.Position.Y > 0 {
		p.Position.Y -= float32(p.Speed)
	}

	if rl.IsKeyDown(rl.KeyS) && p.Position.Y < float32(rl.GetScreenHeight()-int(p.Texture.Height)) {
		p.Position.Y += float32(p.Speed)
	}

	if rl.IsKeyDown(rl.KeyA) && p.Position.X > 0 {
		p.Position.X -= float32(p.Speed)
	}

	if rl.IsKeyDown(rl.KeyD) && p.Position.X < float32(rl.GetScreenWidth()-int(p.Texture.Width)) {
		p.Position.X += float32(p.Speed)
	}
}

func (n *NPC) Spawn() {
	rl.DrawTexture(n.Texture, int32(n.Position.X), int32(n.Position.Y), rl.White)
}

// Helper function to generate a random float value between min and max.
func randomFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

// CreateNewNPC creates a new NPC with random direction.
func CreateNewNPC(texture rl.Texture2D) NPC {
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	x := randomFloat(0, float32(screenWidth-int(texture.Width)))
	y := randomFloat(0, float32(screenHeight-int(texture.Height)))

	randXDirection := randomFloat(-1, 1)
	randYDirection := randomFloat(-1, 1)

	return NPC{
		Texture:  texture,
		Position: rl.Vector2{X: x, Y: y},
		Velocity: rl.Vector2{X: randXDirection, Y: randYDirection},
		Width:    texture.Width,
		Height:   texture.Height,
		Speed:    3.0,
		MaxSpeed: 10.0,
	}
}

// NewNPC creates a new NPC inheriting properties of the existing NPC.
func NewNPC(existingNPC *NPC) NPC {
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	x := randomFloat(0, float32(screenWidth-int(existingNPC.Texture.Width)))
	y := randomFloat(0, float32(screenHeight-int(existingNPC.Texture.Height)))

	randXDirection := randomFloat(-1, 1)
	randYDirection := randomFloat(-1, 1)

	return NPC{
		Texture:  existingNPC.Texture,
		Position: rl.Vector2{X: x, Y: y},
		Velocity: rl.Vector2{X: randXDirection, Y: randYDirection},
		Width:    existingNPC.Width,
		Height:   existingNPC.Height,
		Speed:    existingNPC.Speed,
		MaxSpeed: existingNPC.MaxSpeed,
	}
}

func (n *NPC) Move() bool {
	n.Position.X += n.Velocity.X
	n.Position.Y += n.Velocity.Y

	bounced := false
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()
	// Perbaiki batas layar dan gunakan Texture.Height untuk Y
	if n.Position.X+float32(n.Texture.Width) >= float32(screenWidth) {
		n.Velocity.X *= -1
		n.Position.X = float32(screenWidth) - float32(n.Texture.Width)
		bounced = true
	} else if n.Position.X <= 0 {
		n.Velocity.X *= -1
		n.Position.X = 0
		bounced = true
	}

	// Perbaiki batas layar dan gunakan Texture.Height untuk Y dan GetScreenHeight
	if n.Position.Y+float32(n.Texture.Height) >= float32(screenHeight) {
		n.Velocity.Y *= -1
		n.Position.Y = float32(screenHeight) - float32(n.Texture.Height)
		bounced = true
	} else if n.Position.Y <= 0 {
		n.Velocity.Y *= -1
		n.Position.Y = 0
		bounced = true
	}

	// Tingkatkan kecepatan jika memantul, tetapi jangan melebihi MaxSpeed
	if bounced && n.Speed < n.MaxSpeed {
		n.Speed += 1 // Tingkatkan kecepatan
		// Sesuaikan vektor kecepatan agar sesuai dengan kecepatan baru
		n.Velocity.X = normalizeAndScale(n.Velocity.X, n.Speed)
		n.Velocity.Y = normalizeAndScale(n.Velocity.Y, n.Speed)

		fmt.Printf("Kecepatan meningkat: %.2f\n", n.Speed) // Debug
	}

	return bounced
}

// Helper function to normalize a vector and scale it to the desired speed
func normalizeAndScale(value float32, newSpeed float32) float32 {
	// Normalisasikan vektor (hitung panjang)
	length := rl.Vector2Length(rl.Vector2{X: value, Y: 0}) // Hanya X karena kita menormalisasi per komponen
	if length > 0 {
		return (value / length) * newSpeed
	}
	return 0 // Jika panjangnya 0, kembalikan 0 untuk menghindari pembagian dengan nol
}

func main() {
	// Inisialisasi
	rl.InitWindow(800, 600, windowTitle) // Jendela awal - akan diubah ukurannya

	rl.SetWindowState(rl.FlagBorderlessWindowedMode) // Menghapus border window

	screenWidth := rl.GetScreenWidth()   // Get current screen width
	screenHeight := rl.GetScreenHeight() // Get current screen height

	rl.SetWindowSize(screenWidth, screenHeight)
	rl.MaximizeWindow() // Maximazed Window

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	// Load the collision sound effect
	collisionSound := rl.LoadSound(collisionSoundFile)
	defer rl.UnloadSound(collisionSound)

	// Load image
	catImage := rl.LoadImage("cat.png")
	defer rl.UnloadImage(catImage)

	robotImage := rl.LoadImage("robot.png")
	defer rl.UnloadImage(robotImage)

	// Resize robot image
	rl.ImageResize(robotImage, 100, 100)

	// Convert image to texture
	catTexture := rl.LoadTextureFromImage(catImage)
	defer rl.UnloadTexture(catTexture)

	robotTexture := rl.LoadTextureFromImage(robotImage)
	defer rl.UnloadTexture(catTexture)

	// Set FPS
	rl.SetTargetFPS(60)

	cat := Player{
		Texture:  catTexture,
		Position: rl.Vector2{X: 100, Y: 100},
		Velocity: rl.Vector2{X: 1, Y: 1},
		Width:    catTexture.Width,
		Height:   catTexture.Height,
		Speed:    3.0,
	}

	// Create a slice to hold all the NPCs
	npcs := []NPC{CreateNewNPC(robotTexture)} // Initial NPC

	// init catSpeed and robot exist state
	catMaxSpeed := float32(30.0)

	// Define the new background color (dark blue-grey as an example)
	newBackgroundColor := rl.Color{R: 20, G: 25, B: 40, A: 255} // Adjust these values!

	for !rl.WindowShouldClose() {
		cat.Move()

		// Init collision rect
		catRect := rl.NewRectangle(cat.Position.X, cat.Position.Y, float32(cat.Texture.Width), float32(cat.Texture.Height))

		// Loop through the NPCs
		for i := 0; i < len(npcs); i++ {
			// Init robot and collision
			robot := &npcs[i]
			robotRect := rl.NewRectangle(robot.Position.X, robot.Position.Y, float32(robot.Texture.Width), float32(robot.Texture.Height))

			// Move and spawn
			robot.Spawn()
			if bounced := robot.Move(); bounced {
				// If a robot bounced, create a *new* robot but only if we are below maxNPCs
				if len(npcs) < maxNPCs {
					fmt.Println("NPC Bounced! New NPC created with same properties")

					// Create a new NPC using same properties as existing NPC
					newRobot := NewNPC(robot)

					// Append the new robot to the npcs slice
					npcs = append(npcs, newRobot)
				} else {
					fmt.Println("Max NPC count reached. No new NPC created")
				}
			}

			// Check Collision
			if rl.CheckCollisionRecs(catRect, robotRect) {
				// Play the collision sound
				rl.PlaySound(collisionSound)

				if cat.Speed < catMaxSpeed {
					cat.Speed += 1
				}

				// Reset NPC hit code - We're always deleting a member of NPCs so always add it back unless we're already at maxNPCs
				npcs = append(npcs[:i], npcs[i+1:]...) // Delete the NPC from the slice
				i--                                    // Adjust the index so we don't skip one

				// Create a new NPC with default properties if we aren't already at maxNPCs
				if len(npcs) < maxNPCs {
					newRobot := CreateNewNPC(robotTexture) // Create the NPC with default properties
					npcs = append(npcs, newRobot)          // Append it to slice
					fmt.Println("robot hit!")              // Test hit!
				} else {
					fmt.Println("Max NPC count reached. No new NPC created")
				}
			}
		}

		// Start draw
		rl.BeginDrawing()

		// Clear background with the new color
		rl.ClearBackground(newBackgroundColor) // Use the new background color

		// Draw cat Texture
		cat.Spawn()

		// Draw text catSpeed
		rl.DrawText(fmt.Sprintf("Speed %d", int(cat.Speed)), 50, int32(rl.GetScreenHeight()-50), 40, rl.White)
		rl.DrawText(fmt.Sprintf("Number of NPCs: %d", len(npcs)), int32(rl.GetScreenWidth()-800), int32(rl.GetScreenHeight()-50), 40, rl.Red)

		rl.EndDrawing()
	}
}
