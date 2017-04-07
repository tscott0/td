package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jroimartin/gocui"
)

const (
	descHeight = 10
	metaHeight = 4
)

var (
	tl taskList
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func maxCursor() int {
	return len(tl.Task) - 1
}

func (t *task) metaString() string {
	var s string
	s += "Added: " + t.Created.Format("Mon Jan _2 15:04:05 2006") + "\n"
	s += "Due:   " + t.Deadline.Format("Mon Jan _2 15:04:05 2006") + "\n"
	s += "URL:   " + t.URL

	return s
}
func listString() string {
	var s string
	for i, t := range tl.Task {
		s += "[" + strconv.Itoa(i+1) + "] " + t.Name + "\n"
	}
	return s
}

func drawList(g *gocui.Gui) error {
	l, err := g.View("list")
	if err != nil {
		return err
	}

	l.Clear()
	fmt.Fprintln(l, listString())

	return nil
}

func showItem(g *gocui.Gui, i int) error {
	if i < 0 || i > maxCursor() {
		return gocui.ErrQuit
	}

	t := tl.Task[i]

	// Update description
	d, err := g.View("desc")
	if err != nil {
		return err
	}

	d.Clear()
	fmt.Fprintf(d, "%s", t.Desc)

	// Update Meta
	m, err := g.View("meta")
	if err != nil {
		return err
	}

	m.Clear()
	fmt.Fprintln(m, t.metaString())

	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy < maxCursor() {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
			if err := showItem(g, cy+1); err != nil {
				return err
			}
		}
	}

	return nil
}

func listDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy < maxCursor() {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}

			// Swap positions of a and b in array
			a := tl.Task[cy]
			b := tl.Task[cy+1]
			tl.Task[cy+1] = a
			tl.Task[cy] = b

			if err := drawList(g); err != nil {
				return err
			}

			if err := showItem(g, cy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if cy > 0 {
			if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
				if err := v.SetOrigin(ox, oy-1); err != nil {
					return err
				}
			}
			if err := showItem(g, cy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func listUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if cy > 0 {
			if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
				if err := v.SetOrigin(ox, oy-1); err != nil {
					return err
				}
			}

			// Swap positions of a and b in array
			a := tl.Task[cy]
			b := tl.Task[cy-1]
			tl.Task[cy-1] = a
			tl.Task[cy] = b

			if err := drawList(g); err != nil {
				return err
			}

			if err := showItem(g, cy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func editDesc(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("desc"); err != nil {
		return err
	}
	return nil
}

func descFinish(g *gocui.Gui, v *gocui.View) error {
	l, err := g.View("list")
	if err != nil {
		return err
	}
	_, cy := l.Cursor()

	tl.Task[cy].Desc = v.Buffer()

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(tl); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("config.toml", buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	if _, err := g.SetCurrentView("list"); err != nil {
		return err
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowRight, gocui.ModNone, listDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowLeft, gocui.ModNone, listUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, editDesc); err != nil {
		return err
	}

	if err := g.SetKeybinding("desc", gocui.KeyCtrlSpace, gocui.ModNone, descFinish); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("list", -1, -1, maxX, maxY-descHeight-metaHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		if len(tl.Task) > 0 {
			fmt.Fprintln(v, listString())
			if err := v.SetCursor(1, 0); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}

		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("desc", 0, maxY-descHeight-metaHeight, maxX-1, maxY-metaHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
		v.Title = "Description"

		if len(tl.Task) > 0 {
			t := tl.Task[0]
			fmt.Fprintln(v, t.Desc)
		}
	}
	if v, err := g.SetView("meta", 0, maxY-metaHeight, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if len(tl.Task) > 0 {
			t := tl.Task[0]
			fmt.Fprintln(v, t.metaString())
		}
	}
	return nil
}

func main() {

	if _, err := toml.DecodeFile("config.toml", &tl); err != nil {
		//fmt.Println(err)
		return
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
