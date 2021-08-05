package presenters

import (
	g "github.com/AllenDang/giu"
	"gonesis_gui/controllers"
	"gonesis_gui/models"
	"image"
	"image/color"
)

type WorkspacePresenter struct {
	Width  int
	Height int
	Title  string

	CurrentController controllers.WorkspaceController

	CurrentModel *models.WorkspaceModel

	currentWindow *g.MasterWindow
}

func (this *WorkspacePresenter) Init(title string, controller controllers.WorkspaceController) *WorkspacePresenter {
	this.Title = title
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
					g.InputInt(&this.CurrentModel.Settings.WorldSettings.AgentsCount),
				),
				g.Row(
					g.Label("Start energy"),
					g.InputInt(&this.CurrentModel.Settings.ReproductionSettings.DefaultEnergyVolume),
				),
			),
			g.TabItem("Terrain").Layout(
				g.TreeNode("Generate").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.RadioButton("Moore", this.CurrentModel.Settings.TerrainSettings.TerrainType == 0).OnChange(func() { go this.CurrentController.OnChangeTerrainTypeHandler(0) }),
						g.RadioButton("Neumann", this.CurrentModel.Settings.TerrainSettings.TerrainType == 1).OnChange(func() { go this.CurrentController.OnChangeTerrainTypeHandler(1) }),
						g.RadioButton("Hex", this.CurrentModel.Settings.TerrainSettings.TerrainType == 2).OnChange(func() { go this.CurrentController.OnChangeTerrainTypeHandler(2) }),
					),
					g.Row(
						g.Label("Width"),
						g.InputInt(&this.CurrentModel.Settings.TerrainSettings.Width),

						g.Label("Height"),
						g.InputInt(&this.CurrentModel.Settings.TerrainSettings.Height),
					),
					g.Row(
						g.Label("OrganicProbability"),
						g.InputInt(&this.CurrentModel.Settings.TerrainSettings.OrganicProbability),
					),
					g.Button("Generate").OnClick(func() {
						go this.CurrentController.OnGenerateTerrainHandler(true)
					}),
				),
				g.TreeNode("FromFile").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.InputText(&this.CurrentModel.TerrainFilePath),
						g.Button("Select...").OnClick(this.CurrentController.OnLoadTerrainHandler),
					),
				),
			),
			g.TabItem("Reproduction").Layout(
				g.TreeNode("Budding").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.Label("Reproduction power"),
						g.InputInt(&this.CurrentModel.Settings.ReproductionSettings.BuddingReproductionSettings.ReproductionPower),
					),
					g.Row(
						g.Label("Mutation probability"),
						g.InputInt(&this.CurrentModel.Settings.ReproductionSettings.BuddingReproductionSettings.MutationProbability),
					),
				),
				g.TreeNode("Mitosis").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.Label("Reproduction power"),
						g.InputInt(&this.CurrentModel.Settings.ReproductionSettings.MitosisReproductionSettings.ReproductionPower),
					),
					g.Row(
						g.Label("Mutation probability"),
						g.InputInt(&this.CurrentModel.Settings.ReproductionSettings.MitosisReproductionSettings.MutationProbability),
					),
					g.Row(
						g.Label("Generation capacity"),
						g.InputInt(&this.CurrentModel.Settings.ReproductionSettings.MitosisReproductionSettings.GenerationCapacity),
					),
				),
			),
			g.TabItem("Agent").Layout(
				g.TreeNode("Reproduction").Flags(g.TreeNodeFlagsCollapsingHeader|g.TreeNodeFlagsDefaultOpen).Layout(
					g.Row(
						g.RadioButton("Budding", this.CurrentModel.Settings.ReproductionSettings.ReproductionType == 0).OnChange(func() { go this.CurrentController.OnChangeReproductionTypeHandler(0) }),
						g.RadioButton("Mitosis", this.CurrentModel.Settings.ReproductionSettings.ReproductionType == 1).OnChange(func() { go this.CurrentController.OnChangeReproductionTypeHandler(1) }),
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
			if this.CurrentModel.Texture != nil {
				canvas.AddImage(this.CurrentModel.Texture, image.Pt(0, 0), image.Pt(1920, 1080))
			}
		}),
	)
	return &layout
}

func (this *WorkspacePresenter) loop() {

	g.SingleWindowWithMenuBar().Layout(
		g.MenuBar().Layout(
			g.Menu(this.Title).Layout(
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
