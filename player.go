package main

import (
	"fmt"
	"math"

	"github.com/mothfuzz/dyndraw/framework/input"
	"github.com/mothfuzz/dyndraw/framework/render"
	"github.com/mothfuzz/dyndraw/framework/transform"
)

type Player struct {
	transform.Transform
	hp      int8
	xspeed  float32
	yspeed  float32
	gravity float32
}

var CurrentLevel *TileMap = nil

const pw = 16
const ph = 16

func (p *Player) Init() {
	p.Transform = transform.Origin2D(pw, ph)
	p.Transform.SetPosition(640/2, 0, 0)
	p.hp = 10
	p.gravity = 0.1
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
	if input.IsKeyPressed("up") {
		p.yspeed = -4
	}
	p.yspeed += p.gravity

	/*if input.IsKeyDown("left") {
		p.xspeed = -1
	}
	if input.IsKeyDown("right") {
		p.xspeed = 1
	}
	if input.IsKeyDown("up") {
		p.yspeed = -1
	}
	if input.IsKeyDown("down") {
		p.yspeed = 1
	}*/

	p.xspeed *= 0.8
	//p.yspeed *= 0.8
	if math.Abs(float64(p.xspeed)) < 0.1 {
		p.xspeed = 0
	}
	if math.Abs(float64(p.yspeed)) < 0.1 {
		p.yspeed = 0
	}
	if CurrentLevel != nil {
		//TODO: add raycasting
		// 0 g on ground
		// and also movement speed adjusted based on normal of current floor triangle

		//p.xspeed, p.yspeed = MoveAgainstLines(&p.Transform, CurrentLevel.Planes, p.xspeed, p.yspeed, pw/2-0.5)

		//p.xspeed, p.yspeed = MoveAgainstTiles(&p.Transform, CurrentLevel, p.xspeed, p.yspeed, pw-0.5, ph)

		//p.xspeed, p.yspeed, _ = MoveAgainstPlanes(&p.Transform, CurrentLevel.Planes, pw/2-0.5, p.xspeed, p.yspeed, 0)
		//p.Translate2D(p.xspeed, p.yspeed)

		//apply collisions per-axis to avoid getting 'stuck' at 'seams'
		xadj, yadj := float32(0), float32(0)
		xadj, yadj, _ = MoveAgainstPlanes(&p.Transform, CurrentLevel.Planes, pw/2-0.5, p.xspeed, 0, 0)
		p.xspeed = xadj
		p.Translate2D(p.xspeed, yadj)
		xadj, yadj, _ = MoveAgainstPlanes(&p.Transform, CurrentLevel.Planes, pw/2-0.5, 0, p.yspeed, 0)
		if yadj != p.yspeed {
			p.gravity = 0
		} else {
			p.gravity = 0.1
		}
		p.yspeed = yadj
		p.Translate2D(xadj, p.yspeed)
	}
	if p.GetPositionV().Y()+8 > 400 {
		p.Translate2D(0, 400-(p.GetPositionV().Y()+8))
		p.yspeed = 0
	}
}
func (p *Player) Destroy() {
	fmt.Println("game over bro!!")
}
func (p *Player) Draw() {
	render.DrawSprite("player.png", p.Transform.Mat4())
}
