package creature

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type Creature struct {
	Id          int
	Style       tcell.Style
	Name        rune
	Type        int
	X, Y        int
	Tx, Ty      int
	Max, Health int
	Damage      int
}

func (c *Creature) Draw(s tcell.Screen) {
	s.SetContent(c.X, c.Y, c.Name, nil, c.Style)
}

func (c *Creature) Move(x, y int) {
	c.X, c.Y = x, y
}

func (c *Creature) TurnRound(x, y int) (tx, ty int, canMove bool) {
	if c.Health <= 0 {
		return
	}
	tx, ty = c.X, c.Y
	xdis, ydis := math.Abs(float64(c.X-x)), math.Abs(float64(c.Y-y))
	if xdis > ydis {
		if c.X > x {
			tx--
		} else {
			tx++
		}
	} else {
		if c.Y > y {
			ty--
		} else {
			ty++
		}
	}
	c.Tx, c.Ty = tx, ty
	if math.Sqrt(math.Pow(float64(c.X-x), 2)+math.Pow(float64(c.Y-y), 2)) < 2 {
		return
	}
	if xdis < 3 && ydis == 0 {
		return
	} else if ydis < 3 && xdis == 0 {
		return
	}
	canMove = true
	return
}

func (c *Creature) Heal(v int) int {
	c.Health += v
	if c.Health < (c.Max / 20) {
		c.Style = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.Color196)
	}
	return c.Health
}

type Corpse struct {
	Id    int
	Style tcell.Style
	Name  rune
	Type  int
	X, Y  int
}

func (c *Corpse) Draw(s tcell.Screen) {
	s.SetContent(c.X, c.Y, c.Name, nil, c.Style)
}
