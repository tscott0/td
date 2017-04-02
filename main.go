package main

import (
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	ui "github.com/gizak/termui"
)

type task struct {
	Name     string
	Desc     string
	URL      string
	Created  time.Time
	Deadline time.Time
}

type taskList struct {
	Task []task
}

func main() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	tt := NewTidyTable()

	par2 := ui.NewPar(" q quit")
	par2.Height = 6
	par2.BorderLabel = "Controls"

	desc := ui.NewPar("This is a long\ndescription")
	desc.Height = 11
	desc.BorderLabel = "Description"

	url := ui.NewPar("A URL string")
	url.Height = 3
	url.BorderLabel = "URL"

	created := ui.NewPar("Mon Jan _2 15:04:05 2006")
	created.Height = 3
	created.BorderLabel = "Added"

	deadline := ui.NewPar("Mon Jan _2 15:04:05 2006")
	deadline.Height = 3
	deadline.BorderLabel = "Deadline"

	var tl taskList
	if _, err := toml.DecodeFile("config.toml", &tl); err != nil {
		//fmt.Println(err)
		return
	}

	//for _, t := range tl.Task {
	// TODO
	//}

	tt.Height = 3 + len(tl.Task)

	// build layout
	ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, tt)))
	ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, desc)))
	ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, url)))
	ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, created)))
	ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, deadline)))
	ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, par2)))

	// calculate layout
	ui.Body.Align()

	rows1 := make([][]string, 1+len(tl.Task))

	draw := func(z int) {
		i := 0
		for _, t := range tl.Task {
			rows1[i] = []string{t.Name}
			if tt.cursor == i {
				desc.Text = t.Desc
				url.Text = t.URL
				created.Text = t.Created.Format("Mon Jan _2 15:04:05 2006")
				deadline.Text = t.Deadline.Format("Mon Jan _2 15:04:05 2006")
			}
			if z != 0 {
				par2.Text = strconv.Itoa(z)
			}
			i++
		}
		tt.Rows = rows1
		tt.Highlight()
		ui.Render(ui.Body)
	}

	// TODO: Bit of a hack. Not sure why it doesn't draw immediately
	draw(0)
	//tt.BgColors[0] = ui.ColorGreen

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		tt.Up()
		draw(0)
	})

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		tt.Down()
		draw(0)
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		draw(int(t.Count))
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()

}
