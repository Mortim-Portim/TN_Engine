package TNE

import (
	"fmt"
	"math/rand"

	"github.com/mortim-portim/GraphEng/GE"
)

type EnupPassive struct {
	Speed int
	nodes [][2]int
}

func (enup *EnupPassive) Update(e *Entity, world *World) {
	if enup.nodes == nil {
		enup.NewRoute(e, world)
	}
}

func (enup *EnupPassive) NewRoute(e *Entity, world *World) {
	var x, y int
	bnds := e.Eobj.Hitbox.Bounds()
	for {
		x = rand.Intn(10)
		y = rand.Intn(10 - x)

		if rand.Intn(2) == 1 {
			x *= -1
		}

		if rand.Intn(2) == 1 {
			y *= -1
		}

		if world.Structure.Collides(float64(x), float64(y), bnds.X, bnds.Y) {
			continue
		}

		fmt.Println([2]int{int(e.xPos), int(e.yPos)})

		enup.nodes = GE.FindPathMat(world.Structure.ObjMat, [2]int{int(e.xPos), int(e.yPos)}, [2]int{10, 10}, false)

		if len(enup.nodes) != 0 {
			break
		}

		fmt.Println(enup.nodes)
		panic("Hallo")
	}

	fmt.Println(enup.nodes)
	panic("Hi")
}
