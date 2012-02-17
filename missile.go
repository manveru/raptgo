/*
Missile
-------
appearance: same as one of your Air\Air Missile
damage: 4 bars (2 bars in training mode)
projectile speed: average

Very similar to your own Air\Air Missile, even the damage is the same. If you
look carefully they don't quite move in the same way though. Your missiles
start out slower but accelerate a bit, while the enemy's version has a more
constant velocity. But overall the average speed is similar. Not usually too
troublesome except in a few cases.
*/

// Missile, an enemy weapon.
package main

var (
  MissileDamage   float64 = 2
  MissileRate     int64  = 20
  MissileVelocity = &Vector{X: 0, Y: 5.0}
  MissileCost     = 175600
  MissileSurface  *Surface
)

type Missile struct {
  position              *Vector
  velocity              *Vector
  shape                 *RectangleShape
  DrawablesId, ActorsId int64
}

func NewMissile(pos Vector) *Missile {
  position := &pos

  m := &Missile{
    position: position,
    velocity: MissileVelocity,
  }

  if MissileSurface == nil {
    MissileSurface = CreateImage("gfx/missile_top_down.png")
  }

  m.shape = NewRectangleShape(
    position,
    int16(MissileSurface.W),
    int16(MissileSurface.H),
  )

  Register(m)

  return m
}

func (m *Missile) Draw() { Draw(MissileSurface, m.position) }

func (m *Missile) Act() {
  m.position.Replace(m.position.Plus(m.velocity))

  if int(m.position.Y) > Height {
    Unregister(m)
    return
  }

  if m.shape.Collides(Player.shape) {
    Player.TakeDamage(MissileDamage)
    Unregister(m)
  }
}
