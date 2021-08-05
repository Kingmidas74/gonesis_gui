package controllers

import (
	g "github.com/AllenDang/giu"
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/core/world"
	"gonesis_gui/models"
	"gonesis_gui/services"
	"os"
	"time"
)

type WorkspaceController struct {
	CurrentModel     *models.WorkspaceModel
	GeneratorService *services.GeneratorService
}

func (this *WorkspaceController) OnExitHandler() {
	os.Exit(0)
}

func (this *WorkspaceController) OnChangeReproductionTypeHandler(reproductionType int) {
	this.CurrentModel.Settings.ReproductionSettings.ReproductionType = reproductionType
}

func (this *WorkspaceController) OnChangeTerrainTypeHandler(terrainType int) {
	this.CurrentModel.Settings.TerrainSettings.TerrainType = terrainType
}

func (this *WorkspaceController) OnRunEvolutionHandler() {
	if this.CurrentModel.CurrentTerrain == nil {
		this.generateTerrain(false)
	}
	currentWorld := world.World{
		this.CurrentModel.CurrentTerrain,
	}
	go currentWorld.Action(1, func(terrain contracts.ITerrain, currentDay int) {

		defer time.AfterFunc(time.Duration(100)*time.Millisecond, func() {
			this.updateTexture(terrain)
		}).Stop()
		time.Sleep(100 * time.Millisecond)
	})
	this.CurrentModel.CurrentTerrain = nil
}

func (this *WorkspaceController) OnGenerateTerrainHandler(withDraw bool) {
	this.generateTerrain(withDraw)
}

func (this *WorkspaceController) OnLoadTerrainHandler() {
	this.CurrentModel.CurrentTerrain = this.GeneratorService.GetTerrainFromFile(this.CurrentModel.TerrainFilePath, this.CurrentModel.Settings.TerrainSettings)

	this.CurrentModel.Settings.TerrainSettings.Width = int32(this.CurrentModel.CurrentTerrain.GetWidth())
	this.CurrentModel.Settings.TerrainSettings.Height = int32(this.CurrentModel.CurrentTerrain.GetHeight())

	this.CurrentModel.CurrentTerrain = this.GeneratorService.PlaceAgents(this.GeneratorService.GetAgents(int(this.CurrentModel.Settings.WorldSettings.AgentsCount), this.CurrentModel.Settings.ReproductionSettings), this.CurrentModel.CurrentTerrain)
	go this.updateTexture(this.CurrentModel.CurrentTerrain)
}

func (this *WorkspaceController) generateTerrain(withDraw bool) {
	cells := this.GeneratorService.GetCells(this.CurrentModel.Settings.TerrainSettings)
	this.CurrentModel.CurrentTerrain = this.GeneratorService.GetTerrain(this.CurrentModel.Settings.TerrainSettings, cells, int(this.CurrentModel.Settings.TerrainSettings.Width), int(this.CurrentModel.Settings.TerrainSettings.Height))
	this.CurrentModel.CurrentTerrain = this.GeneratorService.PlaceAgents(this.GeneratorService.GetAgents(int(this.CurrentModel.Settings.WorldSettings.AgentsCount), this.CurrentModel.Settings.ReproductionSettings), this.CurrentModel.CurrentTerrain)
	if withDraw {
		this.updateTexture(this.CurrentModel.CurrentTerrain)
	}
}

func (this *WorkspaceController) updateTexture(terrain contracts.ITerrain) {
	this.CurrentModel.Texture, _ = g.NewTextureFromRgba(services.DrawFrame(terrain, this.CurrentModel.ScaleTextureValue))
}
