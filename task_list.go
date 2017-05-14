package main

import (
	"strconv"

	"github.com/jroimartin/gocui"
)

type taskList struct {
	Task []task
}

func (l *taskList) insertTask(g *gocui.Gui, index int, value task) error {

	new := append(l.Task[:index], append([]task{value}, l.Task[index:]...)...)
	l.Task = new

	return nil
}

func (l *taskList) listString() string {
	var s string
	for i, t := range l.Task {
		s += "[" + strconv.Itoa(i+1) + "] " + t.Name + "\n"
	}
	return s
}

func (l *taskList) maxCursor() int {
	return len(l.Task) - 1
}
