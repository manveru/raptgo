package main

type Goodie struct {
  position              *Vector
  velocity              *Vector
  sprites               []*Surface
  shape                 *CircleShape
  onHit                 func()
  DrawablesId, ActorsId int64
}

func NewGoodie(pos, vel *Vector, size int16, sprites []*Surface, onHit func()) *Goodie {
  goodie := &Goodie{
    position: pos,
    velocity: vel,
    sprites:  sprites,
    onHit:    onHit,
    shape:    NewCircleShape(pos, size),
  }

  Register(goodie) // we have shape and surface now.

  return goodie
}

func (g *Goodie) Draw() {
  n := int64(len(g.sprites))
  surface := g.sprites[FrameTicks%n]
  Draw(surface, g.position)
}

func (g *Goodie) Act() {
  if g.shape.Collides(Player.shape) {
    g.onHit()
    Unregister(g)
  } else {
    g.position.Replace(g.position.Plus(g.velocity))
  }
}

var ShieldpackSurface [6]*Surface

func MakeShieldpackSurface() {
  ShieldpackSurface = [6]*Surface{
    CreateImage("gfx/shieldpack_0.png"),
    CreateImage("gfx/shieldpack_1.png"),
    CreateImage("gfx/shieldpack_2.png"),
    CreateImage("gfx/shieldpack_3.png"),
    CreateImage("gfx/shieldpack_4.png"),
    CreateImage("gfx/shieldpack_5.png"),
  }

  // ShieldpackSurface = &Surface{RenderText("*", Blue)}
}

func NewShieldpack(pos, vel Vector) *Goodie {
  return NewGoodie(
    &pos, &vel, 5, ShieldpackSurface[:],
    func() {
      if Player.shields < Player.maxShields {
        Player.shields += 1
      } else {
        Player.shield = 100
      }
    },
  )
}

var MoneybagSurface [1]*Surface

func MakeMoneybagSurface() {
  MoneybagSurface = [1]*Surface{
    &Surface{RenderText("$", Yellow)},
  }
}

func NewMoneybag(pos, vel Vector) *Goodie {
  return NewGoodie(
    &pos, &vel, 5, MoneybagSurface[:],
    func() { Player.money += 1000 },
  )
}
