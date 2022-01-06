package main

import (
	"embed"

	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/app"
	"github.com/mothfuzz/letsgo/collision"
	"github.com/mothfuzz/letsgo/input"
	"github.com/mothfuzz/letsgo/render"
	"github.com/mothfuzz/letsgo/transform"
	. "github.com/mothfuzz/letsgo/vecmath"
)

//go:embed resources
var Resources embed.FS

type RayTest struct{}

func (r *RayTest) Init()    {}
func (r *RayTest) Update()  {}
func (r *RayTest) Destroy() {}
func (r *RayTest) Draw() {
	mx, my := input.GetMousePosition()
	startPoint := Vec3{640.0 / 2, 400.0 / 2, 0}
	endPoint := render.RelativeToCamera(mx, my)

	t := transform.Origin2D(4, 4)
	t.SetPosition(640/2, 400/2, -1)
	render.DrawSprite("pointg.png", t.Mat4())

	ray := endPoint.Sub(startPoint).Normalize()
	for _, p := range collision.RayCast(startPoint, ray) {
		t.SetPosition(p.I.X(), p.I.Y(), -1)
		render.DrawSprite("point.png", t.Mat4())
	}
	if hit, ok := collision.RayCastLen(startPoint, ray, 640/2); ok {
		t := transform.Origin2D(4, 4)
		t.SetPosition(hit.I.X(), hit.I.Y(), -2)
		render.DrawSprite("pointg.png", t.Mat4())
	}
}

func main() {
	render.Resources = Resources
	app.Init()
	defer app.Quit()

	app.SetWindowSize(640, 400)

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
	CurrentLevel = []collision.Collider{t.Collider}
	actors.Spawn(t)
	actors.Spawn(&Player{})
	actors.Spawn(&RayTest{})
	actors.SpawnAt(ItemDictionary("thingy.xml"), transform.Location2D(640/2+16, 480/2, 16, 16))
	actors.SpawnAt(ItemDictionary("otherthingy.json"), transform.Location2D(640/2+32, 480/2, 16, 16))

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}
}
