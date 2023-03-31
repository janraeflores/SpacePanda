package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	playerSpeed  = 5
)

type Game struct {
	BackgroundImage *ebiten.Image
	PlayerImage     *ebiten.Image
	PlayerX         float64
	PlayerY         float64
	ScrollX         float64
}

type Obstacle struct {
	ObstacleImage *ebiten.Image
	X             float64
	Y             float64
}

var obstacles []Obstacle

func createObstacles() {
	// Load obstacle image
	obstacleImage, _, err := ebitenutil.NewImageFromFile("cookieroid.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	// Create obstacles
	obstacle, err := resizeImage(obstacleImage, 400, 450)
	if err != nil {
		panic(err)
	}
	obstacles = []Obstacle{
		{ObstacleImage: obstacle, X: 1000, Y: -120},
		{ObstacleImage: obstacle, X: 2000, Y: 400},
		{ObstacleImage: obstacle, X: 3000, Y: 500},
		{ObstacleImage: obstacle, X: 3700, Y: 100},
	}
}

func resizeImage(img *ebiten.Image, width, height int) (*ebiten.Image, error) {
	newImg, err := ebiten.NewImage(width, height, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	scaleX := float64(width) / float64(img.Bounds().Max.X)
	scaleY := float64(height) / float64(img.Bounds().Max.Y)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scaleX, scaleY)
	newImg.DrawImage(img, opts)
	return newImg, nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	// Move player with arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.PlayerY -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.PlayerY += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.PlayerX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.PlayerX += 5
	}

	// Scroll background and player (if necessary)
	if g.PlayerX < screenWidth/2 && g.ScrollX > 0 {
		g.ScrollX -= playerSpeed
	} else if g.PlayerX > screenWidth/2 && g.ScrollX < float64(g.BackgroundImage.Bounds().Max.X)-screenWidth {
		g.ScrollX += playerSpeed
	}

	// Check for collision with obstacles
	for _, obstacle := range obstacles {
		if g.PlayerX+float64(g.PlayerImage.Bounds().Max.X) > obstacle.X && g.PlayerX < obstacle.X+float64(obstacle.ObstacleImage.Bounds().Max.X) &&
			g.PlayerY+float64(g.PlayerImage.Bounds().Max.Y) > obstacle.Y && g.PlayerY < obstacle.Y+float64(obstacle.ObstacleImage.Bounds().Max.Y) {

			// If player collides with obstacle, game restarts
			g.PlayerX = screenWidth / 2
			g.PlayerY = screenHeight / 2
			g.ScrollX = 0
			createObstacles()
			break
		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-g.ScrollX, 0)
	screen.DrawImage(g.BackgroundImage, op)

	// Draw player image
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.PlayerX-g.ScrollX, g.PlayerY)
	screen.DrawImage(g.PlayerImage, op)

	// Draw obstacles
	for _, obstacle := range obstacles {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(obstacle.X-g.ScrollX, obstacle.Y)
		screen.DrawImage(obstacle.ObstacleImage, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Load backgorund image
	backgroundImage, _, err := ebitenutil.NewImageFromFile("space2.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	// Load player image
	playerImage, _, err := ebitenutil.NewImageFromFile("panda.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	createObstacles()

	// Create game instance
	game := &Game{
		BackgroundImage: backgroundImage,
		PlayerImage:     playerImage,
		PlayerX:         screenWidth / 2,
		PlayerY:         screenHeight / 2,
	}

	// Start the game loop
	ebiten.RunGame(game)
}
