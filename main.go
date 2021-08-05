package main

import (
	"gonesis_gui/controllers"
	"gonesis_gui/models"
	"gonesis_gui/presenters"
	"gonesis_gui/services"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	model := &models.WorkspaceModel{
		Settings: models.EvolutionSettings{
			WorldSettings:   models.WorldSettings{AgentsCount: 1},
			TerrainSettings: models.TerrainSettings{TerrainType: 0, OrganicProbability: 50},
			ReproductionSettings: models.ReproductionSettings{
				ReproductionType:    0,
				DefaultEnergyVolume: 22,
				MitosisReproductionSettings: models.MitosisReproductionSettings{
					MutationProbability: 0,
					ReproductionPower:   20,
					GenerationCapacity:  2,
				},
				BuddingReproductionSettings: models.BuddingReproductionSettings{
					MutationProbability: 0,
					ReproductionPower:   20,
				},
			},
		},
		ScaleTextureValue: 1,
		CurrentTerrain:    nil,
		TerrainFilePath:   "",
		Texture:           nil,
	}

	currentWorkspace := presenters.WorkspacePresenter{
		Width:        1920,
		Height:       1080,
		CurrentModel: model,
	}

	currentWorkspace.
		Init("Gonesis", controllers.WorkspaceController{
			GeneratorService: &services.GeneratorService{},
			CurrentModel:     model,
		}).
		Start()
}
