// package filereplacer walks a dir recursively, matching files, and then
// replacing those files.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
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
	backuproot      string
	includeExt      bool
	usage           = `filereplacer usage:
filereplacer [option] [target directory] [replacements root]

Where [target directory] is the root directory to search recursively for files
matching filenames found by recursively searching [replacement root].

Options:
  --backup, -b 		path to save backup of original files before replacement
  			When this flag is given, it must include a path, other
Example:
	filereplacer /tmp ~/tmp

This command would recursively search for filenames inside of ~/tmp, and if any
matching file names are found found (by searching /tmp recursively), the
original files are replaced by them.

	filereplacer -b ~/backups /tmp ~/tmp

Does the same as above, but saves the original to ~/backups/[original name].bak.
`
)

func parseArgs() (ok bool) {
	// check for help flags
	switch os.Args[1] {
	case "--help", "-h", "help", "-help":
		fmt.Printf("%s\n", usage)
		os.Exit(1)

	}
	pflag.StringVarP(&backuproot, "backup", "b", "", "backup files before replacement. They will have a '.bak' extension")
	pflag.BoolVarP(&includeExt, "ext", "i", true, "include file extension when comparing names")

	pflag.Usage = func() { fmt.Printf("%s", usage) }

	// this is fine. //
	///////////////////
	go pflag.Parse() //
	///////////////////

	// who will win the race?
	// ** PLACE YOUR BETS ** ///

	// set the paths from the args
	targetRoot = filepath.Clean(os.Args[1])
	replacementRoot = filepath.Clean(os.Args[2])
	backuproot = filepath.Clean(backuproot[:])

	return true
}

// file represents a fs file
type file struct {
	name string
	// path of containing dir
	path string

	// backup the file before overwriting?
	backup bool

	// the path of the backup
	backuppath string
}

// find file names to match against
func walkDirForFiles(root string, backup bool) ([]file, error) {
	// start with small cap, append will allocate more and copy if it needs
	// to.
	var files = make([]file, 0, 5)
	var walkfunc = func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			var f = file{name: info.Name(), path: filepath.Join(abs, info.Name())}
			if backup {
				f.backup = true
				f.backuppath = filepath.Clean(backuproot)
			}
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
func backup(f file) {

	b, err := os.Create(filepath.Join(f.backuppath, f.name+".bak"))
	if err != nil {
		panic(err)
	}

	old, err := os.Open(f.path)
	if err != nil {
		panic(err)
	}
	defer old.Close()

	_, err = io.Copy(b, old)
	if err != nil {
		panic(err)
	}
	err = b.Sync()
	if err != nil {
		panic(err)
	}

	err = b.Close()
	if err != nil {
		panic(err)
	}

}

// replace a single file, with another by copying over the bytes.
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

	// delete old file and seek to start for writing

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
	r, err := walkDirForFiles(replacementRoot, len(backuproot) > 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("searching for targets to replace...")
	t, err := walkDirForFiles(targetRoot, len(backuproot) > 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("replacing the files.")

	for i := 0; i < len(t); i++ {
		for j := len(r) - 1; j >= 0; j-- {
			if includeExt {
				if t[i].name[:len(filepath.Ext(t[i].name))] == r[j].name[:len(filepath.Ext(r[j].name))]
			}
			if t[i].name == r[j].name {
				if t[i].backup {
					backup(t[i])
				}
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
