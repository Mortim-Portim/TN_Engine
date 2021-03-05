package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
)

var (
	Width, Height float64
	ScreenScale   = 32.0
)

func main() {
	GE.Init("", 30)
	x, y := ebiten.ScreenSizeInFullscreen()
	Width, Height = float64(x), float64(y)
	ps := GE.GetNewParticleSystem(10)
	sprites, err := GE.LoadEbitenImg("./res/arrow.png")
	GE.ShitImDying(err)
	psAnim := GE.GetAnimation(0, 0, 1, 1, 16, 6, sprites)
	pf := GE.GetNewParticleFactory(100, 30, psAnim)
	//ps.Add(pf.GetNewRandom(0, 1.0, 50, 50, 1, 1))

	game := &TestGame{ps, pf, 0}
	ebiten.SetWindowSize(int(Width), int(Height))
	ebiten.SetWindowTitle("Particles")
	ebiten.SetMaxTPS(30)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetScreenTransparent(true)
	//ebiten.SetWindowFloating(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type TestGame struct {
	ps           *GE.ParticleSystem
	pf           *GE.ParticleFactory
	frameCounter int
}

func (g *TestGame) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.ps.Add(g.pf.GetNewRandom(g.frameCounter, 1.0, float64(x)/ScreenScale, float64(y)/ScreenScale, 1, 1))
	}
	g.ps.Update(g.frameCounter)
	g.frameCounter++
	return nil
}
func (g *TestGame) Draw(screen *ebiten.Image) {
	g.ps.DrawOnPlainScreen(screen, 0, 0, ScreenScale)
}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(Width), int(Height)
}
