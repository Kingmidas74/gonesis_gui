package services

import (
	"bufio"
	"fmt"
	preparedCommands "github.com/Kingmidas74/gonesis_engine/commands"
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/core"
	"github.com/Kingmidas74/gonesis_engine/core/agents"
	"github.com/Kingmidas74/gonesis_engine/core/commands"
	"github.com/Kingmidas74/gonesis_engine/core/primitives"
	"github.com/Kingmidas74/gonesis_engine/core/reproductions"
	"github.com/Kingmidas74/gonesis_engine/core/terrains"
	"github.com/Kingmidas74/gonesis_engine/utils"
	"gonesis_gui/models"
	"os"
	"strconv"
	"strings"
)

type GeneratorService struct {
}

func (this *GeneratorService) GetCommands() map[int]contracts.ICommand {
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

func (this *GeneratorService) Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func (this *GeneratorService) GetTerrainFromFile(filePath string, settings models.TerrainSettings) contracts.ITerrain {
	cells := make([]contracts.ICell, 0)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	s, e := this.Readln(r)
	currentRowIndex := 0
	for e == nil {
		for i, data := range strings.Split(s, ",") {
			currentCell := terrains.Cell{
				Coords: primitives.Coords{
					X: i,
					Y: currentRowIndex,
				},
				CellType: contracts.EmptyCell,
				Cost:     0,
				Agent:    nil,
			}
			if data == "*" {
				currentCell.SetCellType(contracts.ObstacleCell)
			} else {
				weight, _ := strconv.Atoi(data)
				if weight != 0 {
					currentCell.SetCellType(contracts.OrganicCell)
					currentCell.SetCost(weight)
				}
			}
			cells = append(cells, &currentCell)
		}
		s, e = this.Readln(r)
		currentRowIndex++
	}

	return this.GetTerrain(settings, cells, len(cells)/currentRowIndex, currentRowIndex)
}

func (this *GeneratorService) GetTerrain(settings models.TerrainSettings, cells []contracts.ICell, width, height int) contracts.ITerrain {

	var terrain contracts.ITerrain

	baseTerrain := terrains.Terrain{
		Cells:  cells,
		Width:  width,
		Height: height,
	}

	switch settings.TerrainType {
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

func (this *GeneratorService) GetAgents(agentsCount int, settings models.ReproductionSettings) []contracts.IAgent {
	result := make([]contracts.IAgent, 0)

	for i := 0; i < agentsCount; i++ {
		agent := &agents.Agent{
			IBrain: &core.Brain{
				CommandList: commands.CommandList{
					Commands: this.GetCommands(),
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
			Energy:     int(settings.DefaultEnergyVolume),
			Generation: 0,
		}
		switch settings.ReproductionType {
		case 0:
			agent.IReproduction = &reproductions.BuddingReproduction{
				ReproductionPower:   int(settings.BuddingReproductionSettings.ReproductionPower),
				MutationProbability: int(settings.BuddingReproductionSettings.MutationProbability),
			}
			break
		case 1:
			agent.IReproduction = &reproductions.MitosisReproduction{
				ReproductionPower:   int(settings.MitosisReproductionSettings.ReproductionPower),
				MutationProbability: int(settings.MitosisReproductionSettings.MutationProbability),
				GenerationPower:     int(settings.MitosisReproductionSettings.GenerationCapacity),
			}
			break
		}
		result = append(result, agent)
	}
	return result
}

func (this *GeneratorService) PlaceAgents(currentAgents []contracts.IAgent, terrain contracts.ITerrain) contracts.ITerrain {
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

func (this *GeneratorService) GetCells(settings models.TerrainSettings) []contracts.ICell {
	cells := make([]contracts.ICell, 0)
	for y := 0; y < int(settings.Height); y++ {
		for x := 0; x < int(settings.Width); x++ {
			currentCell := &terrains.Cell{
				Coords: primitives.Coords{
					X: x,
					Y: y,
				},
				CellType: contracts.EmptyCell,
				Cost:     0,
				Agent:    nil,
			}
			if utils.RandomIntBetween(0, 100) <= int(settings.OrganicProbability) {
				currentCell.SetCellType(contracts.OrganicCell)
				currentCell.SetCost(utils.RandomIntBetween(-20, 20))
			}
			cells = append(cells, currentCell)
		}
	}
	return cells
}
