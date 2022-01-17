package main

import (
	"embed"

	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/app"
	"github.com/mothfuzz/letsgo/resources"
	"github.com/mothfuzz/letsgo/transform"
)

//go:embed resources
var Resources embed.FS

func main() {
	resources.Resources = Resources
	app.Init()
	defer app.Quit()

	app.SetWindowSize(640, 400)
	//app.SetFullScreen(true)
	app.SetVSync(false)
	app.SetBackground(0.55, 0.75, 0.95)

	LoadItemDictionary()

	t := &TileMap{
		TileSet: TileSet{"tileset.png", 4, 4, 16, 16},
		Data: [][]uint8{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 3, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 3, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 3, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 3, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 0, 3, 0, 3, 0, 2, 0, 0, 0, 0, 0, 0, 1, 0, 3, 0, 3, 0, 2, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 3, 0, 3, 0, 3, 0, 2, 0, 0, 0, 0, 1, 0, 3, 0, 3, 0, 3, 0, 2, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 0, 3, 0, 3, 0, 3, 0, 3, 0, 2, 0, 0, 1, 0, 3, 0, 3, 0, 3, 0, 3, 0, 2, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 2, 0, 0, 0, 0},
			{0, 0, 0, 1, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 2, 0, 0, 0},
			{0, 0, 1, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 2, 0, 0},
			{0, 1, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 2, 0},
			{1, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 2},
		},
	}
	CurrentLevel = t
	actors.Spawn(t)
	actors.Spawn(&Player{})
	actors.SpawnAt(ItemDictionary("thingy.xml"), transform.Location2D(640/2+16, 480/2, 16, 16))
	actors.SpawnAt(ItemDictionary("otherthingy.json"), transform.Location2D(640/2+32, 480/2, 16, 16))

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}
}
