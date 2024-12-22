package main

import "time"

var Writings []Writing

type Writing struct {
	Title       string
	Slug        string
	Featured    bool
	Body        string
	PublishedAt time.Time
}
