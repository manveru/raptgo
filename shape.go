package main

import "sdl"

type Shapeish interface {
  Collides(Shapeish) bool
  Draw()
}

type CircleShape struct {
  radius   int16
  position *Vector
}

func NewCircleShape(position *Vector, radius int16) *CircleShape {
  return &CircleShape{
    radius:   radius,
    position: position,
  }
}

func (circle *CircleShape) IsCircle() bool { return true }

func (circle *CircleShape) Collides(shape Shapeish) bool {
  switch shape.(type) {
  case *CircleShape:
    return Circle2Circle(circle, shape.(*CircleShape))
  case *RectangleShape:
    return Circle2Rectangle(circle, shape.(*RectangleShape))
  }

  return false
}

func (shape *CircleShape) ToSurface(color uint32) *Surface {
  radius := shape.radius
  length := int((shape.radius * 2) + 2)
  surface := CreateRGBSurface(length, length)
  half := int16(length / 2)
  surface.DrawCircleFill(half, half, radius, color)

  return surface
}

func (shape *CircleShape) Draw() {}

type RectangleShape struct {
  width, height int16
  position      *Vector
}

func NewRectangleShape(position *Vector, width, height int16) *RectangleShape {
  return &RectangleShape{
    position: position,
    width:    width,
    height:   height,
  }
}

func (rect *RectangleShape) Collides(shape Shapeish) bool {
  switch shape.(type) {
  case *CircleShape:
    return Circle2Rectangle(shape.(*CircleShape), rect)
  case *RectangleShape:
    return Rectangle2Rectangle(shape.(*RectangleShape), rect)
  }

  return false
}

func (shape *RectangleShape) Draw() {
  pos := shape.position
  rect := &sdl.Rect{
    X: int16(pos.X) - shape.width/2,
    Y: int16(pos.Y) - shape.height/2,
    W: uint16(shape.width),
    H: uint16(shape.height),
  }
  Screen.FillRect(rect, 0xff0000)
}

func (rect *RectangleShape) Top() float64 { return rect.position.Y - float64(rect.height/2) }
func (rect *RectangleShape) Bottom() float64 {
  return rect.position.Y + float64(rect.height/2)
}
func (rect *RectangleShape) Left() float64  { return rect.position.X - float64(rect.width/2) }
func (rect *RectangleShape) Right() float64 { return rect.position.X + float64(rect.width/2) }

func Circle2Circle(circle1, circle2 *CircleShape) bool {
  diff := circle1.position.Minus(circle2.position)
  return diff.Length() <= float64(circle1.radius+circle2.radius)
}

func Fclamp(value, min, max float64) float64 {
  if value < min {
    return min
  }
  if value > max {
    return max
  }
  return value
}

// Thank you, StackOverflow.
func Circle2Rectangle(circle *CircleShape, rect *RectangleShape) bool {
  cpos := circle.position
  rpos := rect.position
  rectLeft := rpos.X - float64(rect.width/2)
  rectRight := rpos.X + float64(rect.width/2)
  rectTop := rpos.Y - float64(rect.height/2)
  rectBottom := rpos.Y + float64(rect.height/2)
  radius := float64(circle.radius)

  // Find the closest point to the circle within the rectangle
  closestX := Fclamp(cpos.X, rectLeft, rectRight)
  closestY := Fclamp(cpos.Y, rectTop, rectBottom)

  // Calculate the distance between the circle's center and this closest point
  distanceX := cpos.X - closestX
  distanceY := cpos.Y - closestY

  // If the distance is less than the circle's radius, an intersection occurs
  distanceSquared := (distanceX * distanceX) + (distanceY * distanceY)
  return distanceSquared < (radius * radius)
}

func Rectangle2Rectangle(rect1, rect2 *RectangleShape) bool {
  return rect1.Left() < rect2.Right() &&
    rect1.Right() > rect2.Left() &&
    rect1.Top() < rect2.Bottom() &&
    rect1.Bottom() > rect2.Top()
}
