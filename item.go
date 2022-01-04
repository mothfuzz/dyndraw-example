package main

import (
	"github.com/mothfuzz/dyndraw/framework/actors"
	"github.com/mothfuzz/dyndraw/framework/render"
	"github.com/mothfuzz/dyndraw/framework/transform"
)

type Item struct {
	Name        string
	Description string
	Icon        string
	transform.Transform
}

func (i *Item) Init() {
	i.Transform = transform.Origin2D(16, 16)
}
func (i *Item) Destroy() {}
func (i *Item) Update() {
	//listeners for Item can pick me up
	actors.AllListeners(Item{}, func(a actors.Actor) {
		if DistanceSqr(a, i) <= 16*16 {
			actors.Send(a, *i)
			actors.Destroy(i)
		}
	})
}

func (i *Item) Draw() {
	render.DrawSprite(i.Icon, i.Transform.Mat4())
}
