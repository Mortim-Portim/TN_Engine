package main

import (
	"flag"
	"image"
	"log"
	"os"
	"os/signal"

	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"

	"github.com/mortim-portim/TN_Engine/res"
)

const FPS = 60

var (
	Width, Height float64
	ScreenScale   = 32.0
)
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to file")

func Close() {
	GE.StopProfiling(cpuprofile, memprofile)
	log.Fatal("Termination")
}

func main() {
	flag.Parse()
	GE.Init("", FPS)
	GE.StartProfiling(cpuprofile)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		Close()
		return
	}()
	x, y := ebiten.ScreenSizeInFullscreen()
	Width, Height = float64(x), float64(y)

	sprites, err := GE.LoadEbitenImgFromBytes(res.SPARK_IMG)
	GE.ShitImDying(err)
	psAnim := GE.GetAnimation(0, 0, 1, 1, 15, 6, sprites)
	pf := GE.GetNewParticleFactory(100, FPS, psAnim)
	ps := GE.GetNewParticleSystem(10, pf)

	game := &TestGame{ps, pf, 0, 0, 0}
	ebiten.SetWindowSize(int(Width), int(Height))
	ebiten.SetWindowTitle("Fireball")
	ebiten.SetWindowIcon([]image.Image{sprites})
	ebiten.SetMaxTPS(FPS)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowDecorated(false)
	if err := ebiten.RunGame(game); err != nil {
		Close()
		log.Fatal(err)
	}
	Close()
}

type TestGame struct {
	ps           *GE.ParticleSystem
	pf           *GE.ParticleFactory
	lastX, lastY int
	frameCounter int
}

func (g *TestGame) Update() error {
	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && (x != g.lastX || y != g.lastY) {
		g.ps.Spawn(1, g.frameCounter, 100, (&GE.Vector{float64(g.lastX) - float64(x), float64(g.lastY) - float64(y), 0}).Mul(0.004), 0, float64(x)/ScreenScale, float64(y)/ScreenScale, 1, 1)
	}
	g.lastX = x
	g.lastY = y
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
