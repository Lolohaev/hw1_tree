package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeWithPrefix(out, path, printFiles, "")
}

func dirTreeWithPrefix(out io.Writer, path string, printFiles bool, prefix string) error {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if !printFiles {
		dirEntries := make([]fs.DirEntry, 0)
		for _, e := range entries {
			if e.IsDir() {
				dirEntries = append(dirEntries, e)
			}
		}
		entries = dirEntries
	}

	sort.Sort(ByName(entries))

	prefixLast := prefix + "├───"
	prefixNext := prefix
	for i, e := range entries {
		if i == len(entries)-1 {
			prefixLast = prefix + "└───"
			prefixNext = prefix + "\t"
		} else {
			prefixNext = prefix + "│\t"
		}
		endOfLine := "\n"

		if !e.IsDir() {
			info, err := e.Info()
			if err != nil {
				return err
			}
			size := "empty"
			if info.Size() != 0 {
				size = fmt.Sprint(info.Size()) + "b"
			}
			endOfLine = " (" + size + ")\n"
		}

		out.Write([]byte(prefixLast + e.Name() + endOfLine))
		//fmt.Println(prefixLast + e.Name())
		if e.IsDir() {
			dirTreeWithPrefix(out, path+string(os.PathSeparator)+e.Name(), printFiles, prefixNext)
		}
	}
	return nil
}

type ByName []fs.DirEntry

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }
