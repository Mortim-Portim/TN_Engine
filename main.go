package main 

import (
	"marvin/TN_Engine/TNE"
	"marvin/GraphEng/GE"
	
	"github.com/hajimehoshi/ebiten"
	"fmt"
	"time"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 1600
	screenHeight = 900
	FPS = 30
)
func StartGame(g ebiten.Game) {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("TN_Engine Test")
	//ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(FPS)
	//ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
	GE.CloseLogFile()
}


type TestGame struct {
	character *TNE.Creature
	frame int
}
func (g *TestGame) Init(screen *ebiten.Image) {}

var timeTaken int64
func (g *TestGame) Update(screen *ebiten.Image) error {
	startTime := time.Now()
	if g.frame%1 == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.character.ChangeOrientation(0)
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.character.ChangeOrientation(1)
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.character.ChangeOrientation(2)
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.character.ChangeOrientation(3)
		}
	}
	if g.frame%(FPS/2) == 0 {
		g.character.Move(1, FPS/4)
	}
	g.character.Update(g.character, nil)
	g.character.Draw(screen, 255, 0, 0, 0, 0, 100)
	
	g.frame ++
	timeTaken = time.Now().Sub(startTime).Milliseconds()
	fps := ebiten.CurrentTPS()
	msg := fmt.Sprintf(`TPS: %0.2f, Updating took: %v at frame %v`, fps, timeTaken, g.frame-1)
	ebitenutil.DebugPrint(screen, msg)
	GE.LogToFile(msg+"\n")
	fmt.Println(msg)
	return nil
}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func main() {
	GE.Init("")
	GE.SetLogFile("./res/log.txt")
	time.Sleep(time.Second*2)
	
	game := &TestGame{nil, 0}
	
	cf, err := TNE.GetCreatureFactory("./res/creatures/", &game.frame, 3)
	GE.ShitImDying(err)
	
	prepStart := time.Now()
	go cf.Prepare()
	fmt.Println("Preparing took: ", time.Now().Sub(prepStart))
	
	getStart := time.Now()
	c := cf.Get(0)
	//GE.ShitImDying(err)
	fmt.Println("Getting took: ", time.Now().Sub(getStart))
	
	c.RegiserUpdateFunc(func(e TNE.EntityI, world *TNE.World) {
		//fmt.Println("Updating creature, frame: ", frame)
	})
	game.character = c
	
	game.Init(nil)
	StartGame(game)
}

