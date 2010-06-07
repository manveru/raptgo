package main

type EnemySpec struct {
  health      float64
  moneyReward int
  velocity    Vector
  weapon      func(*Enemy)
  sprite      string
  surface     *Surface
}

var EnemySpecs map[string]EnemySpec

func init() {
  EnemySpecs = map[string]EnemySpec{
    "stinger": EnemySpec{
      health:      20,
      moneyReward: 20,
      velocity:    Vector{0, 1},
      weapon:      (*Enemy).FireFlak,
      sprite:      "gfx/stinger.png",
    },
    "bee": EnemySpec{
      health:      20,
      moneyReward: 40,
      velocity:    Vector{0, 0.5},
      weapon:      (*Enemy).FireTrackingFlak,
      sprite:      "gfx/bee.png",
    },
    "hornet": EnemySpec{
      health:      20,
      moneyReward: 40,
      velocity:    Vector{0, 0.5},
      weapon:      (*Enemy).FireMissile,
      sprite:      "gfx/rocketeer.png",
    },
    "miner": EnemySpec{
      health:      10,
      moneyReward: 50,
      velocity:    Vector{0, 2},
      weapon:      (*Enemy).FireMine,
      sprite:      "gfx/miner.png",
    },
  }
}

func CreateEnemy(name string, position Vector) *Enemy {
  spec := EnemySpecs[name]

  if spec.surface == nil {
    spec.surface = CreateImage(spec.sprite)
  }

  enemy := &Enemy{
    health:      spec.health,
    moneyReward: spec.moneyReward,
    position:    &position,
    velocity:    &spec.velocity,
    surface:     spec.surface,
    weapon:      spec.weapon,
  }

  enemy.shape = NewRectangleShape(
    enemy.position,
    int16(enemy.surface.W),
    int16(enemy.surface.H),
  )

  Register(enemy)

  return enemy
}
