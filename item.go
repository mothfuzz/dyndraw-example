package main

import (
	"encoding/json"
	"encoding/xml"
	"io/fs"
	"path/filepath"

	"github.com/mothfuzz/dyndraw/framework/actors"
	"github.com/mothfuzz/dyndraw/framework/render"
	"github.com/mothfuzz/dyndraw/framework/transform"
)

type Item struct {
	Name        string
	Description string
	Icon        string
	Sprite      string
	transform.Transform
	Collider
}

func (i *Item) Init() {
	if i.Transform.GetScaleV().Z() == 0 {
		i.Transform = transform.Origin2D(16, 16)
	}
	i.Collider = NewBoundingSphere(0.5)
	i.IgnoreRaycast = true
}
func (i *Item) Destroy() {}
func (i *Item) Update() {
	//listeners for Item can pick me up
	actors.AllListeners(Item{}, func(a actors.Actor) {
		//if DistanceSqr(a, i) <= 16*16 {
		if Overlaps(a, i) {
			actors.Send(a, *i)
			actors.Destroy(i)
		}
	})
}

func (i *Item) Draw() {
	if i.Sprite != "" {
		render.DrawSprite(i.Sprite, i.Transform.Mat4())
	} else {
		render.DrawSprite(i.Icon, i.Transform.Mat4())
	}
}

var itemDictionary = map[string]Item{}

func LoadItemDictionary() {
	err := fs.WalkDir(render.Resources, "resources/items", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		var item Item
		file, err := fs.ReadFile(render.Resources, path)
		if err != nil {
			return err
		}
		switch filepath.Ext(path) {
		case ".json":
			err = json.Unmarshal(file, &item)
			if err != nil {
				return err
			}
		case ".xml":
			err = xml.Unmarshal(file, &item)
			if err != nil {
				return err
			}
		default:
			return nil
		}
		//itemName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		itemDictionary[filepath.Base(path)] = item
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func ItemDictionary(name string) *Item {
	//just so we don't modify the actual items *in* the dictionary
	newItem := new(Item)
	*newItem = itemDictionary[name]
	return newItem
}
