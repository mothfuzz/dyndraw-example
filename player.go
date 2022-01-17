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

	items []Item
	actors.Mailbox
}

var CurrentLevel []collision.Collider = nil

const pw = 16
const ph = 16

func (p *Player) Init() {
	p.Transform = transform.Origin2D(pw, ph)
	p.Transform.SetPosition(640/2, 0, 0)
	p.state = ground
	p.hp = 10
	p.gravity = 0.1
	p.xspeedMax = 8
	p.yspeedMax = 6
	p.xfriction = 0.8

	p.Collider = collision.NewBoundingSphere(0.5)
	p.Collider.IgnoreRaycast = true

	p.items = []Item{}
	p.Mailbox = actors.Listen(p, Item{}) //listen for items
}

func (p *Player) ProcessInput() {
	if input.IsKeyDown("left") {
		p.xspeed -= 0.25
		p.SetScale2D(-pw, ph)
	}
	if input.IsKeyDown("right") {
		p.xspeed += 0.25
		p.SetScale2D(pw, ph)
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
	xadj, yadj, _ := collision.MoveAgainstPlanes(&p.Transform, CurrentLevel[0].Planes, pw/2-0.5, p.xspeed, 0, 0)
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
	if p.state == ground {
		highestY := p.Y()
		if leftOk && rightOk {
			highestY = float32(math.Min(float64(leftHit.I.Y()), float64(rightHit.I.Y()))) - ph/2
		}
		if leftOk && !rightOk {
			highestY = leftHit.I.Y() - ph/2
		}
		if rightOk && !leftOk {
			highestY = rightHit.I.Y() - ph/2
		}
		if p.Y()+ph/2 >= highestY {
			p.SetPosition2D(p.X(), highestY)
		}
	}
	if (leftOk || rightOk) && p.state != jumping {
		p.yspeed = 0
		p.state = ground
	} else {
		if p.state == ground {
			p.state = falling
		}
	}

	//avoid planes
	xadj, yadj, _ := collision.MoveAgainstPlanes(&p.Transform, CurrentLevel[0].Planes, pw/2-0.5, 0, p.yspeed, 0)
	p.yspeed = yadj
	p.Translate2D(xadj, p.yspeed)
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
		//p.SetPosition2D(float32(mx), float32(my))
		render.ActiveCamera.Look2D(Vec2{p.X() + float32(mx) - 640/2, p.Y() + float32(my) - 400/2})
	} else {
		render.ActiveCamera.Look2D(Vec2{p.X(), p.Y()})
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
			case Item:
				fmt.Printf("Player grabbed a %s! \"%s\"\n", m.Name, m.Description)
				p.items = append(p.items, m)
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
}
