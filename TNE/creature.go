package TNE

import (
	
)


type Creature struct {
	Entity
}
func (c *Creature) Copy() (c2 *Creature) {
	c2 = &Creature{Entity:*c.Entity.Copy()}
	return
}