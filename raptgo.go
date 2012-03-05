package main

import (
  "fmt"
  "github.com/banthar/Go-SDL/sdl"
  "github.com/banthar/Go-SDL/ttf"
  "math/rand"
  "reflect"
  "time"
)

var (
  Width                 = 1024
  Height                = 640
  MouseX, MouseY        int
  Font                  *ttf.Font
  MousePressed          *sdl.MouseButtonEvent
  drawablesId           int64 = 0
  actorsId              int64 = 0
  enemiesId             int64 = 0
  ActorTicks            int64 = 0
  FrameTicks            int64 = 0
  Running                     = true
  Paused                      = false
  GameWin                     = false
  GameOver                    = false
  Player                *Raptor
  Drawables             = DrawableVector{}
  Actors                = ActorVector{}
  Screen, ShieldSurface *Surface

  // some colors
  Black      = sdl.Color{0x00, 0x00, 0x00, 0xff}
  Blue       = sdl.Color{0x00, 0x00, 0xff, 0xff}
  Orange     = sdl.Color{0x66, 0xaa, 0x00, 0xff}
  Red        = sdl.Color{0xff, 0x00, 0x00, 0xff}
  White      = sdl.Color{0xff, 0xff, 0xff, 0xff}
  Yellow     = sdl.Color{0xff, 0xff, 0x00, 0xff}
  Background *Surface
)

func init() {
  sdlSetup()
  gfxSetup()

  rand.Seed(time.Now().Unix())
}

func main() {
  go HandleActors()
  go HandleDrawing()
  go Level2()
  HandleEvents()

  defer Quit()
}

func p(a ...interface{}) {
  fmt.Println(a)
}

func Quit() {
  ttf.Quit()
  sdl.Quit()
  fmt.Println("Have a nice day.")
  fmt.Printf("Your earned: $%d\n", Player.money)
}

func HandleDrawing() {
  for Running {
    for Background == nil {
      sdl.Delay(100)
    }
    DrawBackground()

    for _, drawable := range Drawables {
      drawable.Draw()
    }

    Player.Draw() // separate so it's always on top
    DrawScore()
    DrawShieldIcons()

    if GameOver {
      DrawGameOver()
    } else if Paused {
      DrawPaused()
    } else if GameWin {
      DrawWin()
    }

    Screen.Flip()
    sdl.Delay(24)
    FrameTicks++
  }
}

func HandleActors() {
  for Running {
    if !GameWin && !Paused && !GameOver {
      ActorTicks++

      for _, actor := range Actors {
        actor.Act()
      }
    }

    sdl.Delay(25)
  }
}

func HandleEvents() {
  for Running {
    for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
      switch e := ev.(type) {
      case *sdl.QuitEvent:
        Running = false
      case *sdl.MouseMotionEvent:
        MouseMotion(e)
      case *sdl.MouseButtonEvent:
        switch e.Type {
        case sdl.MOUSEBUTTONDOWN:
          MouseDown(e)
        case sdl.MOUSEBUTTONUP:
          MouseUp(e)
        }
      }
    }

    keys := sdl.GetKeyState()
    for n, i := range keys {
      if i == 1 {
        if Paused {
          switch n {
          case sdl.K_RETURN, sdl.K_SPACE:
            Paused = false
          case sdl.K_q:
            Running = false
          }
        } else if GameOver || GameWin {
          switch n {
          case sdl.K_RETURN, sdl.K_SPACE:
            Player.Reset()
          case sdl.K_q:
            Running = false
          }
        } else {
          switch n {
          case sdl.K_LSHIFT, sdl.K_RSHIFT, sdl.K_LCTRL, sdl.K_RCTRL:
            Player.FireWeapon()
          case sdl.K_0:
            Player.weapon = Player.weapons['0']
          case sdl.K_1:
            Player.weapon = Player.weapons['1']
          case sdl.K_2:
            Player.weapon = Player.weapons['2']
          case sdl.K_3:
            Player.weapon = Player.weapons['3']
          case sdl.K_5:
            Player.weapon = Player.weapons['5']
          case sdl.K_9:
            Player.weapon = Player.weapons['9']
          case sdl.K_MINUS:
            Player.weapon = Player.weapons['-']
          case sdl.K_a, sdl.K_LEFT:
            Player.GoLeft()
          case sdl.K_d, sdl.K_RIGHT:
            Player.GoRight()
          case sdl.K_s, sdl.K_DOWN:
            Player.GoDown()
          case sdl.K_w, sdl.K_UP:
            Player.GoUp()
          case sdl.K_ESCAPE:
            Running = false
          case sdl.K_p:
            Paused = true
          }
        }
      }
    }

    sdl.Delay(25)
  }
}

func MouseDown(button *sdl.MouseButtonEvent) { MousePressed = button }

func MouseUp(button *sdl.MouseButtonEvent) { MousePressed = nil }

func MouseMotion(motion *sdl.MouseMotionEvent) {
  MouseX, MouseY = int(motion.X), int(motion.Y)
}

type DrawableVector map[int64]Drawable
type ActorVector map[int64]Actor

func Register(obj interface{}) {
  value := reflect.Indirect(reflect.ValueOf(obj))
  struc := value

  var f reflect.Value
  if f = struc.FieldByName("DrawablesId"); f.IsValid() {
    f.SetInt(AddDrawable(obj.(Drawable)))
  }
  if f = struc.FieldByName("ActorsId"); f.IsValid() {
    f.Set(reflect.ValueOf(AddActor(obj.(Actor))))
  }
  if f = struc.FieldByName("EnemiesId"); f.IsValid() {
    f.Set(reflect.ValueOf(AddEnemy(obj.(*Enemy))))
  }
}

func Unregister(obj interface{}) {
  value := reflect.Indirect(reflect.ValueOf(obj))
  struc := value

  var f reflect.Value
  if f = struc.FieldByName("DrawablesId"); f.IsValid() {
    DelDrawable(f.Int())
  }
  if f = struc.FieldByName("ActorsId"); f.IsValid() {
    DelActor(f.Int())
  }
  if f = struc.FieldByName("EnemiesId"); f.IsValid() {
    DelEnemy(f.Int())
  }
}

type Drawable interface {
  Draw()
}

func DrawBackground() {
  var dstrect, srcrect *sdl.Rect

  dstrect = &sdl.Rect{
    X: 0,
    Y: 0,
    H: uint16(Height),
    W: uint16(Width),
  }

  src := Background.Surface
  srcrect = &sdl.Rect{
    X: 0,
    Y: int16(src.H) - (int16(Height) + int16(ActorTicks)),
    H: uint16(Height),
    W: uint16(Width),
  }

  Screen.Blit(dstrect, src, srcrect)
}

func DrawScore() {
  var dstrect, srcrect *sdl.Rect

  text := RenderText(fmt.Sprintf("$%d", Player.money), Orange)
  text.GetClipRect(srcrect)

  w, h := uint16(text.W), uint16(text.H)

  dstrect = &sdl.Rect{
    X: 0,
    Y: int16(Height) - int16(h),
    W: w,
    H: h,
  }

  Screen.Blit(dstrect, text, srcrect)

  text.Free()
}

func DrawShieldIcons() {
  var dstrect, srcrect *sdl.Rect

  src := ShieldSurface.Surface
  srcrect = &sdl.Rect{
    X: 0,
    Y: 0,
    W: 18,
    H: 20,
  }

  dstrect = &sdl.Rect{
    X: int16((Width - 20) - int(srcrect.W)),
    Y: int16(Height - int(srcrect.H)),
    W: srcrect.W,
    H: srcrect.H,
  }

  for n := 0; n < Player.shields; n++ {
    Screen.Blit(dstrect, src, srcrect)
    dstrect.X -= int16(srcrect.W)
  }
}

// draw lines of text aligned around center of screen
func DrawTextLines(strings []string) {
  var dstrect, srcrect *sdl.Rect
  var line *sdl.Surface
  var w, h, prevH uint16

  centerX := int16(Width / 2)
  centerY := int16(Height / 2)

  for _, str := range strings {
    line = RenderText(str, Red)
    line.GetClipRect(srcrect)

    w, h = uint16(line.W), uint16(line.H)

    dstrect = &sdl.Rect{
      X: centerX - int16(w/2),
      Y: (centerY + int16(prevH)) - int16(h/2),
      W: w,
      H: h,
    }

    prevH += h

    Screen.Blit(dstrect, line, srcrect)
    line.Free()
  }
}

func DrawGameOver() {
  DrawTextLines([]string{
    "Game Over",
    "Press Return or Spacebar to try again",
    "Press q to exit",
  })
}

func DrawPaused() {
  DrawTextLines([]string{
    "Game paused",
    "Press Return or Spacebar to continue",
    "Press q to exit",
  })
}

func DrawWin() {
  DrawTextLines([]string{
    "You Win!",
    "Press Return or Spacebar to play again",
    "Press q to exit",
  })
}

func AddDrawable(drawable Drawable) (id int64) {
  id = drawablesId
  drawablesId++
  Drawables[id] = drawable
  return id
}

func DelDrawable(id int64) { delete(Drawables, id) }

type Actor interface {
  Act()
}

func AddActor(actor Actor) (id int64) {
  id = actorsId
  actorsId++
  Actors[id] = actor
  return id
}

func DelActor(id int64) { delete(Actors, id) }

func sdlSetup() {
  if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
    panic(sdl.GetError())
  }

  Screen = &Surface{sdl.SetVideoMode(Width, Height, 32, 0)}
  if Screen == nil {
    panic(sdl.GetError())
  }

  if sdl.EnableKeyRepeat(25, 25) != 0 {
    panic(sdl.GetError())
  }

  if ttf.Init() != 0 {
    panic(sdl.GetError())
  }
}

func gfxSetup() {
  Font = NewFont("DroidSansMono.ttf", 20)

  Player = NewRaptor(&Vector{320, 400})

  MakeMoneybagSurface()
  MakeShieldpackSurface()
  ShieldSurface = CreateImage("gfx/shield.png")
}

func NewFont(name string, size int) (font *ttf.Font) {
  font = ttf.OpenFont("/usr/share/fonts/TTF/"+name, size)
  if font == nil {
    panic(sdl.GetError())
  }
  return font
}

func RenderText(text string, color sdl.Color) *sdl.Surface {
  return ttf.RenderText_Blended(Font, text, color)
}

func sleep(delay uint32) { sdl.Delay(delay) }
