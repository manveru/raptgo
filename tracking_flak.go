/*
Tracking Flak
-------------
appearance: orange bullet, exact same as Flak
damage: 1 bar (1 bar in training mode)
projectile speed: slow

The only tracking enemy weapon in the game, so that really helps things. Unless
you count the monkey coconuts. Their version of the Auto-Track Mini-Gun except
much worse. In later stages, the tracking turrets will shoot blobs of Tracking
Flak that look like one bullet but are actually several bunched up.

After testing the monkey coconuts they ended up being the exact same as the
Tracking Flak, just brown instead of orange. They seem way more dangerous but
it's probably just because the monkeys throw so many of them.
*/
package main

var (
  TrackingFlakDamage   float64 = 2
  TrackingFlakRate     int64   = 40
  TrackingFlakVelocity         = &Vector{X: 0, Y: -4.0}
  TrackingFlakSurface  *Surface
  TrackingFlakSize     int16 = 5
)

type TrackingFlak struct {
  position              *Vector
  velocity              *Vector
  shape                 *CircleShape
  DrawablesId, ActorsId int64
}

func NewTrackingFlak(pos Vector) *TrackingFlak {
  position := &pos

  flak := &TrackingFlak{
    position: position,
    shape:    NewCircleShape(position, TrackingFlakSize),
  }

  if TrackingFlakSurface == nil {
    TrackingFlakSurface = CreateImage("gfx/flak.png")
  }

  // the tracking flak uses a dumb approximation of the current position of the
  // player since the player velocity changes constantly and to make it
  // actually possible to survive for the poor player.
  flak.velocity = flak.position.Minus(Player.position).Normalize().MultiplyNum(TrackingFlakVelocity.Length())

  Register(flak)

  return flak
}

func (f *TrackingFlak) Draw() { Draw(TrackingFlakSurface, f.position) }

func (f *TrackingFlak) Act() {
  if f.position.Y > float64(Height) {
    Unregister(f)
    return
  }

  if f.shape.Collides(Player.shape) {
    Player.TakeDamage(TrackingFlakDamage)
    Unregister(f)
    return
  }

  f.position.Replace(f.position.Minus(f.velocity))
}
