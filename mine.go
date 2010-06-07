// Enemy Mine
package main

var (
  MineDamage   float64 = 16
  MineRate     uint64  = 100
  MineVelocity = &Vector{X: 0, Y: 1.0}
  MineCost     = 12000
  MineSurface  [5]*Surface
  MineSize     int16 = 3
)

type Mine struct {
  position              *Vector
  velocity              *Vector
  shape                 *CircleShape
  DrawablesId, ActorsId uint64
}

func NewMine(pos Vector) *Mine {
  position := &pos

  mine := &Mine{
    position: position,
    shape:    NewCircleShape(position, MineSize),
  }

  if MineSurface[0] == nil {
    MineSurface = [5]*Surface{
      CreateImage("gfx/mine_0.png"),
      CreateImage("gfx/mine_1.png"),
      CreateImage("gfx/mine_2.png"),
      CreateImage("gfx/mine_3.png"),
      CreateImage("gfx/mine_4.png"),
    }
  }

  Register(mine)

  return mine
}

func (m *Mine) Draw() {
  tick := FrameTicks % (4 * 2)
  if tick > 4 {
    tick = 4 - (tick - 4)
  }
  surface := MineSurface[tick]
  Draw(surface, m.position)
}

func (m *Mine) Act() {
  m.position.Replace(m.position.Plus(MineVelocity))

  if m.position.Y > float64(Height) {
    Unregister(m)
    return
  }

  if m.shape.Collides(Player.shape) {
    Player.TakeDamage(MineDamage)
    Unregister(m)
    return
  }
}
