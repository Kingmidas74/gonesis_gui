package models

type ReproductionSettings struct {
	ReproductionType            int
	DefaultEnergyVolume         int32
	MitosisReproductionSettings MitosisReproductionSettings
	BuddingReproductionSettings BuddingReproductionSettings
}

type MitosisReproductionSettings struct {
	MutationProbability int32
	ReproductionPower   int32
	GenerationCapacity  int32
}

type BuddingReproductionSettings struct {
	MutationProbability int32
	ReproductionPower   int32
}

type WorldSettings struct {
	AgentsCount int32
}

type TerrainSettings struct {
	TerrainType        int
	Width              int32
	Height             int32
	OrganicProbability int32
}

type EvolutionSettings struct {
	WorldSettings        WorldSettings
	TerrainSettings      TerrainSettings
	ReproductionSettings ReproductionSettings
}
