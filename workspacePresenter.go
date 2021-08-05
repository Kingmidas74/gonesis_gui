package main

import (
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
)

type WorkspacePresenter struct {
	Width             int
	Height            int
	CurrentController WorkspaceController

	currentWindow *g.MasterWindow
}

func (this *WorkspacePresenter) Init(title string, controller WorkspaceController) *WorkspacePresenter {
	this.currentWindow = g.NewMasterWindow(title, this.Width, this.Height, g.MasterWindowFlagsMaximized)
	this.CurrentController = controller
	return this
}

func (this *WorkspacePresenter) Start() {
	this.currentWindow.Run(this.loop)
}

func (this *WorkspacePresenter) drawControls() *g.Layout {
	return &g.Layout{
		g.Row(
			g.Label("Start simulation"),
			g.Style().
				SetColor(g.StyleColorText, color.RGBA{0x36, 0x74, 0xD5, 255}).
				To(
					g.ArrowButton("Start simulation", g.DirectionRight).OnClick(this.CurrentController.OnRunEvolutionHandler),
				),
		),
		g.TabBar().TabItems(
			g.TabItem("World").Layout(
				g.Row(
					g.Label("Agents count"),
					g.InputInt(&this.CurrentController.Settings.worldSettings.agentsCount),
				),
				g.Row(
					g.Label("Start energy"),
					g.InputInt(&this.CurrentController.Settings.reproductionSettings.defaultEnergyVolume),
				),
			),
			g.TabItem("Terrain").Layout(
				g.TreeNode("Generate").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.RadioButton("Moore", this.CurrentController.Settings.terrainSettings.terrainType == 0).OnChange(func() { this.CurrentController.Settings.terrainSettings.terrainType = 0 }),
						g.RadioButton("Neumann", this.CurrentController.Settings.terrainSettings.terrainType == 1).OnChange(func() { this.CurrentController.Settings.terrainSettings.terrainType = 1 }),
						g.RadioButton("Hex", this.CurrentController.Settings.terrainSettings.terrainType == 2).OnChange(func() { this.CurrentController.Settings.terrainSettings.terrainType = 2 }),
					),
					g.Row(
						g.Label("Width"),
						g.InputInt(&this.CurrentController.Settings.terrainSettings.width),

						g.Label("Height"),
						g.InputInt(&this.CurrentController.Settings.terrainSettings.height),
					),
					g.Row(
						g.Label("OrganicProbability"),
						g.InputInt(&this.CurrentController.Settings.terrainSettings.organicProbability),
					),
					g.Button("Generate").OnClick(func() {
						go this.CurrentController.OnGenerateTerrainHandler(true)
					}),
				),
				g.TreeNode("FromFile").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.InputText(&this.CurrentController.terrainFilePath),
						g.Button("Select...").OnClick(this.CurrentController.OnLoadTerrainHandler),
					),
				),
			),
			g.TabItem("Reproduction").Layout(
				g.TreeNode("Budding").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.Label("Reproduction power"),
						g.InputInt(&this.CurrentController.Settings.reproductionSettings.buddingReproductionSettings.reproductionPower),
					),
					g.Row(
						g.Label("Mutation probability"),
						g.InputInt(&this.CurrentController.Settings.reproductionSettings.buddingReproductionSettings.mutationProbability),
					),
				),
				g.TreeNode("Mitosis").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.Label("Reproduction power"),
						g.InputInt(&this.CurrentController.Settings.reproductionSettings.mitosisReproductionSettings.reproductionPower),
					),
					g.Row(
						g.Label("Mutation probability"),
						g.InputInt(&this.CurrentController.Settings.reproductionSettings.mitosisReproductionSettings.mutationProbability),
					),
					g.Row(
						g.Label("Generation capacity"),
						g.InputInt(&this.CurrentController.Settings.reproductionSettings.mitosisReproductionSettings.generationCapacity),
					),
				),
			),
			g.TabItem("Agent").Layout(
				g.TreeNode("Reproduction").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.RadioButton("Budding", this.CurrentController.Settings.reproductionSettings.reproductionType == 0).OnChange(func() { this.CurrentController.Settings.reproductionSettings.reproductionType = 0 }),
						g.RadioButton("Mitosis", this.CurrentController.Settings.reproductionSettings.reproductionType == 1).OnChange(func() { this.CurrentController.Settings.reproductionSettings.reproductionType = 1 }),
					),
				),
			),
		),
	}
}

func (this *WorkspacePresenter) drawCanvas() *g.Layout {
	layout := g.Layout{}
	layout = append(layout,
		g.Custom(func() {
			canvas := g.GetCanvas()
			if this.CurrentController.Texture != nil {
				canvas.AddImage(this.CurrentController.Texture, image.Pt(0, 0), image.Pt(1920, 1080))
			}
		}),
	)
	return &layout
}

func (this *WorkspacePresenter) loop() {

	g.SingleWindowWithMenuBar().Layout(
		g.MenuBar().Layout(
			g.Menu("Gonesis").Layout(
				g.MenuItem("Load"),
				g.MenuItem("Save"),
				g.MenuItem("Close").OnClick(this.CurrentController.OnExitHandler),
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
