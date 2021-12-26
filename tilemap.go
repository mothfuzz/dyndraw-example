package main

import (
	"fmt"
	"math"

	. "github.com/go-gl/mathgl/mgl32"
	"github.com/mothfuzz/dyndraw/framework/render"
	"github.com/mothfuzz/dyndraw/framework/transform"
	//"github.com/veandco/go-sdl2/sdl"
)

type TileSet struct {
	Image  string
	W, H   int
	TW, TH int
}

type Line struct {
	a Vec2
	b Vec2
}

type TileMap struct {
	TileSet
	Data   [][]uint8
	Planes []Line
	transform.Transform
	render.SpriteAnimation
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

func (t *TileMap) Init() {
	t.Transform = transform.Origin2D(t.TileSet.TW, t.TileSet.TH)
	t.SpriteAnimation = render.SpriteAnimation{
		Frames: [][]float32{},
		Tags:   map[string][]int{"tiles": {}},
	}
	totalWidth := float32(t.TileSet.W * t.TileSet.TW)
	totalHeight := float32(t.TileSet.H * t.TileSet.TH)
	for i := 0; i < t.TileSet.H; i++ {
		for j := 0; j < t.TileSet.W; j++ {
			t.SpriteAnimation.Frames = append(t.SpriteAnimation.Frames, []float32{
				float32(j*t.TileSet.TW) / totalWidth,
				float32(i*t.TileSet.TH) / totalHeight,
				float32(t.TileSet.TW) / totalWidth,
				float32(t.TileSet.TH) / totalHeight,
			})
			fmt.Println(t.SpriteAnimation.Frames[len(t.SpriteAnimation.Frames)-1])
			t.SpriteAnimation.Tags["tiles"] = append(t.SpriteAnimation.Tags["tiles"], i+j)
		}
	}
	if t.Planes == nil {
		for i := 0; i < len(t.Data); i++ {
			for j := 0; j < len(t.Data[i]); j++ {
				xOffset := float32(t.TileSet.TW) / 2.0
				yOffset := 400.0 + float32(t.TileSet.TH)/2.0 - float32(len(t.Data)*t.TileSet.TH)
				x := float32(j*t.TileSet.TW) + xOffset
				y := float32(i*t.TileSet.TH) + yOffset
				w := float32(t.TileSet.TW) / 2.0
				h := float32(t.TileSet.TH) / 2.0
				switch t.Data[i][j] {
				case 1:
					//t.Planes = append(t.Planes, Line{Vec2{x + w, y - h}, Vec2{x + w, y + h}})
					//t.Planes = append(t.Planes, Line{Vec2{x + w, y + h}, Vec2{x - w, y + h}})
					t.Planes = append(t.Planes, Line{Vec2{x - w, y + h}, Vec2{x + w, y - h}})
				case 2:
					t.Planes = append(t.Planes, Line{Vec2{x - w, y - h}, Vec2{x + w, y + h}})
					//t.Planes = append(t.Planes, Line{Vec2{x + w, y + h}, Vec2{x - w, y + h}})
					//t.Planes = append(t.Planes, Line{Vec2{x - w, y + h}, Vec2{x - w, y - h}})
				case 3:
					if !tileOccupied(t, j, i-1, []uint8{3, 2, 1}) {
						t.Planes = append(t.Planes, Line{Vec2{x - w, y - h}, Vec2{x + w, y - h}})
					}
					if !tileOccupied(t, j+1, i, []uint8{3, 2}) {
						t.Planes = append(t.Planes, Line{Vec2{x + w, y - h}, Vec2{x + w, y + h}})
					}
					if !tileOccupied(t, j, i+1, []uint8{3}) {
						t.Planes = append(t.Planes, Line{Vec2{x + w, y + h}, Vec2{x - w, y + h}})
					}
					if !tileOccupied(t, j-1, i, []uint8{3, 1}) {
						t.Planes = append(t.Planes, Line{Vec2{x - w, y + h}, Vec2{x - w, y - h}})
					}
				}
			}
		}
	}

}
func (t *TileMap) Update()  {}
func (t *TileMap) Destroy() {}

func (t *TileMap) Draw() {
	xOffset := float32(t.TileSet.TW) / 2.0
	yOffset := 400.0 + float32(t.TileSet.TH)/2.0 - float32(len(t.Data)*t.TileSet.TH)
	for i, row := range t.Data {
		for j, tile := range row {
			if tile != 0 {
				t.Transform.SetPosition2D(
					float32(j*t.TileSet.TW)+xOffset,
					float32(i*t.TileSet.TH)+yOffset,
				)
				render.DrawSpriteAnimated(t.TileSet.Image, t.Transform.Mat4(), t.SpriteAnimation.GetTexCoords("tiles", int(tile)))
			}
		}
	}
}

func distance(a Vec2, b Vec2) float32 {
	//return b.Sub(a).Len()
	return float32(math.Sqrt(
		math.Pow(float64(b.X()-a.X()), 2) +
			math.Pow(float64(b.Y()-a.Y()), 2)))
}

//distance from a point p to a line
func lineDistance(p Vec2, l Line) float32 {
	t := (l.b.X()-l.a.X())*(l.a.Y()-p.Y()) - (l.a.X()-p.X())*(l.b.Y()-l.a.Y())
	return float32(math.Abs(float64(t))) / distance(l.a, l.b)
}

func perpendicular(a Vec2) Vec2 {
	return Vec2{-a.Y(), a.X()}
}

func withinLine(c Vec2, l Line) bool {
	xMin := float32(math.Min(float64(l.a.X()), float64(l.b.X())))
	xMax := float32(math.Max(float64(l.a.X()), float64(l.b.X())))
	yMin := float32(math.Min(float64(l.a.Y()), float64(l.b.Y())))
	yMax := float32(math.Max(float64(l.a.Y()), float64(l.b.Y())))
	//fmt.Printf("%f should be between (%f, %f)\n", c.X(), xMin, xMax)
	//fmt.Printf("%f should be between (%f, %f)\n", c.Y(), yMin, yMax)
	if c.X() > xMax || c.X() < xMin {
		//fmt.Println("x out of range")
		return false
	}
	if c.Y() > yMax || c.Y() < yMin {
		//fmt.Println("y out of range")
		return false
	}
	return true
}
func pointInCircle(p Vec2, c Vec2, r float32) bool {
	d := p.Sub(c).LenSqr()
	return d < r*r
}

//moves a bounding sphere against arbitrary planes
//bounding sphere plane. not binary space partition
func MoveAgainstLines(t *transform.Transform, planes []Line, xspeed float32, yspeed float32, radius float32) (float32, float32) {
	velocity := Vec2{xspeed, yspeed}
	pos := t.GetPositionV().Vec2().Add(velocity)
	//pos = pos.Add(Vec2{radius, radius})
	for _, p := range planes {
		l := p.b.Sub(p.a)
		n := perpendicular(l).Normalize()
		//perpendicular (i.e. shortest) vector from position to plane
		proj := n.Mul(p.b.Sub(pos).Dot(n))
		//if proj.Len() < radius {
		if proj.LenSqr() < radius*radius {
			//project point to line and check if actually within bounds
			pproj := pos.Add(proj)
			if withinLine(pproj, p) ||
				pointInCircle(p.a, pos, radius) ||
				pointInCircle(p.b, pos, radius) {
				//fmt.Println("passed line test", sdl.GetTicks())
				adj := n.Mul(velocity.Dot(n)) //.Mul(2) //bouncy :3
				velocity = velocity.Sub(adj)
				pos = t.GetPositionV().Vec2().Add(velocity)
			}
		}
	}
	return velocity.X(), velocity.Y()
}

/*
//moves a bounding sphere against a series of walls
def move_bsp3(transform, radius, velocity, wall_actors):
    for wa in wall_actors:
        pos = transform().position + velocity
        ppos = wa.transform().position
        #get vector from point to plane
        dist = pos - ppos
        #project it onto normal (assumed to be normalized already)
        #this gives us a vector from the point perpendicular to the plane
        #the length of which is the shortest possible distance
        v = dot(dist, wa.normal) * wa.normal
        if length(v) < radius:
            #find the nearest point on the plane along that vector
            pp = pos + v
            #check if the point is actually within the bounds of the plane
            a = wa.points[0]
            b = wa.points[1]
            c = wa.points[2]
            d = wa.points[3]
            axis1 = a-b
            axis2 = a-d
            p1 = dot(axis1, pp)
            p2 = dot(axis2, pp)
            if p1 < dot(axis1, a) and p1 > dot(axis1, b) and p2 < dot(axis2, a) and p2 > dot(axis2, d):
                #if colliding with a wall, subtract velocity going in the wall's direction
                #to prevent movement
                adj = wa.normal * dot(velocity, wa.normal)
                velocity -= adj
    return velocity
*/
