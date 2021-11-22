// Check tries to implement the output `git check-ignore`

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"nogo"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// DirFs actually implements StatFS, so we can use it.
	wdfs := os.DirFS(wd).(fs.StatFS)

	n := nogo.NewGitignore(nogo.WithMatchParents(), nogo.WithFS(wdfs))
	if err := n.AddAll(); err != nil {
		panic(err)
	}

	files := flag.Args()
	for _, toSearch := range files {
		toSearch = strings.TrimPrefix(toSearch, "./")
		if toSearch == "" {
			toSearch = "."
		}

		info, err := wdfs.Stat(toSearch)
		if err != nil {
			panic(err)
		}

		if info.Name() == ".git" {
			return
		}

		if err != nil {
			panic(err)
		}

		if n.MatchPath(toSearch).Resolve(info.IsDir()) {
			fmt.Printf("./%v\n", toSearch)
		}
	}
}
