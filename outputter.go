package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type outputter interface {
	Parent() string
	OutputError(msg string) error
	OutputParent() error
	Output(fi os.FileInfo, last bool) error
	Terminate() error
	ChildOutputter(path string, last bool) outputter
}

type outputterFunc func(w io.Writer, parent string, opt *option) outputter

func newDefaultOutputter(w io.Writer, parent string, opt *option) outputter {
	return &defaultOutputter{
		w:          w,
		opt:        opt,
		parentPath: parent,
		prefix: map[bool]string{
			false: "├── ",
			true:  "└── ",
		},
		root: true,
	}
}

type defaultOutputter struct {
	w            io.Writer
	opt          *option
	parentPath   string
	parentBranch string
	branch       string
	prefix       map[bool]string
	last         bool
	root         bool
}

func (o *defaultOutputter) Parent() string {
	return o.parentPath
}

func (o *defaultOutputter) parentForOut() string {
	parent := o.parentPath
	if !o.root {
		parent = filepath.Base(parent)
	}
	parent = escape(parent)
	return parent
}

func (o *defaultOutputter) OutputError(msg string) error {
	_, err := io.WriteString(o.w, o.parentBranch+o.parentForOut()+" [error "+msg+"]\n")
	return err
}

func (o *defaultOutputter) OutputParent() error {
	_, err := io.WriteString(o.w, o.parentBranch+o.parentForOut()+"\n")
	return err
}

func (o *defaultOutputter) Output(fi os.FileInfo, last bool) error {
	name := escape(fi.Name())
	_, err := io.WriteString(o.w, o.branch+o.prefix[last]+name+"\n")
	return err
}

func (o *defaultOutputter) Terminate() error {
	return nil
}

func (o *defaultOutputter) ChildOutputter(path string, last bool) outputter {
	branch := o.branch
	if last {
		branch += "    "
	} else {
		branch += "│   "
	}

	return &defaultOutputter{
		w:            o.w,
		opt:          o.opt,
		parentPath:   filepath.Join(o.parentPath, path),
		parentBranch: o.branch + o.prefix[last],
		branch:       branch,
		prefix:       o.prefix,
		last:         last,
	}
}

func escape(s string) string {
	return strings.Replace(s, " ", "\\ ", -1)
}
