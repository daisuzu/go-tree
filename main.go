package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	flaga := flag.Bool("a", false, "All files are listed.")
	flagL := flag.Int("L", 0, "Descend only level directories deep.")
	flagJ := flag.Bool("J", false, "Prints out an JSON representation of the tree.")
	flagV := flag.Bool("V", false, "Prints out the tree for tree.vim.")
	flag.Parse()

	var opts []Option
	if *flaga {
		opts = append(opts, WithAllFiles())
	}
	if *flagL > 0 {
		opts = append(opts, WithLevel(*flagL))
	}
	if *flagJ {
		opts = append(opts, WithJSONOutputter())
	}
	if *flagV {
		opts = append(opts, WithVimOutputter())
	}

	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	if err := BuildTree(os.Stdout, dirs, opts...); err != nil {
		log.Println(err)
	}
}
