// Blitz, my own creation just for fun
package main

import "math"
import "rand"

var (
  BlitzDamage float64 = 3
  BlitzRate   uint64  = 20
  BlitzCost   = 120000
  BlitzColor  uint32 = 0x0000ff
  BlitzJumps  = 10
  BlitzRange  float64 = 200
)

type Blitz struct {
  from, to              *Vector
  jumps                 int
  whatNow               func(*Blitz)
  victim                *Enemy
  lifeTime              uint64
  DrawablesId, ActorsId uint64
}

func NewBlitz(pos Vector) *Blitz {
  from := &pos

  blitz := &Blitz{
    from:     from,
    jumps:    BlitzJumps,
    to:       &Vector{X: from.X, Y: 0},
    whatNow:  (*Blitz).Seek,
    lifeTime: BlitzRate,
  }

  Register(blitz)

  return blitz
}

func (b *Blitz) Draw() {
  if b.whatNow == (*Blitz).Seek {
    cx, cy := int16(b.from.X), int16(b.from.Y)
    radius := BlitzRange

    for i := float64(0.0); i < 2*math.Pi; i += float64(rand.Float()) {
      x := cx + int16(radius*math.Sin(i))
      y := cy + int16(radius*math.Cos(i))
      Screen.DrawLine(int16(cx), int16(cy), x, y, BlitzColor)
    }
  } else {
    origin, vector := b.from, b.to
    x1, y1 := int16(origin.X), int16(origin.Y)
    x2, y2 := int16(vector.X), int16(vector.Y)
    Screen.DrawLine(x1, y1, x2, y2, BlitzColor)
    Screen.DrawLine(x1-1, y1-1, x2-1, y2+1, BlitzColor)
    Screen.DrawLine(x1+1, y1+1, x2+1, y2+1, BlitzColor)
  }
}

func (b *Blitz) Act() {
  if b.jumps < 1 || b.lifeTime < 1 {
    Unregister(b)
    return
  }

  b.lifeTime--
  b.whatNow(b)
}

func (b *Blitz) Seek() {
  circle := NewCircleShape(b.from, int16(BlitzRange))

  for _, enemy := range Enemies {
    if enemy == b.victim {
      continue
    }
    if enemy.shape.Collides(circle) {
      b.to.Replace(enemy.position)
      b.victim = enemy
      b.whatNow = (*Blitz).Jump
      return
    }
  }

  b.jumps--
}

func (b *Blitz) Jump() {
  b.jumps--
  b.victim.TakeDamage(BlitzDamage)
  b.from.Replace(b.victim.position)
  b.whatNow = (*Blitz).Seek
}
