package TNE

import (

)

//SHOULD contain information about the races stats
type Race struct {
	Entity
}
func (r *Race) Copy() (r2 *Race) {
	r2 = &Race{Entity:*r.Entity.Copy()}
	return
}