package main

import (
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/core/world"
	"math/rand"
	"time"
)


func initWorld() contracts.IWorld {
	agent := GetAgent()

	terrain := GetTerrain()
	terrain.GetCell(1, 0).SetCellType(contracts.LockedCell)
	terrain.GetCell(1, 0).SetAgent(agent)

	return &world.World{
		terrain,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	workspace := Workspace{
		Width:         1920,
		Height:        1080,
	}

	workspace.Init(initWorld())
	workspace.Start()
}
