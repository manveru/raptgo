// AIM-31 "Mauler" Air/Air Missile
// It sucks by design...
package main

var (
  MaulerDamage   float64 = 4
  MaulerRate     int64  = 24
  MaulerVelocity = &Vector{X: 0, Y: -4.0}
  MaulerCost     = 63500
  MaulerSurface  *Surface
  MaulerSize     int16 = 3
  MaulerSpread   = 50
)

type Mauler struct {
  position              *Vector
  velocity              *Vector
  speed                 float64
  shape                 *CircleShape
  spread                int
  DrawablesId, ActorsId int64
}

func NewMauler(position, velocity Vector) *Mauler {
  pos, vel := &position, &velocity
  mauler := &Mauler{
    position: pos,
    velocity: vel,
    speed:    0.0,
    spread:   0,
    shape:    NewCircleShape(pos, MaulerSize),
  }

  if MaulerSurface == nil {
    MaulerSurface = CreateImage("gfx/missile_top_up.png")
  }

  Register(mauler)

  return mauler
}

func (m *Mauler) Draw() { Draw(MaulerSurface, m.position) }

func (m *Mauler) Act() {
  if m.position.Y < 0 {
    Unregister(m)
    return
  }

  for _, enemy := range Enemies {
    if m.shape.Collides(enemy.shape) {
      enemy.TakeDamage(MaulerDamage)
      Unregister(m)
      return
    }
  }

  if m.spread < MaulerSpread {
    m.spread++
    m.position.Replace(m.position.Plus(m.velocity))
    return
  }

  m.speed += 0.1
  m.velocity.Replace(MaulerVelocity.Normalize().MultiplyNum(m.speed))

  var closest *Enemy
  closestLength := float64(Width * Height)

  for _, enemy := range Enemies {
    length := m.position.Minus(enemy.position).Length()
    if closestLength > length {
      closestLength = length
      closest = enemy
    }
  }

  // adjust velocity slightly towards that enemy
  if closest != nil {
    pe := closest.position.Minus(m.position).Normalize().MultiplyNum(2.0)
    m.velocity.Replace(m.velocity.Plus(pe))
  }

  m.position.Replace(m.position.Plus(m.velocity))
}
