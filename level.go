package main

import "math/rand"

// import . "fmt"

func Level1() {
  Background = CreateImage("gfx/level1.png")

  for i := float64(102.4); i < 922; i += 102.4 {
    NewHornet(Vector{i, 0})
    sleep(1000)
  }

  for i := float64(102.4); i < 922; i += 102.4 {
    NewStinger(Vector{i, 0})
    sleep(1000)
  }

  for i := float64(102.4); i < 922; i += 102.4 {
    NewBee(Vector{i, 0})
    sleep(2000)
  }

  for i := float64(102.4); i < 922; i += 102.4 {
    NewStinger(Vector{i, 0})
    sleep(500)
  }

  for i := float64(102.4); i < 922; i += 102.4 {
    NewBee(Vector{i, 0})
    sleep(1000)
  }

  // and ten more hard waves
  for i := 0; i < 10; i++ {
    for i := float64(102.4); i < 922; i += 102.4 {
      NewStinger(Vector{i, 0})
    }

    sleep(10000)

    for i := float64(102.4); i < 922; i += 102.4 {
      NewBee(Vector{i, 0})
    }

    sleep(15000)
  }

  for len(Enemies) > 0 {
    sleep(1000)
  }

  Player.Wins()
}

// spawn 100 random enemies at random locations with random pauses.
// rightfully you could call this the random level
func Level2() {
  Background = CreateImage("gfx/level1.png")

  for t := 0; t < 100; t++ {
    r := rand.Float64()
    x := float64(rand.Intn(Width))
    v := Vector{x, 0}

    if r < 0.25 {
      NewHornet(v)
    } else if r < 0.5 {
      NewStinger(v)
    } else if r < 0.75 {
      NewBee(v)
    } else {
      NewMiner(v)
    }

    sleep(uint32(rand.Intn(3000)))
  }

  for len(Enemies) > 0 {
    sleep(1000)
  }

  Player.Wins()
}
