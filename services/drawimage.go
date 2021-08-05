package services

import (
	"github.com/Kingmidas74/gonesis_engine/contracts"
	"github.com/Kingmidas74/gonesis_engine/utils"
	"image"
	"image/color"
	"image/draw"
)

func DrawFrame(terrain contracts.ITerrain, zoom int) *image.RGBA {

	m := image.NewRGBA(image.Rect(0, 0, (terrain.GetWidth()+1)*zoom, (terrain.GetHeight()+1)*zoom))
	foodCellColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	poisonCellColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	freeCellColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	currentColor := freeCellColor

	for currentLatitude := 0; currentLatitude < terrain.GetWidth(); currentLatitude++ {
		for currentLongitude := 0; currentLongitude < terrain.GetHeight(); currentLongitude++ {
			if currentCell := terrain.GetCell(currentLatitude, currentLongitude); currentCell.GetCellType() == contracts.EmptyCell {
				currentColor = freeCellColor
			} else if currentCell.GetCellType() == contracts.LockedCell {
				currentColor = color.RGBA{R: 0, G: 0, B: 255, A: uint8(utils.ModLikePython(255+currentCell.GetAgent().GetGeneration()*10, 255))}
			} else if currentCell.GetCellType() == contracts.ObstacleCell {
				currentColor = color.RGBA{R: 50, G: 50, B: 50, A: 255}
			} else if currentCell.GetCellType() == contracts.OrganicCell {
				if currentCell.GetCost() > 0 {
					currentColor = foodCellColor
				} else {
					currentColor = poisonCellColor
				}
			}
			draw.Draw(m, image.Rect(currentLatitude*zoom, currentLongitude*zoom, currentLatitude*zoom+zoom, currentLongitude*zoom+zoom), &image.Uniform{C: currentColor}, image.Point{}, draw.Src)
		}
	}

	return m
}
