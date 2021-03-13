package TNE

import (
	"fmt"
	"testing"

	"github.com/mortim-portim/GraphEng/GE"
)

func TestSyncAction(t *testing.T) {
	sAs := NewActionStack()
	fmt.Printf("%p: %v\n", sAs, sAs)
	sAs2 := sAs.Copy()
	fmt.Printf("%p: %v\n", sAs2, sAs2)

	fc := 0
	testEnt, err := LoadEntity("../../TerraNomina_Client/res/Entities/Goblin", &fc)
	GE.ShitImDying(err)
	testEnt.SetMiddle(10.324, 34.243)
	testEnt.setIntPos()
	testEnt.actions = sAs
	fmt.Println(testEnt.Print())
	fmt.Println(testEnt.Actions().Print())
	testEnt.MoveLengthAndFrame(1.2, 10)
	fmt.Println(testEnt.Actions().Print())
	testEnt.ChangeOrientation(GetNewRandomDirection())
	fmt.Println(testEnt.Actions().Print())
	testEnt.KeepMoving(true)
	fmt.Println(testEnt.Actions().Print())
	testEnt.SetAnim(3)
	fmt.Println(testEnt.Actions().Print())
	testEnt.AddPos()
	fmt.Println(testEnt.Actions().Print())

	bs := testEnt.Actions().GetAll()
	fmt.Println(testEnt.Actions().Print())
	testEnt.Actions().Reset()
	fmt.Println(testEnt.Actions().Print())
	fmt.Println(testEnt.Print())

	testEnt2, err := LoadEntity("../../TerraNomina_Client/res/Entities/Goblin", &fc)
	GE.ShitImDying(err)
	testEnt2.actions = sAs2
	fmt.Println(testEnt2.Print())
	fmt.Println(testEnt2.Actions().Print())
	testEnt2.Actions().AppendAndApply(bs, testEnt2, nil)
	fmt.Println(testEnt2.Print())
	fmt.Println(testEnt2.Actions().Print())
}
