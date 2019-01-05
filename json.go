package main

import (
	"io"
	"os"
	"path/filepath"
)

func newJSONOutputter(w io.Writer, parent string, opt *option) outputter {
	return &jsonOutputter{
		w:          w,
		opt:        opt,
		parentPath: parent,
		indent:     "  ",
		last:       true,
		root:       true,
	}
}

type jsonOutputter struct {
	w          io.Writer
	opt        *option
	parentPath string
	indent     string
	last       bool
	root       bool
	level      int
}

func (o *jsonOutputter) Parent() string {
	return o.parentPath
}

func (o *jsonOutputter) parentForOut() string {
	parent := o.parentPath
	if !o.root {
		parent = filepath.Base(parent)
	}
	return parent
}

func (o *jsonOutputter) OutputError(msg string) error {
	if o.root {
		if _, err := io.WriteString(o.w, "[\n"); err != nil {
			return err
		}
	}

	parent := o.parentForOut()
	if _, err := io.WriteString(o.w, o.indent+`{"type":"directory","name":"`+parent+`","contents":[{"error": "`+msg+`"}`+"\n"); err != nil {
		return err
	}
	if err := o.Terminate(); err != nil {
		return err
	}

	return nil
}

func (o *jsonOutputter) OutputParent() error {
	if o.root {
		if _, err := io.WriteString(o.w, "[\n"); err != nil {
			return err
		}
	}

	parent := o.parentForOut()
	if _, err := io.WriteString(o.w, o.indent+`{"type":"directory","name":"`+parent+`","contents":[`+"\n"); err != nil {
		return err
	}
	if o.opt.level > 0 && o.opt.level <= o.level {
		if err := o.Terminate(); err != nil {
			return err
		}
	}

	return nil
}

func (o *jsonOutputter) Output(fi os.FileInfo, last bool) error {
	if _, err := io.WriteString(o.w, o.indent+"  "+`{"type":"file","name":"`+fi.Name()+`"}`); err != nil {
		return err
	}
	if !last {
		if _, err := io.WriteString(o.w, ","); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(o.w, "\n"); err != nil {
		return err
	}

	return nil
}

func (o *jsonOutputter) Terminate() error {
	if _, err := io.WriteString(o.w, o.indent+`]}`); err != nil {
		return err
	}
	if !o.last {
		if _, err := io.WriteString(o.w, ","); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(o.w, "\n"); err != nil {
		return err
	}

	if o.root {
		if _, err := io.WriteString(o.w, "]\n"); err != nil {
			return err
		}
	}

	return nil
}

func (o *jsonOutputter) ChildOutputter(path string, last bool) outputter {
	return &jsonOutputter{
		w:          o.w,
		opt:        o.opt,
		parentPath: filepath.Join(o.parentPath, path),
		indent:     o.indent + "  ",
		last:       last,
		level:      o.level + 1,
	}
}
