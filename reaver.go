// MG21C Twin Reaver Machine Guns
package main

var (
  ReaverDamage   float64 = 2
  ReaverRate     uint64  = 5
  ReaverVelocity = &Vector{X: 0, Y: 5.0}
  ReaverCost     = 12000
  ReaverSurface  [4]*Surface
  ReaverSize     int16 = 3
)

type Reaver struct {
  position              *Vector
  velocity              *Vector
  shape                 *CircleShape
  DrawablesId, ActorsId uint64
}

func NewReaver(pos Vector) *Reaver {
  position := &pos

  reaver := &Reaver{
    position: position,
    shape:    NewCircleShape(position, ReaverSize),
  }

  if ReaverSurface[0] == nil {
    ReaverSurface = [4]*Surface{
      CreateImage("gfx/reaver0.png"),
      CreateImage("gfx/reaver1.png"),
      CreateImage("gfx/reaver2.png"),
      CreateImage("gfx/reaver3.png"),
    }
  }

  Register(reaver)

  return reaver
}

func (r *Reaver) Draw() {
  surface := ReaverSurface[FrameTicks%4]
  Draw(surface, r.position)
}

func (r *Reaver) Act() {
  r.position.Replace(r.position.Minus(ReaverVelocity))

  if r.position.Y < 0 {
    Unregister(r)
    return
  }

  for _, enemy := range Enemies {
    if r.shape.Collides(enemy.shape) {
      enemy.TakeDamage(ReaverDamage)
      Unregister(r)
      return
    }
  }
}
