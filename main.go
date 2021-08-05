package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	workspace := WorkspacePresenter{
		Width:  1920,
		Height: 1080,
	}

	workspace.Init("Gonesis", WorkspaceController{
		Settings: EvolutionSettings{
			worldSettings:   WorldSettings{agentsCount: 1},
			terrainSettings: TerrainSettings{terrainType: 0, organicProbability: 50},
			reproductionSettings: ReproductionSettings{
				reproductionType:    0,
				defaultEnergyVolume: 22,
				mitosisReproductionSettings: MitosisReproductionSettings{
					mutationProbability: 0,
					reproductionPower:   20,
					generationCapacity:  2,
				},
				buddingReproductionSettings: BuddingReproductionSettings{
					mutationProbability: 0,
					reproductionPower:   20,
				},
			},
		},
	}).Start()
}
