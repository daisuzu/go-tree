package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBuildTree(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	contents := filepath.Join(path, "contents")
	os.MkdirAll(filepath.Join(contents, "a", "a2"), 0700)
	os.Create(filepath.Join(contents, "a", "a2", "aa"))
	os.Create(filepath.Join(contents, "a", "a1"))
	os.Create(filepath.Join(contents, "b"))
	os.MkdirAll(filepath.Join(contents, "a", "a3"), 0700)
	os.MkdirAll(filepath.Join(contents, "c 0"), 0700)
	os.Create(filepath.Join(contents, "c 0", "c c"))

	notExists := filepath.Join(path, "not_exists")
	dirs := []string{contents, notExists}

	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "default",
			args: args{},
			want: strings.Join([]string{
				contents,
				"├── a",
				"│   ├── a1",
				"│   ├── a2",
				"│   │   └── aa",
				"│   └── a3",
				"├── b",
				"└── c\\ 0",
				"    └── c\\ c",
				notExists + " [error opening dir]",
				"",
			}, "\n"),
		},
		{
			name: "default L1",
			args: args{opts: []Option{WithLevel(1)}},
			want: strings.Join([]string{
				contents,
				"├── a",
				"├── b",
				"└── c\\ 0",
				notExists + " [error opening dir]",
				"",
			}, "\n"),
		},
		{
			name: "json",
			args: args{opts: []Option{WithJSONOutputter()}},
			want: strings.Join([]string{
				`[`,
				`  {"type":"directory","name":"` + contents + `","contents":[`,
				`    {"type":"directory","name":"a","contents":[`,
				`      {"type":"file","name":"a1"},`,
				`      {"type":"directory","name":"a2","contents":[`,
				`        {"type":"file","name":"aa"}`,
				`      ]},`,
				`      {"type":"directory","name":"a3","contents":[`,
				`      ]}`,
				`    ]},`,
				`    {"type":"file","name":"b"},`,
				`    {"type":"directory","name":"c 0","contents":[`,
				`      {"type":"file","name":"c c"}`,
				`    ]}`,
				`  ]}`,
				`]`,
				`[`,
				`  {"type":"directory","name":"` + notExists + `","contents":[{"error": "opening dir"}`,
				`  ]}`,
				`]`,
				"",
			}, "\n"),
		},
		{
			name: "json L1",
			args: args{opts: []Option{WithJSONOutputter(), WithLevel(1)}},
			want: strings.Join([]string{
				`[`,
				`  {"type":"directory","name":"` + contents + `","contents":[`,
				`    {"type":"directory","name":"a","contents":[`,
				`    ]},`,
				`    {"type":"file","name":"b"},`,
				`    {"type":"directory","name":"c 0","contents":[`,
				`    ]}`,
				`  ]}`,
				`]`,
				`[`,
				`  {"type":"directory","name":"` + notExists + `","contents":[{"error": "opening dir"}`,
				`  ]}`,
				`]`,
				"",
			}, "\n"),
		},
		{
			name: "vim",
			args: args{opts: []Option{WithVimOutputter()}},
			want: strings.Join([]string{
				fmt.Sprintf("%s%[2]c..%[2]c", contents, filepath.Separator),
				fmt.Sprintf("%s%[2]c.%[2]c", contents, filepath.Separator),
				fmt.Sprintf("%s%[2]ca%[2]c{{{", contents, filepath.Separator),
				fmt.Sprintf("  %s%[2]ca%[2]ca1", contents, filepath.Separator),
				fmt.Sprintf("  %s%[2]ca%[2]ca2%[2]c{{{", contents, filepath.Separator),
				fmt.Sprintf("    %s%[2]ca%[2]ca2%[2]caa}}}", contents, filepath.Separator),
				fmt.Sprintf("  %s%[2]ca%[2]ca3%[2]c{{{}}}}}}", contents, filepath.Separator),
				fmt.Sprintf("%s%[2]cb", contents, filepath.Separator),
				fmt.Sprintf("%s%[2]cc\\ 0%[2]c{{{", contents, filepath.Separator),
				fmt.Sprintf("  %s%[2]cc\\ 0%[2]cc\\ c}}}", contents, filepath.Separator),
				"",
			}, "\n"),
		},
		{
			name: "vim L1",
			args: args{opts: []Option{WithVimOutputter(), WithLevel(1)}},
			want: strings.Join([]string{
				fmt.Sprintf("%s%[2]c..%[2]c", contents, filepath.Separator),
				fmt.Sprintf("%s%[2]c.%[2]c", contents, filepath.Separator),
				fmt.Sprintf("%s%[2]ca%[2]c", contents, filepath.Separator),
				fmt.Sprintf("%s%[2]cb", contents, filepath.Separator),
				fmt.Sprintf("%s%[2]cc\\ 0%[2]c", contents, filepath.Separator),
				"",
			}, "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := BuildTree(w, dirs, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("BuildTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := w.String(); got != tt.want {
				t.Errorf("BuildTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getContents(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	os.Create(filepath.Join(path, ".hidden"))
	os.Create(filepath.Join(path, "normal"))
	os.Create(filepath.Join(path, "ok"))
	os.Create(filepath.Join(path, "more"))

	type args struct {
		opt *option
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "default",
			args: args{opt: &option{sort: defaultSort}},
			want: []string{"more", "normal", "ok"},
		},
		{
			name: "with all",
			args: args{opt: &option{all: true, sort: defaultSort}},
			want: []string{".hidden", "more", "normal", "ok"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contents, err := getContents(path, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("getContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := make([]string, len(contents))
			for i := range contents {
				got[i] = contents[i].Name()
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getContents() = %v, want %v", got, tt.want)
			}
		})
	}
}
