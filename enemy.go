package main

import "rand"

type EnemyVector map[uint64]*Enemy
var Enemies = EnemyVector{}

func AddEnemy(e *Enemy) (id uint64) {
  id = enemiesId
  enemiesId++
  Enemies[id] = e
  return id
}

func DelEnemy(id uint64) { Enemies[id] = nil, false }

type Enemy struct {
  dead                             bool
  EnemiesId, ActorsId, DrawablesId uint64
  health                           float64
  moneyReward                      int
  position, velocity               *Vector
  shape                            Shapeish
  surface                          *Surface
  weapon                           func(*Enemy)
}

func NewStinger(position Vector) *Enemy { return CreateEnemy("stinger", position) }

func NewBee(position Vector) *Enemy { return CreateEnemy("bee", position) }

func NewHornet(position Vector) *Enemy { return CreateEnemy("hornet", position) }

func NewMiner(position Vector) *Enemy { return CreateEnemy("miner", position) }

func (e *Enemy) Shape() Shapeish { return e.shape }
func (e *Enemy) Draw()           { Draw(e.surface, e.position) }

func (e *Enemy) TakeEnergy(amount float64) bool {
  if ActorTicks%20 == 0 {
    return true
  }
  return false
}

func (e *Enemy) TakeDamage(amount float64) {
  e.health -= amount
  if !e.dead && e.health <= 0 {
    e.dead = true
    Unregister(e)
    Player.money += e.moneyReward

    // every 10 enemies or so, drop goodies
    if rand.Float() < 0.1 {
      e.MaybeDropsGoodie()
    }
  }
}

func (e *Enemy) MaybeDropsGoodie() {
  r := rand.Float()
  if r < 0.5 {
    NewMoneybag(*e.position, *e.velocity)
  } else {
    NewShieldpack(*e.position, *e.velocity)
  }
}

// destroy the player with powerful bullets
// if we get to the lower edge, self-destruct.
func (e *Enemy) Act() {
  if e.position.Y >= float64(Height) {
    e.dead = true
    Unregister(e)
    return
  }

  if e.shape.Collides(Player.shape) {
    Player.TakeDamage(e.health)
    e.dead = true
    Unregister(e)
    return
  }

  e.position.Replace(e.position.Plus(e.velocity))

  e.FireWeapon()
}

func (e *Enemy) UseWeapon(weapon func(*Enemy)) {
  e.weapon = weapon
}

func (e *Enemy) FireWeapon() {
  if e.weapon != nil {
    e.weapon(e)
  }
}

// Fire Flak from each side of the wings
func (e *Enemy) FireFlak() {
  if ActorTicks%FlakRate == 0 {
    NewFlak(*e.position.Plus(&Vector{X: -10, Y: 0}))
    NewFlak(*e.position.Plus(&Vector{X: 10, Y: 0}))
  }
}

// Fire Tracking Flak from each side of the wings
func (e *Enemy) FireTrackingFlak() {
  if ActorTicks%TrackingFlakRate == 0 {
    NewTrackingFlak(*e.position.Plus(&Vector{X: -10, Y: 0}))
    NewTrackingFlak(*e.position.Plus(&Vector{X: 10, Y: 0}))
  }
}

// Fire Missile from each side of the wings
func (e *Enemy) FireMissile() {
  if ActorTicks%MissileRate == 0 {
    NewMissile(*e.position.Plus(&Vector{X: -10, Y: 0}))
    NewMissile(*e.position.Plus(&Vector{X: 10, Y: 0}))
  }
}

// Place a mine
func (e *Enemy) FireMine() {
  if ActorTicks%MineRate == 0 {
    NewMine(*e.position)
  }
}
