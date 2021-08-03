package main

import (
	g "github.com/AllenDang/giu"
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/core/world"
	"image"
	"image/color"
	"os"
	"strconv"
	"time"
)

const (
	title string = "Gonesis"
)

type EvolutionSettings struct {
	agentsCount string
}

type Workspace struct {
	Width  int
	Height int

	currentWindow *g.MasterWindow

	texture      *g.Texture
	currentWorld contracts.IWorld

	settings EvolutionSettings
}

func (this *Workspace) initWorld(settings EvolutionSettings) contracts.IWorld {
	agentsCount, _ := strconv.Atoi(settings.agentsCount)
	currentAgents := GetAgents(agentsCount)

	terrain := GetTerrain(currentAgents)

	return &world.World{
		terrain,
	}
}

func (this *Workspace) Init() {

	this.currentWindow = g.NewMasterWindow(title, this.Width, this.Height, g.MasterWindowFlagsNotResizable)
	this.settings = EvolutionSettings{agentsCount: strconv.Itoa(1)}
}

func (this *Workspace) Start() {
	this.currentWindow.Run(this.loop)
}

func (this *Workspace) runEvolution() {
	this.currentWorld = this.initWorld(this.settings)
	go this.currentWorld.Action(1, func(terrain contracts.ITerrain, currentDay int) {
		img := DrawFrame(terrain, 100)

		defer time.AfterFunc(time.Duration(100)*time.Millisecond, func() {
			this.texture, _ = g.NewTextureFromRgba(img)
		}).Stop()

		time.Sleep(100 * time.Millisecond)
	})
}

func (this *Workspace) drawControls() *g.Layout {
	return &g.Layout{
		g.Row(
			g.Label("Agents count"),
			g.InputText(&this.settings.agentsCount),
		),
		g.Row(
			g.Label("Start simulation"),
			g.Style().
				SetColor(g.StyleColorText, color.RGBA{0x36, 0x74, 0xD5, 255}).
				To(
					g.ArrowButton("Start simulation", g.DirectionRight).OnClick(this.runEvolution),
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
			g.Menu("World").Layout(
				// You could add any kind of widget here, not just menu item.
				g.Menu("Terrain").Layout(
					g.MenuItem("Load from file..."),
					g.MenuItem("Generate"),
				),
				g.Menu("Commands").Layout(
					g.MenuItem("Show"),
					g.MenuItem("Edit"),
					g.MenuItem("Add"),
				),
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
