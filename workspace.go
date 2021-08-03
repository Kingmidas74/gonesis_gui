package main

import (
	g "github.com/AllenDang/giu"
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"image"
	"time"
)

const (
	title string = "Gonesis"
)

type Workspace struct {
	Width  int
	Height int

	currentWindow *g.MasterWindow

	texture      *g.Texture
	currentWorld contracts.IWorld
}

func (this *Workspace) Init(world contracts.IWorld) {
	this.currentWorld = world
	this.currentWindow = g.NewMasterWindow(title, this.Width, this.Height, g.MasterWindowFlagsNotResizable)
}

func (this *Workspace) Start() {
	this.currentWindow.Run(this.loop)
}

func (this *Workspace) runEvolution() {

	go this.currentWorld.Action(1, func(terrain contracts.ITerrain, currentDay int) {
		img := DrawFrame(terrain, 100)

		defer time.AfterFunc(time.Duration(1)*time.Second, func() {
			this.texture, _ = g.NewTextureFromRgba(img)
		}).Stop()

		time.Sleep(1 * time.Second)
	})
}

func (this *Workspace) drawControls() *g.Layout {
	layout := g.Layout{}
	layout = append(layout, g.Label("Hello there"),
		g.Button("Start simulation").OnClick(this.runEvolution))
	return &layout
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

func (this *Workspace) loop() {
	controls := this.drawControls()
	canvas := this.drawCanvas()

	g.SingleWindow().Layout(g.SplitLayout(g.DirectionHorizontal, true, 100, controls, canvas))
}
