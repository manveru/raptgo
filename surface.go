package main

import (
  "github.com/banthar/Go-SDL/sdl"
  "os"
  "strings"
)

type Surface struct {
  *sdl.Surface
}

func Draw(src *Surface, pos *Vector) {
  if src == nil {
    panic("Invalid src *Surface == nil")
  }
  var srcrect, dstrect *sdl.Rect
  src.GetClipRect(srcrect)

  // center surface
  w, h := src.W, src.H
  dstrect = &sdl.Rect{
    X: int16(int32(pos.X) - (w / 2)),
    Y: int16(int32(pos.Y) - (h / 2)),
    W: uint16(w),
    H: uint16(h),
  }

  Screen.Blit(dstrect, src.Surface, srcrect)
}

func CreateRGBSurface(width, height int) *Surface {
  var rmask, gmask, bmask, amask uint32

  byteorder := 0
  bigendian := 0

  if byteorder == bigendian {
    rmask = 0xff000000
    gmask = 0x00ff0000
    bmask = 0x0000ff00
    amask = 0x000000ff
  } else {
    rmask = 0x000000ff
    gmask = 0x0000ff00
    bmask = 0x00ff0000
    amask = 0xff000000
  }

  inner := sdl.CreateRGBSurface(
    sdl.SRCALPHA, // flags
    width, height,
    32, // depth
    rmask, gmask, bmask, amask,
  )

  if inner == nil {
    panic(sdl.GetError())
  }

  return &Surface{inner}
}

func CreateImage(path string) *Surface {
  gopaths := strings.Split(os.Getenv("GOPATH"), ":")

  var img *sdl.Surface
  for _, gopath := range gopaths {
    img = sdl.Load(gopath + "/src/github.com/manveru/raptgo/" + path)
    if img != nil {
      return &Surface{img}
    }
  }

  panic(sdl.GetError())
}

func (self *Surface) CircleOutlinePoints(cx, cy, x, y int16, color uint32) {
  if x == 0 {
    self.DrawPixel(cx, cy+y, color)
    self.DrawPixel(cx, cy-y, color)
    self.DrawPixel(cx+y, cy, color)
    self.DrawPixel(cx-y, cy, color)
  } else if x == y {
    self.DrawPixel(cx+x, cy+y, color)
    self.DrawPixel(cx-x, cy+y, color)
    self.DrawPixel(cx+x, cy-y, color)
    self.DrawPixel(cx-x, cy-y, color)
  } else if x < y {
    self.DrawPixel(cx+x, cy+y, color)
    self.DrawPixel(cx-x, cy+y, color)
    self.DrawPixel(cx+x, cy-y, color)
    self.DrawPixel(cx-x, cy-y, color)
    self.DrawPixel(cx+y, cy+x, color)
    self.DrawPixel(cx-y, cy+x, color)
    self.DrawPixel(cx+y, cy-x, color)
    self.DrawPixel(cx-y, cy-x, color)
  }
}

func (self *Surface) CircleFillPoints(cx, cy, x, y int16, color uint32) {
  if x == 0 {
    self.DrawLine(cx, cy+y, cx, cy-y, color)
    self.DrawLine(cx+y, cy, cx-y, cy, color)
  } else if x == y {
    self.DrawLine(cx+x, cy+y, cx-x, cy+y, color)
    self.DrawLine(cx+x, cy-y, cx-x, cy-y, color)
  } else if x < y {
    self.DrawLine(cx+x, cy+y, cx-x, cy+y, color)
    self.DrawLine(cx+x, cy-y, cx-x, cy-y, color)
    self.DrawLine(cx+y, cy+x, cx-y, cy+x, color)
    self.DrawLine(cx+y, cy-x, cx-y, cy-x, color)
  }
}

func (self *Surface) DrawCircleOutline(x, y, radius int16, color uint32) {
  self.drawCircleWith(x, y, radius, color, (*Surface).CircleOutlinePoints)
}

func (self *Surface) DrawCircleFill(x, y, radius int16, color uint32) {
  self.drawCircleWith(x, y, radius, color, (*Surface).CircleFillPoints)
}

type drawCircleFun func(*Surface, int16, int16, int16, int16, uint32)

func (self *Surface) drawCircleWith(cx, cy, radius int16, color uint32, fun drawCircleFun) {
  var x, y, p int16
  x = 0
  y = radius
  p = (5 - radius*4) / 4

  fun(self, cx, cy, x, y, color)

  for x < y {
    x++
    if p < 0 {
      p += 2*x + 1
    } else {
      y--
      p += 2*(x-y) + 1
    }
    fun(self, cx, cy, x, y, color)
  }
}

func (self *Surface) DrawPixel(x, y int16, color uint32) {
  rect := &sdl.Rect{X: x, Y: y, W: 1, H: 1}
  self.FillRect(rect, color)
}

func (self *Surface) DrawLine(x1, y1, x2, y2 int16, color uint32) {
  var i, deltax, deltay, numpixels,
    d, dinc1, dinc2,
    x, xinc1, xinc2,
    y, yinc1, yinc2 int16
  // calculate deltax and deltay for initialization
  deltax = abs(x2 - x1)
  deltay = abs(y2 - y1)
  // initialize all vars based on which is the independent variable
  if deltax >= deltay {
    // x is independent variable
    numpixels = deltax + 1
    d = (2 * deltay) - deltax
    dinc1 = deltay * 2
    dinc2 = (deltay - deltax) * 2
    xinc1 = 1
    xinc2 = 1
    yinc1 = 0
    yinc2 = 1
  } else {
    // y is independent variable
    numpixels = deltay + 1
    d = (2 * deltax) - deltay
    dinc1 = deltax * 2
    dinc2 = (deltax - deltay) * 2
    xinc1 = 0
    xinc2 = 1
    yinc1 = 1
    yinc2 = 1
  }

  // make sure x and y move in the right directions
  if x1 > x2 {
    xinc1 = -xinc1
    xinc2 = -xinc2
  }
  if y1 > y2 {
    yinc1 = -yinc1
    yinc2 = -yinc2
  }

  // start drawing at
  x = x1
  y = y1

  // draw the pixels
  for i = 1; i <= numpixels; i++ {
    self.DrawPixel(x, y, color)
    if d < 0 {
      d += dinc1
      x += xinc1
      y += yinc1
    } else {
      d += dinc2
      x += xinc2
      y += yinc2
    }
  }
}

func abs(x int16) int16 {
  switch {
  case x < 0:
    return -x
  case x == 0:
    return 0
  }
  return x
}
