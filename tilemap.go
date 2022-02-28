package main

import (
	"math"

	"github.com/mothfuzz/letsgo/collision"
	"github.com/mothfuzz/letsgo/render"
	"github.com/mothfuzz/letsgo/transform"
	. "github.com/mothfuzz/letsgo/vecmath"
)

type TileSet struct {
	Image  string
	W, H   int
	TW, TH int
}

//-
func NewXPlane(x, y, l float32) collision.Plane {
	a := Vec3{x, y, -l}
	b := Vec3{x + l/2, y, 0}
	c := Vec3{x - l/2, y, 0}
	return collision.NewPlaneAt(x, y, 0, a, b, c)
}

//|
func NewYPlane(x, y, l float32) collision.Plane {
	a := Vec3{x, y, -l}
	b := Vec3{x, y - l/2, 0}
	c := Vec3{x, y + l/2, 0}
	return collision.NewPlaneAt(x, y, 0, a, b, c)
}

type TileMap struct {
	TileSet
	Data [][]uint8
	collision.Collider
	render.SpriteAnimation
	tileTransform transform.Transform
}

func tileOccupied(t *TileMap, x, y int, mask []uint8) bool {
	if x < 0 {
		return false
	}
	if x >= len(t.Data[0]) {
		return false
	}
	if y < 0 {
		return false
	}
	if y >= len(t.Data) {
		return false
	}
	for _, val := range mask {
		if t.Data[y][x] == val {
			return true
		}
	}
	return false
}

func (t *TileMap) offsets() (float32, float32) {
	xOffset := float32(t.TileSet.TW) / 2.0
	yOffset := 400.0 + float32(t.TileSet.TH)/2.0 - float32(len(t.Data)*t.TileSet.TH)
	return xOffset, yOffset
}

func (t *TileMap) Init() {
	t.tileTransform = transform.Origin2D()
	t.SpriteAnimation = render.SpriteAnimation{
		Frames: []render.Frame{},
		Tags:   map[string][]int{"tiles": {}},
	}
	for i := 0; i < t.TileSet.H; i++ {
		for j := 0; j < t.TileSet.W; j++ {
			t.SpriteAnimation.Frames = append(t.SpriteAnimation.Frames, render.Frame{
				X: float32(j) / float32(t.TileSet.W),
				Y: float32(i) / float32(t.TileSet.H),
				W: 1.0 / float32(t.TileSet.W),
				H: 1.0 / float32(t.TileSet.H),
			})
			t.SpriteAnimation.Tags["tiles"] = append(t.SpriteAnimation.Tags["tiles"], i*t.TileSet.W+j)
		}
	}
	if t.Planes == nil {
		for i := 0; i < len(t.Data); i++ {
			for j := 0; j < len(t.Data[i]); j++ {
				xOffset, yOffset := t.offsets()
				x := float32(j*t.TileSet.TW) + xOffset
				y := float32(i*t.TileSet.TH) + yOffset
				w := float32(t.TileSet.TW)
				h := float32(t.TileSet.TH)
				switch t.Data[i][j] {
				case 1:
					t.Planes = append(t.Planes, collision.NewPlaneAt(x, y, 0, Vec3{x, y, -w}, Vec3{x + w/2, y - h/2, 0}, Vec3{x - w/2, y + h/2, 0}))
				case 2:
					t.Planes = append(t.Planes, collision.NewPlaneAt(x, y, 0, Vec3{x, y, -w}, Vec3{x + w/2, y + h/2, 0}, Vec3{x - w/2, y - h/2, 0}))
				case 3:
					if !tileOccupied(t, j, i-1, []uint8{1, 2, 3}) {
						t.Planes = append(t.Planes, NewXPlane(x, y-h/2, h))
					}
					if !tileOccupied(t, j, i+1, []uint8{1, 2, 3}) {
						t.Planes = append(t.Planes, NewXPlane(x, y+h/2, h))
					}
					if !tileOccupied(t, j-1, i, []uint8{1, 3}) {
						t.Planes = append(t.Planes, NewYPlane(x-w/2, y, w))
					}
					if !tileOccupied(t, j+1, i, []uint8{2, 3}) {
						t.Planes = append(t.Planes, NewYPlane(x+w/2, y, w))
					}
				}
			}
		}
	}

}
func (t *TileMap) Update()  {}
func (t *TileMap) Destroy() {}

func (t *TileMap) Draw() {
	xOffset, yOffset := t.offsets()
	for i, row := range t.Data {
		for j, tile := range row {
			if tile != 0 {
				t.tileTransform.SetPosition2D(
					float32(j*t.TileSet.TW)+xOffset,
					float32(i*t.TileSet.TH)+yOffset,
				)
				render.DrawSpriteAnimated(t.TileSet.Image, t.tileTransform.Mat4(), t.SpriteAnimation.GetTexCoords("tiles", int(tile)))
			}
		}
	}
	/*for _, p := range t.Planes {
		t := transform.Origin2D(4, 4)
		t.SetPosition(p.origin.X(), p.origin.Y(), p.origin.Z())
		render.DrawSprite("pointg.png", t.Mat4())
		t.SetPosition(p.points[0].X(), p.points[0].Y(), p.points[0].Z())
		render.DrawSprite("point.png", t.Mat4())
		t.SetPosition(p.points[1].X(), p.points[1].Y(), p.points[1].Z())
		render.DrawSprite("point.png", t.Mat4())
		t.SetPosition(p.points[2].X(), p.points[2].Y(), p.points[2].Z())
		render.DrawSprite("point.png", t.Mat4())
		out := p.origin.Add(p.normal.Mul(8))
		t.SetPosition(out.X(), out.Y(), out.Z())
		render.DrawSprite("pointg.png", t.Mat4())
	}*/
}

func CheckTile(t *transform.Transform, tm *TileMap, w, h float32, mask []uint8) bool {
	xOffset, yOffset := tm.offsets()
	pos := t.GetPositionV().Vec2().Sub(Vec2{xOffset, yOffset})
	//since we're using center coords we have to offset it
	//by a half-width for the rounding to work
	//otherwise pretty straightforward
	leftTile := int(math.Floor(float64(pos.X()-w/2)/float64(tm.TW) + 0.5))
	topTile := int(math.Floor(float64(pos.Y()-h/2)/float64(tm.TH) + 0.5))
	rightTile := int(math.Ceil(float64(pos.X()+w/2)/float64(tm.TW)+0.5)) - 1
	bottomTile := int(math.Ceil(float64(pos.Y()+h/2)/float64(tm.TH)+0.5)) - 1

	for y := topTile; y <= bottomTile; y++ {
		for x := leftTile; x <= rightTile; x++ {
			if tileOccupied(tm, x, y, mask) {
				return true
			}
		}
	}
	return false
}
func MoveAgainstTiles(t *transform.Transform, tm *TileMap, xspeed, yspeed float32, w, h float32) (float32, float32) {
	pos := t.GetPositionV()
	tw := float32(tm.TW)
	th := float32(tm.TH)
	t.Translate2D(xspeed, 0)
	if CheckTile(t, tm, w, h, []uint8{3}) {
		x, y := t.GetPositionV().Vec2().Elem()
		if xspeed > 0 {
			x = float32(math.Round(float64((pos.X()+w/2)/tw)))*tw - w/2
		} else {
			x = float32(math.Round(float64((pos.X()-w/2)/tw)))*tw + w/2
		}
		xspeed = 0
		t.SetPosition2D(x, y)
	}
	t.Translate2D(0, yspeed)
	if CheckTile(t, tm, w, h, []uint8{3}) {
		x, y := t.GetPositionV().Vec2().Elem()
		if yspeed > 0 {
			y = float32(math.Round(float64((pos.Y()+h/2)/th)))*th - h/2
		} else {
			y = float32(math.Round(float64((pos.Y()-h/2)/th)))*th + h/2
		}
		yspeed = 0
		t.SetPosition2D(x, y)
	}
	return xspeed, yspeed
}
