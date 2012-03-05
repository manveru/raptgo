// CAL-10 "Eclipse" Twin Lasers
package main

var (
  EclipseDamage   float64 = 40
  EclipseRate     int64   = 30
  EclipseCost             = 950000
  EclipseLifeTime         = 10
  EclipseWidth    int16   = 5
  EclipseColor    uint32  = 0xff6600
)

type Eclipse struct {
  lifeTime              int
  offsetX, offsetY      float64
  shape                 *RectangleShape
  DrawablesId, ActorsId int64
}

func NewEclipse(offsetX, offsetY float64) *Eclipse {
  pp := Player.position

  eclipse := &Eclipse{
    lifeTime: EclipseLifeTime,
    offsetX:  offsetX,
    offsetY:  offsetY,
    shape: NewRectangleShape(
      &Vector{pp.X + offsetX, (pp.Y - pp.Y/2) + offsetY},
      EclipseWidth,
      int16(pp.Y),
    ),
  }

  Register(eclipse)

  return eclipse
}

func (e *Eclipse) Draw() {
  rect := e.shape
  t := int16(rect.Top())
  b := int16(rect.Bottom())
  l := int16(rect.Left())
  r := int16(rect.Right())

  for x := l; x <= r; x++ {
    Screen.DrawLine(x, t, x, b, EclipseColor)
  }
}

func (e *Eclipse) Act() {
  if e.lifeTime < 1 {
    Unregister(e)
  } else {
    pp := Player.position
    e.shape.position = &Vector{pp.X + e.offsetX, (pp.Y - pp.Y/2) + e.offsetY}
    e.shape.height = int16(pp.Y)
    e.CheckCollision()
    e.lifeTime--
  }
}

func (e *Eclipse) CheckCollision() {
  for _, enemy := range Enemies {
    if enemy.shape.Collides(e.shape) {
      enemy.TakeDamage(EclipseDamage)
    }
  }
}
