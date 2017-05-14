package main

import (
	"fmt"
	"log"
	"strings"
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
	tc int
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func redrawList(g *gocui.Gui) error {
	l, err := g.View("list")
	if err != nil {
		return err
	}

	l.Clear()
	fmt.Fprintln(l, tl.listString())

	return nil
}

func showTask(g *gocui.Gui, i int) error {
	if i < 0 || i > tl.maxCursor() {
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

func newDialog(g *gocui.Gui, v *gocui.View) error {

	maxX, maxY := g.Size()
	if v, err := g.SetView("new", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true

		if _, err := g.SetCurrentView("new"); err != nil {
			return err
		}
	}
	return nil
}

func newTask(g *gocui.Gui, v *gocui.View) error {
	defer closeNewDialog(g, v)

	newName := strings.TrimSpace(v.Buffer())

	if newName == "" {
		return nil
	}

	// TODO: Default values for new tasks
	t := task{
		Name:     newName,
		Desc:     "",
		URL:      "",
		Created:  time.Now(),
		Deadline: time.Now(),
	}

	tl.insertTask(g, tc, t)
	redrawList(g)
	showTask(g, tc)

	return nil
}

func closeNewDialog(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("new"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("list"); err != nil {
		return err
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy < tl.maxCursor() {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}

			tc = cy + 1

			if err := showTask(g, tc); err != nil {
				return err
			}
		}
	}

	return nil
}

func taskDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy < tl.maxCursor() {
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

			if err := redrawList(g); err != nil {
				return err
			}

			tc = cy + 1

			if err := showTask(g, tc); err != nil {
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

			tc = cy - 1

			if err := showTask(g, tc); err != nil {
				return err
			}
		}
	}
	return nil
}

func taskUp(g *gocui.Gui, v *gocui.View) error {
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

			if err := redrawList(g); err != nil {
				return err
			}

			tc = cy - 1

			if err := showTask(g, tc); err != nil {
				return err
			}
		}
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	// Global
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// Task List
	if err := g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowRight, gocui.ModNone, taskDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowLeft, gocui.ModNone, taskUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, editDesc); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyCtrlN, gocui.ModNone, newDialog); err != nil {
		return err
	}

	// Description
	if err := g.SetKeybinding("desc", gocui.KeyCtrlS, gocui.ModNone, saveDesc); err != nil {
		return err
	}
	if err := g.SetKeybinding("desc", gocui.KeyEsc, gocui.ModNone, cancelDesc); err != nil {
		return err
	}

	// New task dialog
	if err := g.SetKeybinding("new", gocui.KeyEnter, gocui.ModNone, newTask); err != nil {
		return err
	}
	if err := g.SetKeybinding("new", gocui.KeyEsc, gocui.ModNone, closeNewDialog); err != nil {
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
		v.BgColor = gocui.ColorDefault

		if len(tl.Task) > 0 {
			fmt.Fprintln(v, tl.listString())
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
	if v, err := g.SetView("desc", -1, maxY-descHeight-metaHeight, maxX, maxY-metaHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
		//v.Title = "Description"
		//v.BgColor = gocui.AttrReverse
		v.BgColor = gocui.ColorDefault

		if len(tl.Task) > 0 {
			t := tl.Task[0]
			fmt.Fprintln(v, t.Desc)
		}
	}
	if v, err := g.SetView("meta", -1, maxY-metaHeight, maxX, maxY); err != nil {
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
		return
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.InputEsc = true
	g.BgColor = gocui.ColorDefault

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
