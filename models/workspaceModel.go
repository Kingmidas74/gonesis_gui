package models

import (
	g "github.com/AllenDang/giu"
	"github.com/Kingmidas74/gonesis_engine/contracts"
)

type WorkspaceModel struct {
	Settings EvolutionSettings
	Texture  *g.Texture

	CurrentTerrain  contracts.ITerrain
	TerrainFilePath string
}
