package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jroimartin/gocui"
)

func editDesc(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("desc"); err != nil {
		return err
	}
	return nil
}

func saveDesc(g *gocui.Gui, v *gocui.View) error {
	tl.Task[tc].Desc = strings.TrimSpace(v.Buffer())

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

func cancelDesc(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("list"); err != nil {
		return err
	}
	showTask(g, tc)
	return nil
}
