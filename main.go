// package filereplacer walks a dir recursively, matching files, and then
// replacing those files.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func init() {

	// check arg lens
	if len(os.Args) < 3 {
		fmt.Println("wrong number of arguments!")
		// for spacing
		fmt.Println(" ")
		fmt.Printf("%s\n", usage)
		os.Exit(1)
	}

}

var (
	targetRoot      string
	replacementRoot string
	usage           = `filereplacer usage:
filereplacer [target directory] [replacements root]

Where [target directory] is the root directory to search recursively for files
matching filenames found by recursively searching [replacement root].

Example:
	filereplacer /tmp ~/tmp

This command would recursively search for filenames inside of ~/tmp, and if any
matching file names are found found (by searching /tmp recursively), the
original files are replaced by them.
`
)

func parseArgs() {
	// check for help flags
	switch os.Args[1] {
	case "--help", "-h", "help", "-help":
		fmt.Printf("%s\n", usage)
		os.Exit(1)

	}
	// set the paths from the args
	targetRoot = filepath.Clean(os.Args[1])
	replacementRoot = filepath.Clean(os.Args[2])
	return
}

// file represents a fs file
type file struct {
	name string
	// path of containing dir
	path string
}

// find file names to match against
func walkDirForFiles(root string) ([]file, error) {
	// start with small cap, append will allocate more and copy if it needs
	// to.
	var files = make([]file, 0, 5)
	var walkfunc = func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			var f = file{name: info.Name(), path: abs}
			files = append(files, f)
		}
		return nil
	}
	err := filepath.Walk(root, walkfunc)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// make backup of file in the same dir with '.bak' extention.
func backup(files []file) error {
	for _, f := range files {
		err := os.Rename(filepath.Join(f.path, f.name), filepath.Join(f.path, f.name+".bak"))
		if err != nil {
			return err
		}
	}
	return nil
}

// replace a single file, with another by copying the bytes.
func replace(f string, with string) error {
	r, err := os.Open(with)
	if err != nil {
		return err
	}
	defer r.Close()

	old, err := os.Create(f)
	if err != nil {
		return err
	}
	defer old.Close()

	inf, err := old.Stat()
	if err != nil {
		return fmt.Errorf("failed to get old filesize for comparison: %w\n", err)
	}
	oldSize := inf.Size()

	err = old.Truncate(0)
	if err != nil {
		return err
	}

	_, err = old.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error seeking file to start: %w\n", err)
	}

	n, err := io.Copy(old, r)
	if err != nil {
		return err
	}

	// make sure the new file is smaller
	if oldSize < n {
		fmt.Println("warning: the new file is larger than the one it replaced!")
	}

	return nil
}

func run() error {

	fmt.Println("searching for filenames of replacements...")
	r, err := walkDirForFiles(replacementRoot)
	if err != nil {
		panic(err)
	}

	fmt.Println("searching for targets to replace...")
	t, err := walkDirForFiles(targetRoot)
	if err != nil {
		panic(err)
	}

	fmt.Println("replacing the files.")

	for i := 0; i < len(t); i++ {
		for j := len(r) - 1; j >= 0; j-- {
			if t[i].name == r[j].name {
				err = replace(t[i].path, r[j].path)
				if err != nil {
					return err
				}
			}
		}
	}

	fmt.Println("done")
	return nil
}

func main() {

	parseArgs()
	if err := run(); err != nil {
		panic(err)
	}
	os.Exit(0)
}
