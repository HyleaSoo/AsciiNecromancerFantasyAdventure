package game

import (
	"math/rand"
	"necromancer/creature"
	"necromancer/global"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Game struct {
	//maximum int
	Screen tcell.Screen
	Style  tcell.Style
	At     creature.Character
	Pet    creature.Character
	House  House
	Msts   []*creature.Monster
	Corps  []*creature.Corpse
	Scores int
}

func InitGameObject(x, y int) (g Game) {
	msts := []*creature.Monster{
		creature.NewMonster(global.MstStyle[global.MstZombie], x-x, y-y, global.AsciiZombie, global.MstZombie),
		creature.NewMonster(global.MstStyle[global.MstZombie], x-3, y-4, global.AsciiZombie, global.MstZombie),
	}
	return Game{
		At:    creature.NewCharacter(global.CptStyle[global.CptHero], x/2, y/2, global.AsciiHero, global.CptHero),
		Pet:   creature.NewCharacter(global.CptStyle[global.CptPet], x/2+1, y/2+1, global.AsciiPet, global.CptPet),
		House: NewHouse(10, 10, 20, 20, []DoorPosition{{10, 15}, {20, 15}}),
		Msts:  msts,
	}
}

func NewGame() *Game {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	game := InitGameObject(s.Size())
	game.Screen = s
	game.Style = defStyle
	game.Scores = 0
	return &game
}

func (g *Game) Turn(act tcell.Key) {
	tx, ty := g.At.X, g.At.Y
	switch act {
	case tcell.KeyUp:
		ty--
	case tcell.KeyDown:
		ty++
	case tcell.KeyLeft:
		tx--
	case tcell.KeyRight:
		tx++
	/* case tcell.KeyUpLeft:
		tx--
		ty--
	case tcell.KeyUpRight:
		tx++
		ty--
	case tcell.KeyDownLeft:
		tx--
		ty++
	case tcell.KeyDownRight:
		tx++
		ty++ */
	default:
		if act == 'a' {
			go g.UseSkill(SKID_CuiDuBiShou)
		} else if act == 's' {
			go g.UseSkill(SKID_ShiBao)
		}
	}

	g.turnSched(tx, ty)
}

func (g *Game) UseSkill(sk uint8) {
	SkillMap[sk](g)
}

func (g *Game) turnSched(tx, ty int) {
	if g.At.TurnRound(tx, ty); CanMove(g.Screen, tx, ty) {
		g.At.Move(tx, ty)
	}
	if ptx, pty, canMove := g.Pet.Creature.TurnRound(tx, ty); canMove {
		g.Pet.Follow(ptx, pty)
	}
	for k, v := range g.Msts {
		if mtx, mty, canMove := v.TurnRound(g.At.Tx, g.At.Ty); canMove && CanMove(g.Screen, mtx, mty) {
			v.Move(mtx, mty)
		}
		if v.X == g.At.Tx && v.Y == g.At.Ty {
			if hp := v.Heal(-g.At.Damage); hp <= 0 {
				g.MstDeath(k)
				continue
			}
		}
		if v.Tx == g.At.X && v.Ty == g.At.Y {
			if hp := g.At.Heal(-v.Damage); hp <= 0 {
				g.Screen.Clear()
				g.DrawText(tx, ty, "Necromancer shuld know what is death")
				g.Screen.Show()
				time.Sleep(5 * time.Second)
			}
		}
	}
	if len(g.Msts) < 2 {
		g.GenerateMst()
	}
}

/*
s(0x, ny)
s(mx, ny)
s(0y, nx)
s(my, nx)
*/
func (g *Game) GenerateMst() {
	mx, my := g.Screen.Size()
	mx, my = mx-1, my-4
	x, y := 0, 0
	if rand.Intn(4) < 2 {
		x = []int{0, mx}[rand.Intn(2)]
		y = rand.Intn(my)
	} else {
		x = rand.Intn(mx)
		y = []int{0, my}[rand.Intn(2)]
	}
	g.Msts = append(g.Msts, creature.NewMonster(global.MstStyle[global.MstZombie], x, y, global.AsciiZombie, global.MstZombie))
}

func (g *Game) MstDeath(i int) {
	if len(g.Msts) < (i+1) || g.Msts == nil {
		return
	}
	g.Corps = append(g.Corps, &creature.Corpse{Id: 10001, Style: global.CorpseStyle, Name: global.AsciiCorpse, Type: global.MstZombie, X: g.Msts[i].X, Y: g.Msts[i].Y})
	g.Msts = append(g.Msts[:i], g.Msts[i+1:]...)
	g.Scores += 1
}

func (g *Game) Graph() {
	g.Screen.Clear()
	g.StateBox()

	// g.House.Draw(g.Screen)
	for _, v := range g.Corps {
		v.Draw(g.Screen)
	}
	g.Pet.Draw(g.Screen)
	g.At.Draw(g.Screen)
	for _, v := range g.Msts {
		v.Draw(g.Screen)
	}
}

func (g *Game) StateBox() {
	sx, sy := g.Screen.Size()
	for i := 0; i < sx; i++ {
		g.Screen.SetContent(i, sy-3, tcell.RuneHLine, nil, g.Style)
	}
	g.DrawText(0, sy-2, "AT:"+strconv.Itoa(g.At.Health))
	g.DrawText(0, sy-1, "M:"+strconv.Itoa(len(g.Msts))+" C:"+strconv.Itoa(len(g.Corps))+" S:"+strconv.Itoa(g.Scores))
	for k, v := range g.Msts {
		g.DrawText(60, sy-(1+k), strconv.Itoa(k)+":"+strconv.Itoa(v.X)+","+strconv.Itoa(v.Y))
	}
}

func (g *Game) DrawBox(x1, y1, x2, y2 int, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	tstyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGray)
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			g.Screen.SetContent(col, row, '.', nil, tstyle)
		}
	}

	for col := x1; col <= x2; col++ { //tcell.RuneHLine
		g.Screen.SetContent(col, y1, '-', nil, g.Style)
		g.Screen.SetContent(col, y2, '-', nil, g.Style)
	}
	for row := y1 + 1; row < y2; row++ { //tcell.RuneVLine
		g.Screen.SetContent(x1, row, '|', nil, g.Style)
		g.Screen.SetContent(x2, row, '|', nil, g.Style)
	}

	g.DrawBoxText(x1+1, y1+1, x2-1, y2-1, text)
}

func (g *Game) DrawBoxText(x1, y1, x2, y2 int, text string) {
	row, col := y1, x1
	for _, r := range text {
		g.Screen.SetContent(col, row, r, nil, g.Style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func (g *Game) DrawText(x, y int, text string) {
	for i := 0; i < len(text); i++ {
		g.Screen.SetContent(x+i, y, rune(text[i]), nil, g.Style)
	}
}

func CanMove(s tcell.Screen, x, y int) bool {
	//超出屏幕
	sx, sy := s.Size()
	if x < 0 || y < 0 || x >= sx || y >= sy {
		return false
	}
	//碰撞
	if mainc, _, _, _ := s.GetContent(x, y); !IsPassable(mainc) {
		return false
	}
	return true
}

func IsPassable(p rune) bool {
	_, ok := map[rune]struct{}{
		global.AsciiHorizon:  {},
		global.AsciiDoor:     {},
		global.AsciiFloor:    {},
		global.AsciiFloorLow: {},
		global.AsciiCorpse:   {},
		global.AsciiPet:      {},
	}[p]
	return ok
}

func IsObstacle(p rune) bool {
	_, ok := map[rune]struct{}{
		global.AsciiHWall: {},
		global.AsciiVWall: {},
	}[p]
	return ok
}
