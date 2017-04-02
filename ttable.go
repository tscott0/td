package main

import ui "github.com/gizak/termui"

type ttable struct {
	*ui.Table
	cursor int
}

func NewTidyTable() *ttable {

	tt := &ttable{
		Table:  ui.NewTable(),
		cursor: 0,
	}

	tt.Separator = false
	//tt.Border = false
	//tt.PaddingLeft = -1
	//tt.PaddingTop = 1
	//tt.PaddingBottom = 1

	return tt
}

func (t *ttable) Highlight() {

	for i, _ := range t.BgColors {
		if i == t.cursor {
			t.BgColors[i] = ui.ColorGreen
		} else {
			t.BgColors[i] = ui.ColorDefault
		}
	}
}

func (t *ttable) Up() {
	if t.cursor > 0 {
		t.cursor -= 1
	}
	t.Highlight()
}

func (t *ttable) Down() {
	maxCursor := len(t.Rows) - 2

	if t.cursor < maxCursor {
		t.cursor += 1
	}
	t.Highlight()
}
