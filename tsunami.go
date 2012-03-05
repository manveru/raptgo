// "Tsunami" Pulse Cannon
package main

var (
  TsunamiDamage   float64 = 5
  TsunamiRate     int64   = 10
  TsunamiVelocity         = &Vector{X: 0, Y: 10.0}
  TsunamiCost             = 725000
  TsunamiSurface  [14]*Surface
)

type Tsunami struct {
  position              *Vector
  velocity              *Vector
  shape                 *RectangleShape
  DrawablesId, ActorsId int64
  drawFlip              bool
}

func NewTsunami(pos Vector) *Tsunami {
  position := &pos

  t := &Tsunami{
    position: position,
    velocity: &Vector{0, 0},
  }

  if TsunamiSurface[0] == nil {
    TsunamiSurface = [14]*Surface{
      CreateImage("gfx/tsunami_00.png"),
      CreateImage("gfx/tsunami_01.png"),
      CreateImage("gfx/tsunami_02.png"),
      CreateImage("gfx/tsunami_03.png"),
      CreateImage("gfx/tsunami_04.png"),
      CreateImage("gfx/tsunami_05.png"),
      CreateImage("gfx/tsunami_06.png"),
      CreateImage("gfx/tsunami_07.png"),
      CreateImage("gfx/tsunami_08.png"),
      CreateImage("gfx/tsunami_09.png"),
      CreateImage("gfx/tsunami_10.png"),
      CreateImage("gfx/tsunami_11.png"),
      CreateImage("gfx/tsunami_12.png"),
      CreateImage("gfx/tsunami_13.png"),
    }
  }

  t.shape = NewRectangleShape(
    position,
    int16(TsunamiSurface[0].W),
    int16(TsunamiSurface[0].H),
  )

  Register(t)

  return t
}

func (t *Tsunami) Draw() {
  tick := FrameTicks % (13 * 2)
  if tick > 13 {
    tick = 13 - (tick - 13)
  }
  surface := TsunamiSurface[tick]
  Draw(surface, t.position)
}

func (t *Tsunami) Act() {
  // accelerate until we reach TsunamiVelocity
  vel := t.velocity
  if vel.Length() < TsunamiVelocity.Length() {
    vel = vel.Plus(TsunamiVelocity.DivideNum(10.0))
    t.velocity.Replace(vel)
  }
  t.position.Replace(t.position.Minus(vel))

  if t.position.Y < 0 {
    Unregister(t)
    return
  }

  for _, enemy := range Enemies {
    if t.shape.Collides(enemy.shape) {
      enemy.TakeDamage(TsunamiDamage)
      Unregister(t)
      return
    }
  }
}
