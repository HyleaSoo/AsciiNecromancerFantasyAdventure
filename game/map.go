package game

import (
	"necromancer/global"

	"github.com/gdamore/tcell/v2"
)

type House struct {
	style          tcell.Style
	x1, y1, x2, y2 int
	door           []DoorPosition
}
type DoorPosition struct {
	x, y int
}

func NewHouse(x1, y1, x2, y2 int, dp []DoorPosition) House {
	return House{
		global.MapStyle[global.MapHouse],
		x1, y1, x2, y2, dp,
	}
}

func (h *House) Draw(s tcell.Screen) {
	if h.y2 < h.y1 {
		h.y1, h.y2 = h.y2, h.y1
	}
	if h.x2 < h.x1 {
		h.x1, h.x2 = h.x2, h.x1
	}
	//地板
	for row := h.y1; row <= h.y2; row++ {
		for col := h.x1; col < h.x2; col++ {
			s.SetContent(col, row, '·', nil, global.MapStyle[global.MapFloor])
		}
	}
	//横墙
	for col := h.x1; col <= h.x2; col++ {
		s.SetContent(col, h.y1, tcell.RuneHLine, nil, h.style)
		s.SetContent(col, h.y2, tcell.RuneHLine, nil, h.style)
	}
	//竖墙
	for row := h.y1 + 1; row < h.y2; row++ {
		s.SetContent(h.x1, row, tcell.RuneVLine, nil, h.style)
		s.SetContent(h.x2, row, tcell.RuneVLine, nil, h.style)
	}
	//墙角
	if h.y1 != h.y2 && h.x1 != h.x2 {
		s.SetContent(h.x1, h.y1, tcell.RuneULCorner, nil, h.style)
		s.SetContent(h.x2, h.y1, tcell.RuneURCorner, nil, h.style)
		s.SetContent(h.x1, h.y2, tcell.RuneLLCorner, nil, h.style)
		s.SetContent(h.x2, h.y2, tcell.RuneLRCorner, nil, h.style)
	}

	//画门
	for _, p := range h.door {
		s.SetContent(p.x, p.y, '+', nil, h.style)
	}
}
