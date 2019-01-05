package main

import (
	"os"
	"sort"
)

// Option sets options such as listing and output format, etc.
type Option func(*option)

// WithAllFiles returns an Option to list all files.
func WithAllFiles() Option {
	return func(o *option) {
		o.all = true
	}
}

// WithLevel returns an Option which set max depth of the directory tree.
func WithLevel(i int) Option {
	return func(o *option) {
		o.level = i
	}
}

// WithJSONOutputter returns an Option to output as an JSON formatted array.
func WithJSONOutputter() Option {
	return func(o *option) {
		o.outputter = newJSONOutputter
	}
}

// WithVimOutputter returns an Option to output as text formatted for tree.vim.
func WithVimOutputter() Option {
	return func(o *option) {
		o.outputter = newVimOutputter
	}
}

func defaultSort(contents []os.FileInfo) {
	sort.Slice(contents, func(i, j int) bool { return contents[i].Name() < contents[j].Name() })
}

type option struct {
	all       bool
	level     int
	outputter outputterFunc
	sort      func([]os.FileInfo)
}
