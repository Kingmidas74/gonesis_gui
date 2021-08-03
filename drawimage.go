package main

import (
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/utils"
	"image"
	"image/color"
	"image/draw"
)

func DrawFrame(terrain contracts.ITerrain, zoom int) *image.RGBA {

	m := image.NewRGBA(image.Rect(0, 0, (terrain.GetWidth()+1)*zoom, (terrain.GetHeight()+1)*zoom))
	organiccolor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	//agentcolor := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	emptycolor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	currentColor := organiccolor
	//draw.Draw(m, image.Rect(0, 0, terrain.GetWidth()*zoom+zoom, terrain.GetHeight()*zoom+zoom), &image.Uniform{C: organiccolor}, image.Point{}, draw.Src)

	for currentLatitude := 0; currentLatitude < terrain.GetWidth(); currentLatitude++ {
		for currentLongitude := 0; currentLongitude < terrain.GetHeight(); currentLongitude++ {
			if currentCell := terrain.GetCell(currentLatitude, currentLongitude); currentCell.GetCellType() == contracts.EmptyCell {
				currentColor = emptycolor
			} else if currentCell.GetCellType() == contracts.LockedCell {
				currentColor = color.RGBA{R: 0, G: 0, B: 255, A: uint8(utils.ModLikePython(255+currentCell.GetAgent().GetGeneration()*10, 255))}
			} else if currentCell.GetCellType() == contracts.OrganicCell {
				currentColor = organiccolor
			}
			draw.Draw(m, image.Rect(currentLatitude*zoom, currentLongitude*zoom, currentLatitude*zoom+zoom, currentLongitude*zoom+zoom), &image.Uniform{C: currentColor}, image.Point{}, draw.Src)
		}
	}

	return m
}

