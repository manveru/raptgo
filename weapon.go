package main

type Explosion struct {
  position                    *Vector
  ActorsId, DrawablesId       uint64
  shape                       *CircleShape
  lifeTime                    int
  radius                      int16
  lethality                   float64
  intensity                   uint32
  damagesPlayer, damagesEnemy bool
}

func NewExplosion(pos *Vector, lethality float64, damagesEnemy, damagesPlayer bool) *Explosion {
  explosion := &Explosion{
    position:      pos,
    lifeTime:      10,
    intensity:     0x000000ff,
    radius:        1.0,
    lethality:     lethality,
    damagesEnemy:  damagesEnemy,
    damagesPlayer: damagesPlayer,
  }

  explosion.shape = NewCircleShape(explosion.position, int16(explosion.radius))

  Register(explosion)

  return explosion
}

// DIY to make pretty explosion animation
func (e *Explosion) Draw() {
  x, y, r := int16(e.position.X), int16(e.position.Y), e.radius
  var color uint32 = 0xff0000ff
  // alpha according to intensity
  color |= e.intensity
  Screen.DrawCircleOutline(x, y, r, color)
}

func (e *Explosion) Act() {
  if e.lifeTime <= 0 {
    Unregister(e)
  } else {
    e.radius++
    e.intensity--
    e.shape.radius = int16(e.radius)

    if e.damagesPlayer {
      e.DamagePlayer()
    }
    if e.damagesEnemy {
      e.DamageEnemy()
    }
  }

  e.lifeTime--
}

func (e *Explosion) DamageEnemy() {
  for _, enemy := range Enemies {
    if e.shape.Collides(enemy.shape) {
      enemy.TakeDamage(e.lethality)
    }
  }
}

func (e *Explosion) DamagePlayer() {
  if e.shape.Collides(Player.shape) {
    Player.TakeDamage(e.lethality)
  }
}
