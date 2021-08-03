package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	workspace := Workspace{
		Width:  1920,
		Height: 1080,
	}

	workspace.Init()
	workspace.Start()
}
