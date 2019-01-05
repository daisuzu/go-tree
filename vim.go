package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func newVimOutputter(w io.Writer, parent string, opt *option) outputter {
	return &vimOutputter{
		w:          w,
		opt:        opt,
		parentPath: filepath.Clean(parent),
	}
}

type vimOutputter struct {
	w            io.Writer
	opt          *option
	parentPath   string
	parentIndent string
	indent       string
	level        int
}

func (o *vimOutputter) Parent() string {
	return o.parentPath
}

func (o *vimOutputter) OutputError(msg string) error {
	// TODO: NotImplemented
	return nil
}

func (o *vimOutputter) OutputParent() error {
	parent := escape(o.parentPath)
	if o.level == 0 {
		header := fmt.Sprintf("%[1]s%[2]c..%[2]c\n%[1]s%[2]c.%[2]c", parent, filepath.Separator)
		if _, err := io.WriteString(o.w, header); err != nil {
			return err
		}
		return nil
	}

	if _, err := io.WriteString(o.w, "\n"+o.parentIndent+parent+string(filepath.Separator)); err != nil {
		return err
	}
	if !(o.opt.level > 0 && o.opt.level <= o.level) {
		if _, err := io.WriteString(o.w, "{{{"); err != nil {
			return err
		}
	}
	return nil
}

func (o *vimOutputter) Output(fi os.FileInfo, last bool) error {
	path := escape(filepath.Join(o.parentPath, fi.Name()))
	_, err := io.WriteString(o.w, "\n"+o.indent+path)
	return err
}

func (o *vimOutputter) Terminate() error {
	if o.level == 0 {
		_, err := io.WriteString(o.w, "\n")
		return err
	}

	_, err := io.WriteString(o.w, "}}}")
	return err
}

func (o *vimOutputter) ChildOutputter(path string, last bool) outputter {
	return &vimOutputter{
		w:            o.w,
		opt:          o.opt,
		parentPath:   filepath.Join(o.parentPath, path),
		parentIndent: o.indent,
		indent:       o.indent + "  ",
		level:        o.level + 1,
	}
}
