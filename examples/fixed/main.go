package main

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	animator "github.com/ninedraft/ebiten-animator"
)

const runnerSize = 32

func main() {
	runnerBytes := bytes.NewReader(images.Runner_png)
	runner, _, err := ebitenutil.NewImageFromReader(runnerBytes)
	if err != nil {
		panic("parsing runner texture: " + err.Error())
	}

	frames := make([]*ebiten.Image, 0, 8)
	for i := range cap(frames) {
		frame := runner.SubImage(image.Rect(i*runnerSize, runnerSize, i*runnerSize+runnerSize, 2*runnerSize)).(*ebiten.Image)
		frames = append(frames, frame)
	}

	animation := animator.NewAnimationFixed(time.Second/8, frames...)

	game := &Game{
		w: 300, h: 300,
		runner:   runner,
		animator: animator.New(animation),
	}

	ebiten.RunGame(game)
}

type Game struct {
	animator *animator.Animator
	runner   *ebiten.Image
	w, h     float64
}

func (game *Game) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(game.w/2, game.h/2)

	screen.Fill(color.White)
	game.animator.Draw(
		screen,
		&ebiten.DrawImageOptions{
			GeoM: geom,
		})
}

func (game *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	game.animator.Update()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		log.Printf("restarting animator")
		game.animator.Restart()
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		log.Printf("stopping animator")
		game.animator.Stop()
	}

	return nil
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(game.w), int(game.h)
}
