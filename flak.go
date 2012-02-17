/*
Flak
----
appearance: orange bullet
damage: 2 bars (1 bar in training mode)
projectile speed: slow

The basic enemy weapon. Probably their version of your Twin Reaver Machine Guns
except way worse. Good thing too, otherwise you wouldn't be able to dodge it so
easily. Their bullets are much slower, and also bigger and easier to see.
*/
package main

var (
  FlakDamage   float64 = 2
  FlakRate     int64  = 40
  FlakVelocity = &Vector{X: 0, Y: -4.0}
  FlakSurface  *Surface
  FlakSize     int16 = 5
)

type Flak struct {
  position              *Vector
  velocity              *Vector
  shape                 *CircleShape
  DrawablesId, ActorsId int64
}

func NewFlak(pos Vector) *Flak {
  position := &pos

  flak := &Flak{
    position: position,
    shape:    NewCircleShape(position, FlakSize),
  }

  if FlakSurface == nil {
    FlakSurface = CreateImage("gfx/flak.png")
  }

  Register(flak)

  return flak
}

func (f *Flak) Draw() { Draw(FlakSurface, f.position) }

func (f *Flak) Act() {
  f.position.Replace(f.position.Minus(FlakVelocity))

  if f.position.Y > float64(Height) {
    Unregister(f)
    return
  }

  if f.shape.Collides(Player.shape) {
    Player.TakeDamage(FlakDamage)
    Unregister(f)
    return
  }
}
