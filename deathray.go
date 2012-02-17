// MSIL-ATLAS Deathray
package main

var (
  DeathrayDamage   float64 = 24
  DeathrayRate     int64  = 30
  DeathrayCost     = 950000
  DeathrayLifeTime = 10
  DeathrayWidth    int16  = 5
  DeathrayColor    uint32 = 0xff6600
)

type Deathray struct {
  from, to              *Vector
  lifeTime              int
  shape                 *RectangleShape
  DrawablesId, ActorsId int64
}

func NewDeathray(position *Vector) *Deathray {
  ray := &Deathray{
    from:     position,
    lifeTime: DeathrayLifeTime,
    shape: NewRectangleShape(
      &Vector{position.X, position.Y - position.Y/2},
      DeathrayWidth,
      int16(position.Y),
    ),
  }

  Register(ray)

  return ray
}

func (d *Deathray) Draw() {
  rect := d.shape
  t := int16(rect.Top())
  b := int16(rect.Bottom())
  l := int16(rect.Left())
  r := int16(rect.Right())

  for x := l; x <= r; x++ {
    Screen.DrawLine(x, t, x, b, DeathrayColor)
  }
}

func (d *Deathray) Act() {
  if d.lifeTime < 1 {
    Unregister(d)
  } else {
    d.shape.position = &Vector{d.from.X, d.from.Y - d.from.Y/2}
    d.shape.height = int16(d.from.Y)
    d.CheckCollision()
    d.lifeTime--
  }
}

func (d *Deathray) CheckCollision() {
  for _, enemy := range Enemies {
    if enemy.shape.Collides(d.shape) {
      enemy.TakeDamage(DeathrayDamage)
    }
  }
}
