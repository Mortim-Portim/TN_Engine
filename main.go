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
	character *TNE.Player
	world     *TNE.World
	rec  *GE.Recorder
	frame int
}
func (g *TestGame) Init(screen *ebiten.Image) {}

var timeTaken int64
func (g *TestGame) Update(screen *ebiten.Image) error {
	startTime := time.Now()
	moving := false
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.character.ChangeOrientation(0)
		if GE.IsKeyJustDown(ebiten.KeyA) {
			g.character.Move(1, FPS/4)
		}
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.character.ChangeOrientation(1)
		if GE.IsKeyJustDown(ebiten.KeyD) {
			g.character.Move(1, FPS/4)
		}
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.character.ChangeOrientation(2)
		if GE.IsKeyJustDown(ebiten.KeyW) {
			g.character.Move(1, FPS/4)
		}
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.character.ChangeOrientation(3)
		if GE.IsKeyJustDown(ebiten.KeyS) {
			g.character.Move(1, FPS/4)
		}
		moving = true
	}
	g.character.KeepMoving(moving)
	//g.character.UpdateAll(nil)
	//g.character.Draw(screen, 255, 0, 0, 0, 0, 100)
	//fmt.Println(g.character.Print())
	
	g.world.UpdateActivePlayer()
	g.world.UpdateDrawables()
	g.world.UpdateWorldStructure()
	g.world.Draw(screen)
	
	if ebiten.IsKeyPressed(ebiten.KeyC) && !g.rec.IsSaving() {
		g.rec.Save("./res/out")
	}
	g.rec.NextFrame(screen)
	
	g.frame ++
	timeTaken = time.Now().Sub(startTime).Milliseconds()
	fps := ebiten.CurrentTPS()
	msg := fmt.Sprintf(`TPS: %0.2f, Updating took: %v at frame %v, isMoving: %v`, fps, timeTaken, g.frame-1, moving)
	ebitenutil.DebugPrint(screen, msg)
	GE.LogToFile(msg+"\n")
	//fmt.Println(msg)
	return nil
}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func main() {
	GE.Init("")
	GE.SetLogFile("./res/log.txt")
	time.Sleep(time.Second*2)
	
	diceStart := time.Now()
	result := TNE.RollDice(100000, 20)
	fmt.Println("Rolling Dice took: ", time.Now().Sub(diceStart))
	fmt.Println(result)
	
	game := &TestGame{nil, nil, GE.GetNewRecorder(FPS*5, 1280, 720, FPS), 0}
	
	cf, err := TNE.GetCreatureFactory("./res/creatures/", &game.frame, 3)
	GE.ShitImDying(err)
	
	prepStart := time.Now()
	cf.Prepare()
	fmt.Println("Preparing took: ", time.Now().Sub(prepStart))
	
	getStart := time.Now()
	c := cf.Get(1)
	//GE.ShitImDying(err)
	fmt.Println("Getting took: ", time.Now().Sub(getStart))
	
	c.RegiserUpdateFunc(func(e TNE.EntityI, world *TNE.World) {
		//fmt.Println("Updating creature, frame: ", frame)
	})
	game.character = &TNE.Player{TNE.Race{c.Entity}}
	
	game.world = TNE.GetWorld(0,0,screenWidth,screenHeight, 16, 9, 4,6, cf, &game.frame, "./res/Worlds/TestWorld1", "TestMap1", "./res/Worlds/TestWorld1/tiles", "./res/Worlds/TestWorld1/structObjs")
	game.world.AddPlayer(game.character)
	err = game.world.SetActivePlayer(0)
	GE.ShitImDying(err)
	//game.world.Structure.SetMiddleSmooth(0, 0)
	
	fmt.Println(game.world.Print())
	
	game.Init(nil)
	StartGame(game)
}

