package main

import (
	g "github.com/AllenDang/giu"
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/core/primitives"
	"github.com/Kingmidas74/gonesis_engine/core/terrains"
	"github.com/Kingmidas74/gonesis_engine/core/world"
	"github.com/Kingmidas74/gonesis_engine/utils"
	"os"
	"time"
)

type WorkspaceController struct {
	Settings EvolutionSettings
	Texture  *g.Texture

	currentTerrain  contracts.ITerrain
	terrainFilePath string
}

func (this *WorkspaceController) OnExitHandler() {
	os.Exit(0)
}

func (this *WorkspaceController) OnRunEvolutionHandler() {
	if this.currentTerrain == nil {
		this.generateTerrain(false)
	}
	currentWorld := world.World{
		this.currentTerrain,
	}
	go currentWorld.Action(1, func(terrain contracts.ITerrain, currentDay int) {

		defer time.AfterFunc(time.Duration(100)*time.Millisecond, func() {
			this.updateTexture(terrain)
		}).Stop()
		time.Sleep(100 * time.Millisecond)
	})
	this.currentTerrain = nil
}

func (this *WorkspaceController) OnGenerateTerrainHandler(withDraw bool) {
	this.generateTerrain(withDraw)
}

func (this *WorkspaceController) OnLoadTerrainHandler() {
	this.currentTerrain = GetTerrainFromFile(this.terrainFilePath, this.Settings.terrainSettings)

	this.Settings.terrainSettings.width = int32(this.currentTerrain.GetWidth())
	this.Settings.terrainSettings.height = int32(this.currentTerrain.GetHeight())

	this.placeAgents(GetAgents(int(this.Settings.worldSettings.agentsCount), this.Settings.reproductionSettings))
	go this.updateTexture(this.currentTerrain)
}

func (this *WorkspaceController) generateTerrain(withDraw bool) {
	cells := make([]contracts.ICell, 0)
	for y := 0; y < int(this.Settings.terrainSettings.height); y++ {
		for x := 0; x < int(this.Settings.terrainSettings.width); x++ {
			currentCell := &terrains.Cell{
				Coords: primitives.Coords{
					X: x,
					Y: y,
				},
				CellType: contracts.EmptyCell,
				Cost:     0,
				Agent:    nil,
			}
			if utils.RandomIntBetween(0, 100) > int(this.Settings.terrainSettings.organicProbability) {
				currentCell.SetCellType(contracts.OrganicCell)
				currentCell.SetCost(utils.RandomIntBetween(-20, 20))
			}
			cells = append(cells, currentCell)
		}
	}
	this.currentTerrain = GetTerrain(this.Settings.terrainSettings, cells, int(this.Settings.terrainSettings.width), int(this.Settings.terrainSettings.height))
	this.placeAgents(GetAgents(int(this.Settings.worldSettings.agentsCount), this.Settings.reproductionSettings))
	if withDraw {
		this.updateTexture(this.currentTerrain)
	}
}

func (this *WorkspaceController) updateTexture(terrain contracts.ITerrain) {
	this.Texture, _ = g.NewTextureFromRgba(DrawFrame(terrain, 100))
}

func (this *WorkspaceController) placeAgents(currentAgents []contracts.IAgent) {
	placedChildrenCount := 0

	for i := 0; i < len(this.currentTerrain.GetCells()) && placedChildrenCount < len(currentAgents); i++ {
		if this.currentTerrain.GetCells()[i].GetCellType() == contracts.EmptyCell {
			currentAgents[placedChildrenCount].SetX(this.currentTerrain.GetCells()[i].GetX())
			currentAgents[placedChildrenCount].SetY(this.currentTerrain.GetCells()[i].GetY())
			this.currentTerrain.GetCells()[i].SetCellType(contracts.LockedCell)
			this.currentTerrain.GetCells()[i].SetAgent(currentAgents[placedChildrenCount])
			placedChildrenCount++
		}
	}

}
