package main

import (
	"fmt"
	"math"

	"github.com/mothfuzz/dyndraw/framework/input"
	"github.com/mothfuzz/dyndraw/framework/render"
	"github.com/mothfuzz/dyndraw/framework/transform"
	. "github.com/mothfuzz/dyndraw/framework/vecmath"
)

type Player struct {
	transform.Transform
	hp               int8
	xspeed           float32
	xspeedMax        float32
	xfriction        float32
	yspeed           float32
	gravity          float32
	grounded         bool
	groundMultiplier float32
}

var CurrentLevel *TileMap = nil

const pw = 16
const ph = 16

func (p *Player) Init() {
	p.Transform = transform.Origin2D(pw, ph)
	p.Transform.SetPosition(640/2, 0, -0.1)
	p.hp = 10
	p.gravity = 0.1
	p.grounded = false
	p.groundMultiplier = 1
	p.xspeedMax = 8
	p.xfriction = 0.8
}
func (p *Player) Update() {

	if input.IsKeyDown("left") {
		p.xspeed -= 0.25
		p.SetScale2D(-pw, ph)
	}
	if input.IsKeyDown("right") {
		p.xspeed += 0.25
		p.SetScale2D(pw, ph)
	}
	if input.IsKeyPressed("up") && p.grounded {
		p.grounded = false
		p.yspeed = -4
	}
	//if either in the air or on a flat surface, apply gravity
	if !p.grounded || p.groundMultiplier == 1 {
		p.yspeed += p.gravity
	}
	p.xspeed *= p.xfriction
	p.xspeed *= p.groundMultiplier

	if p.xspeed < -p.xspeedMax {
		p.xspeed = -p.xspeedMax
	}
	if p.xspeed > p.xspeedMax {
		p.xspeed = p.xspeedMax
	}
	if math.Abs(float64(p.xspeed)) < 0.1 {
		p.xspeed = 0
	}
	if math.Abs(float64(p.yspeed)) < 0.1 {
		p.yspeed = 0
	}
	if CurrentLevel != nil {

		//apply collisions per-axis to avoid getting 'stuck' at 'seams'
		xadj, yadj := float32(0), float32(0)
		xadj, yadj, _ = MoveAgainstPlanes(&p.Transform, CurrentLevel.Planes, pw/2-0.5, p.xspeed, 0, 0)
		p.xspeed = xadj
		p.Translate2D(p.xspeed, yadj)
		xadj, yadj, _ = MoveAgainstPlanes(&p.Transform, CurrentLevel.Planes, pw/2-0.5, 0, p.yspeed, 0)
		p.yspeed = yadj
		p.Translate2D(xadj, p.yspeed)

		//keep your feet on the ground
		feet := p.GetPositionV().Add(Vec3{0, ph / 2, 0})
		direction := float32(0)
		if p.xspeed > 0 {
			direction = 1
		}
		if p.xspeed < 0 {
			direction = -1
		}
		if hit, ok := RayCastLen(feet.Add(Vec3{pw / 4 * direction, 0, 0}), CurrentLevel.Planes, Vec3{0, 1, 0}, 4); ok {
			p.grounded = true
			dot := hit.Plane.Normal().Dot(Vec3{0, 1, 0})
			p.groundMultiplier = 1.0 / (dot * dot)
			//p.SetPosition2D(p.X(), hit.I.Y()-ph/2) //too forceful, messes with velocity
		} else {
			p.grounded = false
			p.groundMultiplier = 1
		}

	}
	if p.Y()+8 > 400 {
		p.Translate2D(0, 400-(p.Y()+8))
		p.grounded = true
		p.yspeed = 0
	}
}
func (p *Player) Destroy() {
	fmt.Println("game over bro!!")
}
func (p *Player) Draw() {
	render.DrawSprite("player.png", p.Transform.Mat4())
}
