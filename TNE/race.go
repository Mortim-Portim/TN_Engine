package TNE

import (

)

type Race struct {
	Entity
}
func (r *Race) Copy() (r2 *Race) {
	r2 = &Race{Entity:*r.Entity.Copy()}
	return
}