package main

type ReproductionSettings struct {
	reproductionType            int
	defaultEnergyVolume         int32
	mitosisReproductionSettings MitosisReproductionSettings
	buddingReproductionSettings BuddingReproductionSettings
}

type MitosisReproductionSettings struct {
	mutationProbability int32
	reproductionPower   int32
	generationCapacity  int32
}

type BuddingReproductionSettings struct {
	mutationProbability int32
	reproductionPower   int32
}

type WorldSettings struct {
	agentsCount int32
}

type TerrainSettings struct {
	terrainType        int
	organicProbability int32
}

type EvolutionSettings struct {
	worldSettings        WorldSettings
	terrainSettings      TerrainSettings
	reproductionSettings ReproductionSettings
}
