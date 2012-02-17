// TH19 "Thor" Auto-Track Minigun
package main

import (
	"math"
)

var (
	ThorDamage   float64 = 1
	ThorRate     int64  = 4
	ThorVelocity         = &Vector{X: 0, Y: 6.0}
	ThorCost             = 12000
	ThorSurface  [7]*Surface
	ThorSize     int16 = 3
)

type Thor struct {
	position              *Vector
	velocity              *Vector
	shape                 *CircleShape
	DrawablesId, ActorsId int64
}

func NewThor(pos Vector) *Thor {
	position := &pos

	thor := &Thor{
		position: position,
		shape:    NewCircleShape(position, ThorSize),
	}

	// shares animation of reaver
	if ThorSurface[0] == nil {
		ThorSurface = [7]*Surface{
			CreateImage("gfx/reaver_0.png"),
			CreateImage("gfx/reaver_1.png"),
			CreateImage("gfx/reaver_2.png"),
			CreateImage("gfx/reaver_3.png"),
			CreateImage("gfx/reaver_1.png"),
			CreateImage("gfx/reaver_2.png"),
			CreateImage("gfx/reaver_3.png"),
		}
	}

	if velocity, ok := thor.CalculateVelcoity(); ok {
		thor.velocity = &velocity
		Register(thor)
	}

	return thor
}

func (t *Thor) Draw() {
	surface := ThorSurface[FrameTicks%7]
	Draw(surface, t.position)
}

func (t *Thor) Act() {
	t.position.Replace(t.position.Minus(t.velocity))

	if t.position.Y < 0 {
		Unregister(t)
		return
	}

	for _, enemy := range Enemies {
		if t.shape.Collides(enemy.shape) {
			enemy.TakeDamage(ThorDamage)
			Unregister(t)
			return
		}
	}
}

// calculate velocity once when firing.
// takes enemy velocity into account, but not too much so fastest enemies can
// escape.
// From http://www.allegro.cc/forums/print-thread/591292
// Assuming pos_u is the position of the target, pos_p is the position of the
// missile, vel_u is the velocity vector of the target, and speed_p is the
// max speed of the missile.
func (t *Thor) CalculateVelcoity() (vel Vector, ok bool) {
	// first go through enemies, gather the distances.
	minLength := float64(Width * Height) // more than enough
	var lock *Enemy

	for _, enemy := range Enemies {
		length := t.position.Minus(enemy.position).Length()
		if length < minLength {
			minLength = length
			lock = enemy
		}
	}

	if lock == nil {
		return *ThorVelocity, false
	}

	pos_u := lock.position
	pos_p := t.position
	vel_u := lock.velocity
	speed_p := ThorVelocity.Length()

	trans_p := pos_p.Minus(pos_u)
	d_squared := trans_p.LengthSquared()
	target_dir := vel_u.Normalize()
	y := trans_p.Dot(target_dir)
	speed_u_sq := vel_u.LengthSquared()
	speed_p_sq := speed_p * speed_p

	var i, i1, i2 float64

	if math.Abs(speed_u_sq/speed_p_sq-1.0) <= 0.001 {
		i1 = -1.0
		i2 = d_squared / (2 * y)
	} else {
		speed_coeff := speed_p_sq/speed_u_sq - 1.0
		delta := y*y + speed_coeff*d_squared
		if delta < 0.0 {
			// no way we could intercept, the bullet is too slow.
			return *ThorVelocity, false
		}
		delta = math.Sqrt(delta)
		i1 = (-y + delta) / speed_coeff
		i2 = (-y - delta) / speed_coeff
	}

	if i1 > 0.0 && i2 > 0.0 {
		if i1 > i2 {
			i = i2
		} else {
			i = i1
		}
	} else if i1 > 0.0 {
		i = i1
	} else if i2 > 0.0 {
		i = i2
	} else {
		// missile is too slow to intercept
		return *ThorVelocity, false
	}

	intercept_pos := pos_u.Plus(target_dir.MultiplyNum(i))
	vel = *pos_p.Minus(intercept_pos).Normalize().MultiplyNum(speed_p)

	return vel, true
}
