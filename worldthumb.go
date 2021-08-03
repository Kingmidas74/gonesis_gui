package main

import (
	preparedCommands "github.com/Kingmidas74/gonesis_engine/commands"
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/core"
	"github.com/Kingmidas74/gonesis_engine/core/agents"
	"github.com/Kingmidas74/gonesis_engine/core/commands"
	"github.com/Kingmidas74/gonesis_engine/core/primitives"
	"github.com/Kingmidas74/gonesis_engine/core/reproductions"
	"github.com/Kingmidas74/gonesis_engine/core/terrains"
)

func GetCommands() map[int]contracts.ICommand {
	commandsMap := make(map[int]contracts.ICommand)
	commandsMap[0] = &preparedCommands.MoveCommand{
		commands.Command{
			IsInterrupt: true,
		},
	}
	commandsMap[1] = &preparedCommands.EatCommand{
		commands.Command{
			IsInterrupt: false,
		},
	}
	return commandsMap
}

func GetTerrain(currentAgents []contracts.IAgent, terrainType int) contracts.ITerrain {

	var terrain contracts.ITerrain

	baseTerrain := terrains.Terrain{
		Cells:  make([]contracts.ICell, 50),
		Width:  10,
		Height: 5,
	}

	if terrainType == 0 {
		terrain = &terrains.MooreTerrain{
			baseTerrain,
		}
	} else if terrainType == 1 {
		terrain = &terrains.NeumannTerrain{
			baseTerrain,
		}
	}

	for i := 0; i < len(terrain.GetCells()); i++ {
		terrain.GetCells()[i] = &terrains.Cell{
			Coords:   primitives.Coords{},
			CellType: contracts.EmptyCell,
			Cost:     0,
			Agent:    nil,
		}
	}

	for y := 0; y < terrain.GetHeight(); y++ {
		for x := 0; x < terrain.GetWidth(); x++ {
			terrain.GetCell(x, y).SetX(x)
			terrain.GetCell(x, y).SetY(y)
		}
	}

	terrain.GetCells()[21].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[21].SetCost(4)

	terrain.GetCells()[23].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[23].SetCost(4)

	terrain.GetCells()[43].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[43].SetCost(4)

	terrain.GetCells()[45].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[45].SetCost(4)

	terrain.GetCells()[15].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[15].SetCost(4)

	terrain.GetCells()[17].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[17].SetCost(4)

	terrain.GetCells()[37].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[37].SetCost(4)

	terrain.GetCells()[39].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[39].SetCost(4)

	terrain.GetCells()[9].SetCellType(contracts.OrganicCell)
	terrain.GetCells()[9].SetCost(4)

	placedChildrenCount := 0

	for i := 0; i < len(terrain.GetCells()) && placedChildrenCount < len(currentAgents); i++ {
		if terrain.GetCells()[i].GetCellType() == contracts.EmptyCell {
			currentAgents[placedChildrenCount].SetX(terrain.GetCells()[i].GetX())
			currentAgents[placedChildrenCount].SetY(terrain.GetCells()[i].GetY())
			terrain.GetCells()[i].SetCellType(contracts.LockedCell)
			terrain.GetCells()[i].SetAgent(currentAgents[placedChildrenCount])
			placedChildrenCount++
		}
	}
	return terrain
}

func GetAgents(agentsCount int) []contracts.IAgent {
	result := make([]contracts.IAgent, 0)
	for i := 0; i < agentsCount; i++ {
		agent := &agents.Agent{
			IBrain: &core.Brain{
				CommandList: commands.CommandList{
					Commands: GetCommands(),
				},
				Commands: []int{
					0, 4, //down
					1, 4, //eat down
					14,
					11,
					2,
					1,
					0, 2, //right
					1, 2, //eat right
					14,
					11,
					2,
					1,
				},
				CurrentAddress: 0,
			},
			IReproduction: &reproductions.BuddingReproduction{},
			Energy:        22,
			Generation:    0,
		}
		result = append(result, agent)
	}
	return result
}
