package main

import (
	"fmt"
	"math"

	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/collision"
	"github.com/mothfuzz/letsgo/input"
	"github.com/mothfuzz/letsgo/render"
	"github.com/mothfuzz/letsgo/transform"
	. "github.com/mothfuzz/letsgo/vecmath"
)

type Inventory struct {
	render.SpriteAnimation
	Visible  bool
	Capacity int
	Items    []Item
}

func (inv *Inventory) Init() {
	inv.Visible = false
	//inv.Items = make([]Item, inv.Capacity)
	inv.SpriteAnimation = render.SpriteAnimation{
		Frames: []render.Frame{
			{X: 0.0, Y: 0.0, W: 0.5, H: 1.0},
			{X: 0.5, Y: 0.0, W: 0.5, H: 1.0},
		},
		Tags: map[string][]int{
			"border": {0},
			"center": {1},
		},
	}
}
func (*Inventory) Update()  {}
func (*Inventory) Destroy() {}
func (inv *Inventory) Draw() {
	if inv.Visible {
		//draw inventory backdrop
		t := transform.Origin2D()
		t.SetPosition(render.RelativeToCamera(16, 16).Elem())
		t.Translate(0, 0, -0.5)
		render.DrawSpriteAnimated("inventory.png", t.Mat4(), inv.SpriteAnimation.GetTexCoords("border", 0))
		for i := 1; i < inv.Capacity-1; i++ {
			t.Translate2D(32, 0)
			render.DrawSpriteAnimated("inventory.png", t.Mat4(), inv.SpriteAnimation.GetTexCoords("center", 0))
		}
		t.Translate2D(32, 0)
		t.SetScale2D(-1, 1)
		render.DrawSpriteAnimated("inventory.png", t.Mat4(), inv.SpriteAnimation.GetTexCoords("border", 0))
		t.SetScale2D(1, 1)

		//draw items
		t = transform.Origin2D()
		t.SetPosition(render.RelativeToCamera(16, 16).Elem())
		t.Translate(0, 0, -1)
		for i := 0; i < len(inv.Items); i++ {
			render.DrawSprite(inv.Items[i].Icon, t.Mat4())
			t.Translate2D(32, 0)
		}
	}
}
func (inv *Inventory) AddItem(i *Item) {
	if len(inv.Items) < inv.Capacity {
		//remove from world, put in inventory
		actors.Destroy(i)
		inv.Items = append(inv.Items, *i)
		fmt.Printf("grabbed a %s! \"%s\"\n", i.Name, i.Description)
	}
}
func (inv *Inventory) Show() {
	inv.Visible = true
}
func (inv *Inventory) Hide() {
	inv.Visible = false
}
func (inv *Inventory) Toggle() {
	inv.Visible = !inv.Visible
}

type PlayerState int

const (
	ground PlayerState = iota
	jumping
	falling
)

func (ps PlayerState) String() string {
	return []string{"Ground", "Jumping", "Falling"}[ps]
}

type Player struct {
	transform.Transform
	hp               int8
	state            PlayerState
	xspeed           float32
	xspeedMax        float32
	yspeedMax        float32
	xfriction        float32
	yspeed           float32
	gravity          float32
	grounded         bool
	groundMultiplier float32

	collision.Collider

	Inventory
	actors.Mailbox
}

var CurrentLevel *TileMap = nil

const pw = 16
const ph = 16

func (p *Player) Init() {
	p.Transform = transform.Origin2D()
	p.Transform.SetPosition(640/2-128*2, 0, -0.1)
	p.state = ground
	p.hp = 10
	p.gravity = 0.1
	p.xspeedMax = 8
	p.yspeedMax = 6
	p.xfriction = 0.8

	p.Collider = collision.NewBoundingSphere(8)
	p.Collider.IgnoreRaycast = true

	p.Inventory.Capacity = 4
	actors.Spawn(&p.Inventory)
	p.Mailbox = actors.Listen(p, &Item{}) //listen for items
}

func (p *Player) ProcessInput() {
	if input.IsKeyDown("left") {
		p.xspeed -= 0.25
		p.SetScale2D(-1, 1)
	}
	if input.IsKeyDown("right") {
		p.xspeed += 0.25
		p.SetScale2D(1, 1)
	}
	if input.IsKeyPressed("up") && p.state == ground {
		p.state = jumping
		p.yspeed = -4
	}
}

func (p *Player) MoveX() {
	p.xspeed *= p.xfriction
	if p.xspeed < -p.xspeedMax {
		p.xspeed = -p.xspeedMax
	}
	if p.xspeed > p.xspeedMax {
		p.xspeed = p.xspeedMax
	}
	if math.Abs(float64(p.xspeed)) < 0.1 {
		p.xspeed = 0
	}

	//initial movement, avoiding walls
	xadj, yadj, _ := collision.MoveAgainstPlanes(&p.Transform, CurrentLevel.Planes, pw/2-0.5, p.xspeed, 0, 0)
	p.xspeed = xadj
	p.Translate2D(p.xspeed, yadj)
}

func (p *Player) MoveY() {

	//fmt.Println(p.state)
	//fmt.Println(p.xspeed)
	//fmt.Println(p.yspeed)
	//fmt.Printf("%2.2f\t%2.2f\t%2.2f\t%2.2f\t%2.2f\t%v\n", p.X(), p.Y(), p.Z(), p.xspeed, p.yspeed, p.state)

	//if in the air, apply gravity
	if p.state == jumping || p.state == falling {
		p.yspeed += p.gravity
		if p.yspeed > 0 {
			p.state = falling
		}
	}
	//min/max velocity
	if p.yspeed < -p.yspeedMax {
		p.yspeed = -p.yspeedMax
	}
	if p.yspeed > p.yspeedMax {
		p.yspeed = p.yspeedMax
	}
	if math.Abs(float64(p.yspeed)) < 0.1 {
		p.yspeed = 0
	}

	//keep your feet on the ground
	feet := p.GetPositionV() //.Add(Vec3{0, ph / 2, 0})
	leftFoot := feet.Add(Vec3{-pw / 3.0, 0, 0})
	rightFoot := feet.Add(Vec3{pw / 3.0, 0, 0})
	leftHit, leftOk := collision.RayCastLen(leftFoot, Vec3{0, 1, 0}, ph)
	rightHit, rightOk := collision.RayCastLen(rightFoot, Vec3{0, 1, 0}, ph)
	if (leftOk || rightOk) && p.state != jumping {
		p.yspeed = 0
		p.state = ground
		highestY := p.Y()
		if leftOk && rightOk {
			highestY = float32(math.Min(float64(leftHit.Point.Y()), float64(rightHit.Point.Y()))) - ph/2
		}
		if leftOk && !rightOk {
			highestY = leftHit.Point.Y() - ph/2
		}
		if rightOk && !leftOk {
			highestY = rightHit.Point.Y() - ph/2
		}
		if p.Y()+ph/2 >= highestY {
			p.SetPosition2D(p.X(), highestY)
		}
	} else {
		if p.state == ground {
			p.state = falling
		}
	}

	//avoid planes
	_, p.yspeed, _ = collision.MoveAgainstPlanes(&p.Transform, CurrentLevel.Planes, pw/2-0.5, 0, p.yspeed, 0)
	p.Translate2D(0, p.yspeed)
}

func (p *Player) Update() {

	p.ProcessInput()
	//apply collisions per-axis to avoid getting 'stuck' at 'seams'
	p.MoveX()
	p.MoveY()
	//fmt.Println(p.state)
	//fmt.Println(p.xspeed)
	//fmt.Println(p.yspeed)

	if input.IsKeyDown("left ctrl") {
		mx, my := input.GetMousePosition()
		//v := render.RelativeToCamera(mx, my)
		//p.SetPosition2D(v.X(), v.Y())
		render.ActiveCamera.Look2D(Vec2{p.X() + float32(mx) - 640/2, p.Y() + float32(my) - 400/2})
	} else {
		render.ActiveCamera.Look2D(Vec2{p.X(), p.Y()})
	}

	if input.IsKeyPressed("i") {
		p.Inventory.Toggle()
	}

	if p.Y()+ph/2 >= 400 {
		p.Translate2D(0, 400-(p.Y()+ph/2))
		p.yspeed = 0
		p.state = ground
	}
	for {
		select {
		case m := <-p.Mailbox:
			switch m := m.(type) {
			case *Item:
				p.Inventory.AddItem(m)
			}
		default:
			return
		}
	}
}
func (p *Player) Destroy() {
	fmt.Println("game over bro!!")
}
func (p *Player) Draw() {
	render.DrawSprite("player.png", p.Transform.Mat4())
	/*for i := range CurrentLevel.Planes {
		p := CurrentLevel.Planes[i]
		a := p.Points()[0]
		b := p.Points()[1]
		c := p.Points()[2]
		t := transform.Origin2D(4, 4)
		t.SetPosition2D(a.X(), a.Y())
		render.DrawSprite("point.png", t.Mat4())
		t.SetPosition2D(b.X(), b.Y())
		render.DrawSprite("point.png", t.Mat4())
		t.SetPosition2D(c.X(), c.Y())
		render.DrawSprite("point.png", t.Mat4())
	}*/
}
