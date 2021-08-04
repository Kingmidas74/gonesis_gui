package main

import (
	"bufio"
	"fmt"
	g "github.com/AllenDang/giu"
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/core/primitives"
	"github.com/Kingmidas74/gonesis_engine/core/terrains"
	"github.com/Kingmidas74/gonesis_engine/core/world"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	title string = "Gonesis"
)

type Workspace struct {
	Width  int
	Height int

	currentWindow *g.MasterWindow

	texture        *g.Texture
	currentTerrain contracts.ITerrain

	settings        EvolutionSettings
	terrainFilePath string
}

func (this *Workspace) placeAgents(currentAgents []contracts.IAgent) {
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

func (this *Workspace) Init() {

	this.terrainFilePath = "$HOME"
	this.currentWindow = g.NewMasterWindow(title, this.Width, this.Height, g.MasterWindowFlagsMaximized)
	this.settings = EvolutionSettings{
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
	}
}

func (this *Workspace) Start() {
	this.currentWindow.Run(this.loop)
}

func (this *Workspace) generateTerrain(withDraw bool) {
	cells := make([]contracts.ICell, 0)
	this.currentTerrain = GetTerrain(this.settings.terrainSettings, cells, 0, 0)
	this.placeAgents(GetAgents(int(this.settings.worldSettings.agentsCount), this.settings.reproductionSettings))
	if withDraw {
		this.texture, _ = g.NewTextureFromRgba(DrawFrame(this.currentTerrain, 100))
	}
}

func Readln(r *bufio.Reader) (string, error) {
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

func (this *Workspace) loadTerrainFromFile(filePath string) {

	cells := make([]contracts.ICell, 0)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	s, e := Readln(r)
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
		s, e = Readln(r)
		currentRowIndex++
	}

	this.currentTerrain = GetTerrain(this.settings.terrainSettings, cells, len(cells)/currentRowIndex, currentRowIndex)
	this.placeAgents(GetAgents(int(this.settings.worldSettings.agentsCount), this.settings.reproductionSettings))
	go func() {
		this.texture, _ = g.NewTextureFromRgba(DrawFrame(this.currentTerrain, 100))
	}()
}

func (this *Workspace) runEvolution() {
	if this.currentTerrain == nil {
		this.generateTerrain(false)
	}
	currentWorld := world.World{
		this.currentTerrain,
	}
	go currentWorld.Action(1, func(terrain contracts.ITerrain, currentDay int) {
		img := DrawFrame(terrain, 100)

		defer time.AfterFunc(time.Duration(100)*time.Millisecond, func() {
			this.texture, _ = g.NewTextureFromRgba(img)
		}).Stop()
		time.Sleep(100 * time.Millisecond)
	})
	this.currentTerrain = nil
}

func (this *Workspace) drawControls() *g.Layout {
	return &g.Layout{
		g.Row(
			g.Label("Start simulation"),
			g.Style().
				SetColor(g.StyleColorText, color.RGBA{0x36, 0x74, 0xD5, 255}).
				To(
					g.ArrowButton("Start simulation", g.DirectionRight).OnClick(this.runEvolution),
				),
		),
		g.TabBar().TabItems(
			g.TabItem("World").Layout(
				g.Row(
					g.Label("Agents count"),
					g.InputInt(&this.settings.worldSettings.agentsCount),
				),
				g.Row(
					g.Label("Start energy"),
					g.InputInt(&this.settings.reproductionSettings.defaultEnergyVolume),
				),
			),
			g.TabItem("Terrain").Layout(
				g.TreeNode("Generate").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.RadioButton("Moore", this.settings.terrainSettings.terrainType == 0).OnChange(func() { this.settings.terrainSettings.terrainType = 0 }),
						g.RadioButton("Neumann", this.settings.terrainSettings.terrainType == 1).OnChange(func() { this.settings.terrainSettings.terrainType = 1 }),
						g.RadioButton("Hex", this.settings.terrainSettings.terrainType == 2).OnChange(func() { this.settings.terrainSettings.terrainType = 2 }),
					),
					g.Row(
						g.Label("OrganicProbability"),
						g.InputInt(&this.settings.terrainSettings.organicProbability),
					),
					g.Button("Generate").OnClick(func() {
						go this.generateTerrain(true)
					}),
				),
				g.TreeNode("FromFile").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.InputText(&this.terrainFilePath),
						g.Button("Select...").OnClick(func() {
							this.loadTerrainFromFile(this.terrainFilePath)
						}),
					),
				),
			),
			g.TabItem("Reproduction").Layout(
				g.TreeNode("Budding").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.Label("Reproduction power"),
						g.InputInt(&this.settings.reproductionSettings.buddingReproductionSettings.reproductionPower),
					),
					g.Row(
						g.Label("Mutation probability"),
						g.InputInt(&this.settings.reproductionSettings.buddingReproductionSettings.mutationProbability),
					),
				),
				g.TreeNode("Mitosis").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.Label("Reproduction power"),
						g.InputInt(&this.settings.reproductionSettings.mitosisReproductionSettings.reproductionPower),
					),
					g.Row(
						g.Label("Mutation probability"),
						g.InputInt(&this.settings.reproductionSettings.mitosisReproductionSettings.mutationProbability),
					),
					g.Row(
						g.Label("Generation capacity"),
						g.InputInt(&this.settings.reproductionSettings.mitosisReproductionSettings.generationCapacity),
					),
				),
			),
			g.TabItem("Agent").Layout(
				g.TreeNode("Reproduction").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.RadioButton("Budding", this.settings.reproductionSettings.reproductionType == 0).OnChange(func() { this.settings.reproductionSettings.reproductionType = 0 }),
						g.RadioButton("Mitosis", this.settings.reproductionSettings.reproductionType == 1).OnChange(func() { this.settings.reproductionSettings.reproductionType = 1 }),
					),
				),
			),
		),
	}
}

func (this *Workspace) drawCanvas() *g.Layout {
	layout := g.Layout{}
	layout = append(layout,
		g.Custom(func() {
			canvas := g.GetCanvas()
			if this.texture != nil {
				canvas.AddImage(this.texture, image.Pt(0, 0), image.Pt(1920, 1080))
			}
		}),
	)
	return &layout
}

func (this *Workspace) exit() {
	os.Exit(0)
}

func (this *Workspace) loop() {

	g.SingleWindowWithMenuBar().Layout(
		g.MenuBar().Layout(
			g.Menu("Gonesis").Layout(
				g.MenuItem("Load"),
				g.MenuItem("Save"),
				g.MenuItem("Close").OnClick(this.exit),
			),
		),
		g.SplitLayout(g.DirectionHorizontal, true, 1730,
			g.SplitLayout(g.DirectionVertical, true, 900,
				this.drawCanvas(),
				g.Layout{},
			),
			g.Layout{
				this.drawControls(),
			},
		),
	)
}
