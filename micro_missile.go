// Micro-Missile
package main

var (
  MicroMissileDamage   float64 = 2
  MicroMissileRate     int64  = 8
  MicroMissileVelocity = &Vector{X: 0, Y: 5.0}
  MicroMissileCost     = 175600
  MicroMissileSurface  *Surface
  MicroMissileSize     int16 = 3
)

type MicroMissile struct {
  position              *Vector
  velocity              *Vector
  shape                 *CircleShape
  DrawablesId, ActorsId int64
}

func NewMicroMissile(pos Vector) *MicroMissile {
  position := &pos

  mm := &MicroMissile{
    position: position,
    velocity: &Vector{0, 0},
    shape:    NewCircleShape(position, MicroMissileSize),
  }

  if MicroMissileSurface == nil {
    MicroMissileSurface = CreateImage("gfx/micro_missile.png")
  }

  Register(mm)

  return mm
}

func (m *MicroMissile) Draw() { Draw(MicroMissileSurface, m.position) }

func (m *MicroMissile) Act() {
  // accelerate until we reach MicroMissileVelocity
  vel := m.velocity
  if vel.Length() < MicroMissileVelocity.Length() {
    vel = vel.Plus(MicroMissileVelocity.DivideNum(50.0))
    m.velocity.Replace(vel)
  }
  m.position.Replace(m.position.Minus(vel))

  if m.position.Y < 0 {
    Unregister(m)
    return
  }

  for _, enemy := range Enemies {
    if m.shape.Collides(enemy.shape) {
      enemy.TakeDamage(MicroMissileDamage)
      Unregister(m)
      return
    }
  }
}
