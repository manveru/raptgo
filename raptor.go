package main

import (
  . "fmt"
  "github.com/banthar/Go-SDL/sdl"
)

var (
  RaptorSurface       *Surface
  RaptorShieldSurface *Surface
  ShieldEnergy        float64 = 100.0
)

type Fighter interface {
  Shape() Shapeish
  TakeDamage(float64)
}

type Raptor struct {
  position   *Vector
  velocity   *Vector
  shape      *RectangleShape
  speed      float64
  weapon     string
  weapons    [255]string
  money      int
  maxShields int     "number of shields the raptor can hold"
  shields    int     "number of shields left, not including the one in use"
  shield     float64 "energy of the shield currently used"
  lastShot   int64
  lastDamage int64
}

func NewRaptor(position *Vector) *Raptor {
  raptor := &Raptor{
    position:   &Vector{320, 480},
    velocity:   &Vector{0, 0},
    speed:      10.0,
    maxShields: 5,
    shields:    1,
    shield:     ShieldEnergy,
  }

  raptor.weapons = [255]string{
    '1': "Reaver",
    '2': "Thor",
    '3': "Odin",
    '5': "Mauler",
    '9': "Tsunami",
    '0': "Deathray",
    '-': "Eclipse",
  }

  raptor.weapon = "Reaver"

  if RaptorSurface == nil {
    RaptorSurface = CreateImage("gfx/raptor.png")
    RaptorShieldSurface = CreateImage("gfx/raptor_shield.png")
  }
  raptor.shape = NewRectangleShape(
    raptor.position,
    int16(RaptorSurface.W),
    int16(RaptorSurface.H),
  )

  AddActor(raptor) // don't keep the id, we won't delete it.

  return raptor
}

func (r *Raptor) Shape() Shapeish { return r.shape }

func (r *Raptor) GoLeft() { r.velocity.X = -1 }

func (r *Raptor) GoRight() { r.velocity.X = 1 }

func (r *Raptor) GoDown() { r.velocity.Y = 1 }

func (r *Raptor) GoUp() { r.velocity.Y = -1 }

func (r *Raptor) Draw() {
  Draw(RaptorSurface, r.position)
  if ActorTicks-r.lastDamage < 10 {
    Draw(RaptorShieldSurface, r.position)
  }

  r.DrawShieldBar(r.shield, 100.0)
  r.DrawWeapon()
}

func (h *Raptor) DrawShieldBar(current, max float64) {
  var y int16
  var r, g uint32
  var rect *sdl.Rect

  height := int16(Height)
  width := int16(Width)

  rect = &sdl.Rect{
    X: width - 10,
    Y: 0,
    H: uint16(height),
    W: 10,
  }

  Screen.FillRect(rect, 0x000000)

  y = height - int16((float64(height)/max)*current)
  if y%2 != 0 {
    y++
  }

  for ; y < height; y += 2 {
    rect = &sdl.Rect{
      X: width - 7,
      Y: y,
      H: 1,
      W: 5,
    }

    r = 0xff << 16
    g = (0xff - uint32((0xff/float64(height))*float64(y))) << 8

    Screen.FillRect(rect, r+g)
  }
}

func (r *Raptor) DrawWeapon() {
  var dstrect, srcrect *sdl.Rect

  text := RenderText(r.weapon, Red)
  text.GetClipRect(srcrect)

  w, h := uint16(text.W), uint16(text.H)

  dstrect = &sdl.Rect{
    X: (int16(Width) / 2) - (int16(w) / 2),
    Y: int16(Height) - int16(h),
    W: w,
    H: h,
  }

  Screen.Blit(dstrect, text, srcrect)

  text.Free()
}

// Follow the velocity and become slower over time.
func (r *Raptor) Act() {
  vX, vY := r.velocity.X, r.velocity.Y
  if vX > -0.6 && vX < 0.6 {
    r.velocity.X = 0
  }
  if vY > -0.6 && vY < 0.6 {
    r.velocity.Y = 0
  }
  r.velocity.Replace(r.velocity.DivideNum(1.4))

  target := r.position.Plus(r.velocity.Normalize().MultiplyNum(r.speed))

  var top, left, right, bottom float64
  top = 0
  left = 0
  right = float64(Width - 10)
  bottom = float64(Height - int(r.shape.height/2))

  if target.X < left {
    target.X = left
  }
  if target.X > right {
    target.X = right
  }
  if target.Y < top {
    target.Y = top
  }
  if target.Y > bottom {
    target.Y = bottom
  }

  r.position.Replace(target)

  if ActorTicks-r.lastShot > 10 {
    if ActorTicks%50 == 0 {
      r.shield++
    }
  }
}

// receive damage
func (r *Raptor) TakeDamage(amount float64) {
  if amount < r.shield {
    r.shield -= amount
  } else if r.shields > 0 {
    amount -= r.shield
    r.shields--
    r.shield = ShieldEnergy
    r.TakeDamage(amount)
  } else {
    r.Die()
    return
  }

  r.lastDamage = ActorTicks
}

func (r *Raptor) Die() {
  Println("Raptor dies")
  GameOver = true
}

// We finished the level successfully
func (r *Raptor) Wins() {
  Println("You win!")
  GameWin = true
}

// Called after game over.
func (r *Raptor) Reset() {
  r.shield = ShieldEnergy
  r.shields = 1
  GameOver = false
}

func (r *Raptor) FireWeapon() {
  // fire constant weapon
  // r.FireMicroMissile()
  // r.FireTsunami()
  r.FireReaver()

  // fire the selected weapon
  switch r.weapon {
  case "Reaver":
    r.FireReaver()
  case "Thor":
    r.FireThor()
  case "Odin":
    r.FireOdin()
  case "Mauler":
    r.FireMauler()
  case "Tsunami":
    r.FireTsunami()
  case "Deathray":
    r.FireDeathray()
  case "Eclipse":
    r.FireEclipse()
  }

  r.lastShot = ActorTicks
}

// Fire Reavers from each side of the wings
func (r *Raptor) FireReaver() {
  if ActorTicks%ReaverRate == 0 {
    NewReaver(*r.position.Plus(&Vector{X: -10, Y: 0}))
    NewReaver(*r.position.Plus(&Vector{X: 10, Y: 0}))
  }
}

// Fire Micro-Missiles from each side of the wings
func (r *Raptor) FireMicroMissile() {
  if ActorTicks%MicroMissileRate == 0 {
    NewMicroMissile(*r.position.Plus(&Vector{X: -8, Y: 0}))
    NewMicroMissile(*r.position.Plus(&Vector{X: 8, Y: 0}))
  }
}

// Fire Tsunami Pulse Cannon from center
func (r *Raptor) FireTsunami() {
  if ActorTicks%TsunamiRate == 0 {
    NewTsunami(*r.position)
  }
}

// Fire Tsunami Pulse Cannon from center
func (r *Raptor) FireThor() {
  if ActorTicks%ThorRate == 0 {
    NewThor(*r.position.Plus(&Vector{X: -10, Y: 0}))
    NewThor(*r.position.Plus(&Vector{X: 10, Y: 0}))
  }
}

// Fire Blitz, silly weapon o_O
func (r *Raptor) FireBlitz() {
  if ActorTicks%BlitzRate == 0 {
    NewBlitz(*r.position)
  }
}

// Fire Odin laser
func (r *Raptor) FireOdin() {
  if ActorTicks%OdinRate == 0 {
    NewOdin()
  }
}

// Fire Mauler Missile
func (r *Raptor) FireMauler() {
  if ActorTicks%MaulerRate == 0 {
    NewMauler(*r.position.Plus(&Vector{-15, 0}), Vector{-0.5, -0.5})
    NewMauler(*r.position.Plus(&Vector{15, 0}), Vector{0.5, -0.5})
  }
}

// Fire Deathray
func (r *Raptor) FireDeathray() {
  if ActorTicks%DeathrayRate == 0 {
    NewDeathray(r.position)
  }
}

// Fire Eclipse twin laser
func (r *Raptor) FireEclipse() {
  if ActorTicks%EclipseRate == 0 {
    NewEclipse(-15, 10)
    NewEclipse(15, 10)
  }
}
