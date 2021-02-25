package TNE

const ABForH1 = 0.707106781
const INVALID_DIR = -1
const (
	ENTITY_ORIENTATION_L = iota
	ENTITY_ORIENTATION_R
	ENTITY_ORIENTATION_U
	ENTITY_ORIENTATION_D
	ENTITY_ORIENTATION_LU
	ENTITY_ORIENTATION_RD
	ENTITY_ORIENTATION_RU
	ENTITY_ORIENTATION_LD
)

func GetNewDirection() *Direction {
	return &Direction{ID: INVALID_DIR}
}
func GetNewRandomDirection() (d *Direction) {
	d = &Direction{ID: RandomInt(ENTITY_ORIENTATION_L, ENTITY_ORIENTATION_LD)}
	d.FromID()
	return
}
func (d *Direction) Equals(d2 *Direction) bool {
	if d.ID == d2.ID {
		return true
	}
	return false
}
func (d *Direction) ToByte() byte {
	return byte(d.ID)
}
func (d *Direction) FromByte(b byte) {
	d.ID = int(b)
	d.FromID()
}
func (d *Direction) Copy() *Direction {
	return &Direction{d.ID, d.L, d.R, d.U, d.D, d.Dx, d.Dy}
}

type Direction struct {
	ID         int
	L, R, U, D bool
	Dx, Dy     float64
}

func (d *Direction) Print() string {
	if d.ID == ENTITY_ORIENTATION_LU {
		return "Left Up"
	} else if d.ID == ENTITY_ORIENTATION_RD {
		return "Right Down"
	} else if d.ID == ENTITY_ORIENTATION_RU {
		return "Right Up"
	} else if d.ID == ENTITY_ORIENTATION_LD {
		return "Left Down"
	} else if d.ID == ENTITY_ORIENTATION_L {
		return "Left"
	} else if d.ID == ENTITY_ORIENTATION_R {
		return "Right"
	} else if d.ID == ENTITY_ORIENTATION_U {
		return "Up"
	} else if d.ID == ENTITY_ORIENTATION_D {
		return "Down"
	}
	return "Invalid Direction"
}
func (d *Direction) IsValid() bool {
	return d.ID != INVALID_DIR
}
func (d *Direction) FromID() *Direction {
	if d.ID == ENTITY_ORIENTATION_LU {
		d.L = true
		d.U = true
		d.R = false
		d.D = false
		d.Dx = -ABForH1
		d.Dy = -ABForH1
	} else if d.ID == ENTITY_ORIENTATION_RD {
		d.L = false
		d.U = false
		d.R = true
		d.D = true
		d.Dx = ABForH1
		d.Dy = ABForH1
	} else if d.ID == ENTITY_ORIENTATION_RU {
		d.L = false
		d.U = true
		d.R = true
		d.D = false
		d.Dx = ABForH1
		d.Dy = -ABForH1
	} else if d.ID == ENTITY_ORIENTATION_LD {
		d.L = true
		d.U = false
		d.R = false
		d.D = true
		d.Dx = -ABForH1
		d.Dy = ABForH1
	} else if d.ID == ENTITY_ORIENTATION_L {
		d.L = true
		d.U = false
		d.R = false
		d.D = false
		d.Dx = -1
		d.Dy = 0
	} else if d.ID == ENTITY_ORIENTATION_R {
		d.L = false
		d.U = false
		d.R = true
		d.D = false
		d.Dx = 1
		d.Dy = 0
	} else if d.ID == ENTITY_ORIENTATION_U {
		d.L = false
		d.U = true
		d.R = false
		d.D = false
		d.Dx = 0
		d.Dy = -1
	} else if d.ID == ENTITY_ORIENTATION_D {
		d.L = false
		d.U = false
		d.R = false
		d.D = true
		d.Dx = 0
		d.Dy = 1
	} else {
		d.ID = INVALID_DIR
	}
	return d
}
func (d *Direction) FromKeys() *Direction {
	if d.L && d.R || d.U && d.D {
		d.ID = INVALID_DIR
	} else if d.L && d.U {
		d.ID = ENTITY_ORIENTATION_LU
		d.Dx = -ABForH1
		d.Dy = -ABForH1
	} else if d.R && d.D {
		d.ID = ENTITY_ORIENTATION_RD
		d.Dx = ABForH1
		d.Dy = ABForH1
	} else if d.R && d.U {
		d.ID = ENTITY_ORIENTATION_RU
		d.Dx = ABForH1
		d.Dy = -ABForH1
	} else if d.L && d.D {
		d.ID = ENTITY_ORIENTATION_LD
		d.Dx = -ABForH1
		d.Dy = ABForH1
	} else if d.L {
		d.ID = ENTITY_ORIENTATION_L
		d.Dx = -1
		d.Dy = 0
	} else if d.R {
		d.ID = ENTITY_ORIENTATION_R
		d.Dx = 1
		d.Dy = 0
	} else if d.U {
		d.ID = ENTITY_ORIENTATION_U
		d.Dx = 0
		d.Dy = -1
	} else if d.D {
		d.ID = ENTITY_ORIENTATION_D
		d.Dx = 0
		d.Dy = 1
	} else {
		d.ID = INVALID_DIR
	}
	return d
}
func (d *Direction) FromDelta() *Direction {
	if d.Dx < 0 && d.Dy < 0 {
		d.ID = ENTITY_ORIENTATION_LU
		d.L = true
		d.U = true
		d.R = false
		d.D = false
	} else if d.Dx > 0 && d.Dy > 0 {
		d.ID = ENTITY_ORIENTATION_RD
		d.L = false
		d.U = false
		d.R = true
		d.D = true
	} else if d.Dx > 0 && d.Dy < 0 {
		d.ID = ENTITY_ORIENTATION_RU
		d.L = false
		d.U = true
		d.R = true
		d.D = false
	} else if d.Dx < 0 && d.Dy > 0 {
		d.ID = ENTITY_ORIENTATION_LD
		d.L = true
		d.U = false
		d.R = false
		d.D = true
	} else if d.Dx == -1 && d.Dy == 0 {
		d.ID = ENTITY_ORIENTATION_L
		d.L = true
		d.U = false
		d.R = false
		d.D = false
	} else if d.Dx == 1 && d.Dy == 0 {
		d.ID = ENTITY_ORIENTATION_R
		d.L = false
		d.U = false
		d.R = true
		d.D = false
	} else if d.Dx == 0 && d.Dy == -1 {
		d.ID = ENTITY_ORIENTATION_U
		d.L = false
		d.U = true
		d.R = false
		d.D = false
	} else if d.Dx == 0 && d.Dy == 1 {
		d.ID = ENTITY_ORIENTATION_D
		d.L = false
		d.U = false
		d.R = false
		d.D = true
	} else {
		d.ID = INVALID_DIR
	}
	return d
}
