package game

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func Input(g *Game) {
	ox, oy := -1, -1
	for {
		// Update screen
		g.Screen.Show()

		// Poll event
		ev := g.Screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			g.Screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlL:
				g.Screen.Sync()
			case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight,
				tcell.KeyUpLeft, tcell.KeyUpRight, tcell.KeyDownLeft, tcell.KeyDownRight:
				g.Turn(ev.Key())
			}
			g.Graph()
		case *tcell.EventMouse:
			x, y := ev.Position()
			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2:
				if ox < 0 {
					ox, oy = x, y // record location when click started
				}
			case tcell.ButtonNone:
				if ox >= 0 {
					label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
					g.DrawBox(ox, oy, x, y, label)
					ox, oy = -1, -1
				}
			}
		}
	}
}
