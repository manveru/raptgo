// OD55 "Odin" Laser Turret
package main

var (
  OdinDamage float64 = 5
  OdinRate   int64   = 20
  OdinCost           = 512850
  OdinSteps  float64 = 10
  OdinColor  uint32  = 0xaa6600
)

type Odin struct {
  from, to              *Vector
  victim                *Enemy
  stepLength            float64
  DrawablesId, ActorsId int64
}

func NewOdin() *Odin {
  odin := &Odin{}

  if odin.Seek() {
    Register(odin)
  }

  return odin
}

// draw a laser beam flying towards the enemy within the lifetime.
func (o *Odin) Draw() {
  from, to := o.from, o.to
  fx, fy, tx, ty := int16(from.X), int16(from.Y), int16(to.X), int16(to.Y)

  Screen.DrawLine(fx-3, fy, tx, ty, OdinColor)
  Screen.DrawLine(fx-2, fy, tx, ty, OdinColor)
  Screen.DrawLine(fx-1, fy, tx, ty, OdinColor)
  Screen.DrawLine(fx, fy, tx, ty, OdinColor)
  Screen.DrawLine(fx+1, fy, tx, ty, OdinColor)
  Screen.DrawLine(fx+2, fy, tx, ty, OdinColor)
  Screen.DrawLine(fx+3, fy, tx, ty, OdinColor)
}

func (o *Odin) Act() {
  if o.StepVector() {
    o.victim.TakeDamage(OdinDamage)
    Unregister(o)
    return
  }
}

func (o *Odin) StepVector() bool {
  from, to := o.to, o.victim.position
  // get a little closer.
  vector := to.Minus(from)
  length := vector.Length()

  if o.stepLength == 0.0 {
    o.stepLength = length / OdinSteps
  }

  if length > o.stepLength {
    o.from = from
    o.to = from.Plus(vector.Normalize().MultiplyNum(o.stepLength))
    return false
  }

  return true
}

func (o *Odin) Seek() bool {
  for _, enemy := range Enemies {
    o.victim = enemy
    o.to = Player.position
    o.StepVector()
    return true
  }

  return false
}
