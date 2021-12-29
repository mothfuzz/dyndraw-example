package main

import (
	"embed"

	"github.com/mothfuzz/dyndraw/framework/actors"
	"github.com/mothfuzz/dyndraw/framework/app"
	"github.com/mothfuzz/dyndraw/framework/render"
)

//go:embed resources
var Resources embed.FS

func main() {
	render.Resources = Resources
	app.Init()
	defer app.Quit()

	app.SetWindowSize(640, 400)

	t := &TileMap{
		TileSet: TileSet{"tileset.png", 4, 4, 16, 16},
		Data: [][]uint8{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 0, 0, 0, 1, 3, 3, 2, 0, 0},
			{0, 0, 1, 3, 2, 0, 0, 0, 0, 0, 0, 3, 0, 3, 0, 0, 3, 0, 0, 1, 3, 3, 3, 3, 2, 0},
			{1, 3, 3, 3, 3, 3, 2, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 0, 1, 3, 3, 3, 3, 3, 3, 2},
			{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
		},
	}
	CurrentLevel = t
	actors.Spawn(t)
	actors.Spawn(&Player{})

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}
}
