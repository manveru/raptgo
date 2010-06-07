package main

import "math"

type Vector struct {
  X, Y float64
}

func (self *Vector) Replace(other *Vector) *Vector {
  self.X, self.Y = other.X, other.Y
  return self
}

// Normalize the vector to a length of 1, takes care of division by zero.
func (self *Vector) Normalize() *Vector {
  length := self.Length()
  x, y := self.X, self.Y
  v := &Vector{}
  if x == 0.0 {
    v.X = x
  } else {
    v.X = x / length
  }
  if y == 0.0 {
    v.Y = y
  } else {
    v.Y = y / length
  }
  return v
}

func (self *Vector) Length() float64 { return math.Sqrt(self.Dot(self)) }

func (self *Vector) LengthSquared() float64 { return self.Dot(self) }

func (self *Vector) Dot(other *Vector) float64 {
  return self.X*other.X + self.Y*other.Y
}

func (self *Vector) SetLength(length float64) {
  angle := self.Angle()
  self.X = math.Cos(angle) * length
  self.Y = math.Sin(angle) * length
}

func (self *Vector) Angle() (angle float64) {
  angle = math.Atan2(self.Y, self.X)
  if angle < 0 {
    angle += 2 * math.Pi
  }
  return
}

func (self *Vector) SetAngle(angle float64) {
  length := self.Length()
  self.X = math.Cos(angle) * length
  self.Y = math.Sin(angle) * length
}

func (self *Vector) Minus(other *Vector) *Vector {
  return &Vector{X: (self.X - other.X), Y: (self.Y - other.Y)}
}

func (self *Vector) MinusNum(other float64) *Vector {
  return &Vector{X: (self.X - other), Y: (self.Y - other)}
}

func (self *Vector) Plus(other *Vector) *Vector {
  return &Vector{X: (self.X + other.X), Y: (self.Y + other.Y)}
}

func (self *Vector) PlusNum(other float64) *Vector {
  return &Vector{X: (self.X + other), Y: (self.Y + other)}
}

func (self *Vector) Multiply(other *Vector) *Vector {
  return &Vector{X: (self.X * other.X), Y: (self.Y * other.Y)}
}

func (self *Vector) MultiplyNum(other float64) *Vector {
  return &Vector{X: (self.X * other), Y: (self.Y * other)}
}

func (self *Vector) Divide(other *Vector) *Vector {
  return &Vector{X: (self.X / other.X), Y: (self.Y / other.Y)}
}

func (self *Vector) DivideNum(other float64) *Vector {
  return &Vector{X: (self.X / other), Y: (self.Y / other)}
}

func (self *Vector) Invert() *Vector { return &Vector{X: -self.X, Y: -self.Y} }
