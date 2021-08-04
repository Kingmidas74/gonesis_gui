package main

import (
	preparedCommands "github.com/Kingmidas74/gonesis_engine/commands"
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/core"
	"github.com/Kingmidas74/gonesis_engine/core/agents"
	"github.com/Kingmidas74/gonesis_engine/core/commands"
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

func GetTerrain(settings TerrainSettings, cells []contracts.ICell, width, height int) contracts.ITerrain {

	var terrain contracts.ITerrain

	baseTerrain := terrains.Terrain{
		Cells:  cells,
		Width:  width,
		Height: height,
	}

	switch settings.terrainType {
	case 0:
		terrain = &terrains.MooreTerrain{
			baseTerrain,
		}
		break
	case 1:
		terrain = &terrains.NeumannTerrain{
			baseTerrain,
		}
		break
	case 2:
		terrain = &terrains.HexTerrain{
			baseTerrain,
		}
		break
	}

	return terrain
}

func GetAgents(agentsCount int, settings ReproductionSettings) []contracts.IAgent {
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
			Energy:     int(settings.defaultEnergyVolume),
			Generation: 0,
		}
		switch settings.reproductionType {
		case 0:
			agent.IReproduction = &reproductions.BuddingReproduction{
				ReproductionPower:   int(settings.buddingReproductionSettings.reproductionPower),
				MutationProbability: int(settings.buddingReproductionSettings.mutationProbability),
			}
			break
		case 1:
			agent.IReproduction = &reproductions.MitosisReproduction{
				ReproductionPower:   int(settings.mitosisReproductionSettings.reproductionPower),
				MutationProbability: int(settings.mitosisReproductionSettings.mutationProbability),
				GenerationPower:     int(settings.mitosisReproductionSettings.generationCapacity),
			}
			break
		}
		result = append(result, agent)
	}
	return result
}
