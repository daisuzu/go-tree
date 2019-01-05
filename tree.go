package main

import (
	"io"
	"os"
)

// BuildTree builds trees with the given dirs as their root into w.
func BuildTree(w io.Writer, dirs []string, opts ...Option) error {
	opt := option{
		outputter: newDefaultOutputter,
		sort:      defaultSort,
	}
	for _, o := range opts {
		o(&opt)
	}

	for _, root := range dirs {
		if err := build(opt.outputter(w, root, &opt), &opt, 0); err != nil {
			return err
		}
	}
	return nil
}

func build(o outputter, opt *option, level int) error {
	contents, err := getContents(o.Parent(), opt)
	if err != nil {
		err = o.OutputError("opening dir")
		return err
	}
	if err := o.OutputParent(); err != nil {
		return err
	}

	if opt.level > 0 && opt.level <= level {
		return nil
	}

	for i, c := range contents {
		last := i == len(contents)-1
		if c.IsDir() {
			err = build(o.ChildOutputter(c.Name(), last), opt, level+1)
		} else {
			err = o.Output(c, last)
		}
		if err != nil {
			return err
		}
	}

	if err := o.Terminate(); err != nil {
		return err
	}

	return nil
}

func getContents(path string, opt *option) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	contents := make([]os.FileInfo, 0, len(list))
	for _, info := range list {
		if !opt.all && info.Name()[0] == '.' {
			continue
		}
		contents = append(contents, info)
	}

	opt.sort(contents)

	return contents, nil
}
