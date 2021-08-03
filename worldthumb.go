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

func GetTerrain() contracts.ITerrain {
	terrain := terrains.MooreTerrain{
		terrains.Terrain{
			Cells:  make([]contracts.ICell, 0),
			Width:  10,
			Height: 5,
		},
	}

	for y := 0; y < terrain.Height; y++ {
		for x := 0; x < terrain.Width; x++ {
			currentCell := terrains.Cell{
				Coords: primitives.Coords{
					X: x,
					Y: y,
				},
				CellType: contracts.EmptyCell,
				Cost:     0,
			}
			terrain.Cells = append(terrain.Cells, &currentCell)
		}
	}

	terrain.Cells[21].SetCellType(contracts.OrganicCell)
	terrain.Cells[21].SetCost(4)

	terrain.Cells[23].SetCellType(contracts.OrganicCell)
	terrain.Cells[23].SetCost(4)

	terrain.Cells[43].SetCellType(contracts.OrganicCell)
	terrain.Cells[43].SetCost(4)

	terrain.Cells[45].SetCellType(contracts.OrganicCell)
	terrain.Cells[45].SetCost(4)

	terrain.Cells[15].SetCellType(contracts.OrganicCell)
	terrain.Cells[15].SetCost(4)

	terrain.Cells[17].SetCellType(contracts.OrganicCell)
	terrain.Cells[17].SetCost(4)

	terrain.Cells[37].SetCellType(contracts.OrganicCell)
	terrain.Cells[37].SetCost(4)

	terrain.Cells[39].SetCellType(contracts.OrganicCell)
	terrain.Cells[39].SetCost(4)

	terrain.Cells[9].SetCellType(contracts.OrganicCell)
	terrain.Cells[9].SetCost(4)

	return &terrain
}

func GetAgent() contracts.IAgent {
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
		Coords: primitives.Coords{
			X: 1,
			Y: 0,
		},
		IReproduction: &reproductions.BuddingReproduction{},
		Energy:        22,
		Generation:    0,
	}
	return agent
}
