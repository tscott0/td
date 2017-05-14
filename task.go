package main

import (
	"time"
)

type task struct {
	Name     string
	Desc     string
	URL      string
	Created  time.Time
	Deadline time.Time
}

func (t *task) metaString() string {
	var s string
	s += "Added: " + t.Created.Format("Mon Jan _2 15:04:05 2006") + "\n"
	s += "Due:   " + t.Deadline.Format("Mon Jan _2 15:04:05 2006") + "\n"
	s += "URL:   " + t.URL

	return s
}
