package main

import (
	"fmt"
	"math"

	"github.com/mothfuzz/dyndraw/framework/render"
	"github.com/mothfuzz/dyndraw/framework/transform"
	. "github.com/mothfuzz/dyndraw/framework/vecmath"
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

//actually a triangle
type Plane struct {
	origin Vec3
	normal Vec3
	points [3]Vec3
}

func (p *Plane) Normal() Vec3 {
	return p.normal
}

func triNorm(a, b, c Vec3) Vec3 {
	//(B - A) x (C - A)
	return b.Sub(a).Cross(c.Sub(a)).Normalize()
}
func NewPlane(x, y, z float32, a, b, c Vec3) Plane {
	t := Vec3{x, y, z}
	n := triNorm(a, b, c)
	return Plane{t, n, [3]Vec3{a, b, c}}
}

//-
func newXPlane(x, y, l float32) Plane {
	a := Vec3{x, y, -l}
	b := Vec3{x + l/2, y, 0}
	c := Vec3{x - l/2, y, 0}
	return NewPlane(x, y, 0, a, b, c)
}

//|
func newYPlane(x, y, l float32) Plane {
	a := Vec3{x, y, -l}
	b := Vec3{x, y - l/2, 0}
	c := Vec3{x, y + l/2, 0}
	return NewPlane(x, y, 0, a, b, c)
}

type TileMap struct {
	TileSet
	Data   [][]uint8
	Planes []Plane
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

func (t *TileMap) offsets() (float32, float32) {
	xOffset := float32(t.TileSet.TW) / 2.0
	yOffset := 400.0 + float32(t.TileSet.TH)/2.0 - float32(len(t.Data)*t.TileSet.TH)
	return xOffset, yOffset
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
					t.Planes = append(t.Planes, NewPlane(x, y, 0, Vec3{x, y, -w}, Vec3{x + w/2, y - h/2, 0}, Vec3{x - w/2, y + h/2, 0}))
				case 2:
					t.Planes = append(t.Planes, NewPlane(x, y, 0, Vec3{x, y, -w}, Vec3{x + w/2, y + h/2, 0}, Vec3{x - w/2, y - h/2, 0}))
				case 3:
					if !tileOccupied(t, j, i-1, []uint8{1, 2, 3}) {
						t.Planes = append(t.Planes, newXPlane(x, y-h/2, h))
					}
					if !tileOccupied(t, j, i+1, []uint8{1, 2, 3}) {
						t.Planes = append(t.Planes, newXPlane(x, y+h/2, h))
					}
					if !tileOccupied(t, j-1, i, []uint8{1, 3}) {
						t.Planes = append(t.Planes, newYPlane(x-w/2, y, w))
					}
					if !tileOccupied(t, j+1, i, []uint8{2, 3}) {
						t.Planes = append(t.Planes, newYPlane(x+w/2, y, w))
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
				t.Transform.SetPosition2D(
					float32(j*t.TileSet.TW)+xOffset,
					float32(i*t.TileSet.TH)+yOffset,
				)
				render.DrawSpriteAnimated(t.TileSet.Image, t.Transform.Mat4(), t.SpriteAnimation.GetTexCoords("tiles", int(tile)))
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

//distance from a point p to a line
func lineDistance(p Vec2, l Line) float32 {
	t := (l.b.X()-l.a.X())*(l.a.Y()-p.Y()) - (l.a.X()-p.X())*(l.b.Y()-l.a.Y())
	return float32(math.Abs(float64(t))) / l.b.Sub(l.a).Len()
}

func perpendicular(a Vec2) Vec2 {
	return Vec2{-a.Y(), a.X()}
}

func withinLine(c Vec2, l Line) bool {
	xMin := float32(math.Min(float64(l.a.X()), float64(l.b.X())))
	xMax := float32(math.Max(float64(l.a.X()), float64(l.b.X())))
	yMin := float32(math.Min(float64(l.a.Y()), float64(l.b.Y())))
	yMax := float32(math.Max(float64(l.a.Y()), float64(l.b.Y())))
	if c.X() > xMax || c.X() < xMin {
		return false
	}
	if c.Y() > yMax || c.Y() < yMin {
		return false
	}
	return true
}
func pointInCircle(p Vec2, c Vec2, r float32) bool {
	d := p.Sub(c).LenSqr()
	return d < r*r
}

//moves a bounding sphere against arbitrary planes
func MoveAgainstLines(t *transform.Transform, planes []Line, xspeed float32, yspeed float32, radius float32) (float32, float32) {
	velocity := Vec2{xspeed, yspeed}
	pos := t.GetPositionV().Vec2().Add(velocity)
	for _, p := range planes {
		l := p.b.Sub(p.a)
		n := perpendicular(l).Normalize()
		//perpendicular (i.e. shortest) vector from position to plane
		proj := n.Mul(p.b.Sub(pos).Dot(n))
		if proj.LenSqr() < radius*radius {
			//project point to line and check if actually within bounds
			pproj := pos.Add(proj)
			if withinLine(pproj, p) ||
				pointInCircle(p.a, pos, radius) ||
				pointInCircle(p.b, pos, radius) {
				adj := n.Mul(velocity.Dot(n)) //.Mul(2) //bouncy :3
				velocity = velocity.Sub(adj)
				pos = t.GetPositionV().Vec2().Add(velocity)
			}
		}
	}
	t.Translate2D(velocity.X(), velocity.Y())
	return velocity.X(), velocity.Y()
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

func insideTriangleVertices(p Vec3, r float32, a, b, c Vec3) bool {
	r2 := r * r
	if p.Sub(a).LenSqr() <= r2 {
		return true
	}
	if p.Sub(b).LenSqr() <= r2 {
		return true
	}
	if p.Sub(c).LenSqr() <= r2 {
		return true
	}
	return false
}
func sphereEdge(p Vec3, r float32, a, b Vec3) bool {
	r2 := r * r
	//check a
	/*if p.Sub(a).LenSqr() <= r2 {
		return true
	}
	//check b
	if p.Sub(b).LenSqr() <= r2 {
		return true
	}*/
	//check parametric distance
	ab := b.Sub(a)
	t := p.Sub(a).Dot(ab.Normalize())
	if t > 0 && t < 1 {
		x := a.Add(ab.Mul(t))
		if p.Sub(x).LenSqr() <= r2 {
			return true
		}
	}
	return false
}
func insideTriangleEdges(p Vec3, r float32, a, b, c Vec3) bool {
	if sphereEdge(p, r, a, b) {
		return true
	}
	if sphereEdge(p, r, b, c) {
		return true
	}
	if sphereEdge(p, r, c, a) {
		return true
	}
	return false
}
func pointInTriangle(p Vec3, a, b, c Vec3) bool {
	ab := b.Sub(a)
	ac := c.Sub(a)
	ap := p.Sub(a)
	abac := ab.Cross(ac)
	//check if 3 points are coplanar first
	//the floating point errors are strong with this one...
	if math.Abs(float64(ap.Dot(abac))) <= 1e-2 {
		//compute barycentric coords
		//u = ||CAxCP|| / ||ABxAC||
		//v = ||ABxAP|| / ||ABxAC||
		//w = ||BCxBP|| / ||ABxAC||
		abacl := abac.LenSqr()
		u := c.Sub(a).Cross(c.Sub(p)).LenSqr() / abacl
		v := a.Sub(b).Cross(a.Sub(p)).LenSqr() / abacl
		w := b.Sub(c).Cross(b.Sub(p)).LenSqr() / abacl
		if u >= 0 && u <= 1 && v >= 0 && v <= 1 && w >= 0 && w <= 1 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//moves a bounding sphere against a series of walls
func MoveAgainstPlanes(t *transform.Transform, planes []Plane, radius float32, xspeed, yspeed, zspeed float32) (float32, float32, float32) {
	velocity := Vec3{xspeed, yspeed, zspeed}
	for _, p := range planes {
		pos := t.GetPositionV().Add(velocity)
		//get vector from point to plane
		dist := pos.Sub(p.origin)
		//project it onto normal (assumed to be normalized already)
		//this gives us a vector from the point perpendicular to the plane
		//the length of which is the shortest possible distance
		v := p.normal.Mul(dist.Dot(p.normal))
		if v.LenSqr() <= radius*radius {
			a := p.points[0]
			b := p.points[1]
			c := p.points[2]
			//find the nearest point on the plane along that vector
			//then check if the point is actually within the bounds of the triangle
			if pointInTriangle(pos.Add(v), a, b, c) ||
				insideTriangleVertices(pos, radius, a, b, c) ||
				insideTriangleEdges(pos, radius, a, b, c) {
				//if colliding with a wall, subtract velocity going in the wall's direction
				//to prevent movement
				adj := p.normal.Mul(velocity.Dot(p.normal)) //.Mul(2) //bouncy :3
				preserve := velocity.LenSqr()
				velocity = velocity.Sub(adj)
				mag := velocity.LenSqr()
				//attempt to preserve momentum against slopes etc.
				//not physically accurate but it's more fun.
				if mag > 0 {
					velocity = velocity.Mul(preserve / mag / 1.0)
				}
			}
		}
	}
	return velocity.Elem()
}

type RayHit struct {
	Plane
	I Vec3
}

func RayCast(pos Vec3, planes []Plane, ray Vec3) []RayHit {
	hits := []RayHit{}
	for _, p := range planes {
		rdot1 := ray.Dot(p.normal)
		rdot2 := p.origin.Sub(pos).Dot(p.normal)
		t := rdot2 / rdot1
		i := pos.Add(ray.Mul(t))
		if t >= 0 && //t <= 1 && //t > 1 if plane exceeds distance
			pointInTriangle(i, p.points[0], p.points[1], p.points[2]) {
			hits = append(hits, RayHit{p, i})
		}
	}
	//fmt.Println("len(hits)", len(hits))
	return hits
}
func RayCastLen(pos Vec3, planes []Plane, ray Vec3, l float32) (RayHit, bool) {
	ll := l * l
	shortest := ll
	ok, hit := false, RayHit{}
	for _, p := range RayCast(pos, planes, ray) {
		dist := p.I.Sub(pos).LenSqr()
		if dist <= shortest {
			shortest = dist
			ok = true
			hit = p
		}
	}
	return hit, ok
}
