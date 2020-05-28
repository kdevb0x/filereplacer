package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)


func cleanup(f testcase) {
	os.Remove(f.origpath())
	ioutil.WriteFile(f.origpath, f.saveorig, 0644)
}

type testcase struct {
	origpath string
	replacementpath string
	backuppath string
	size int64
	saveorig []byte
	cleanup func(testcase)

}
const (
	targetRoot = "testdata/testroot"
	replacementRoot = "testdata/replaceRoot"

	blahPath = "testdata/testroot/rootsub1/blah.txt"
	blahRePath = "testdata/replaceRoot/blah.txt"

	bleepPath = "testdata/testroot/bleep.txt"
	bleepRePath = "testdata/replaceRoot/another/bleep.txt"

	backupDir = "testdata/backupdir"

)

func testcases(paths [][2]string) []testcase {
	var c = make([]testcase, count)
	for i := 0; i < len(paths); i++ {
		var t testcase
		t.origpath = paths[i][0]
		t.replacementpath = paths[i][1]
		t.backuppath = backupDir
		t.size = filesize(t.origpath)
		if t.size == -1 {
			info, err := os.Stat(t.origpath)
			if err != nil {
				panic("failed to get filesize for comparison")
			}
			t.size = info.Size()
		}


		save, err := ioutil.ReadFile(t.origpath)
		if err != nil {
			panic(err)
		}
		t.saveorig = save
		t.cleanup = cleanup(t)

		c = append(c, t)
	}
	return c


}

func filesize(f string) int64 {
	d, _ := os.Open(f)
	if err != nil {
		return -1
	}
	defer d.Close()
	s, err := d.Stat()
	if err != nil {
		return -1
	}
	return s.Size()
}
func TestBackup(t *testing.T) {

	// backuproot = "testdata/testbackups"
	/*
	targetPath1 = "testdata/testroot/rootsub1/blah.txt"
	targetPath2 = "testdata/testroot/bleep.txt"
	replacePath1 = "testdata/replaceRoot/blah.txt"
	replacePath2 = "testdata/replaceRoot/another/bleep.txt"

	targetRoot = "testdata/testroot"
	replacementRoot = "testdata/replaceRoot"
	*/

	for _, c := range tests {
	abs, err := filepath.Abs(c)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
	}
	bakdir = abs
	backuproot = abs
	bleeppath := filepath.Join(backuproot, "bleep.txt.bak")
	err = run()
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Open(bleeppath); err != nil && os.IsNotExist(err) {
		t.Fail()
	}
	if err := os.Remove(bleeppath); err != nil {
		t.Fatalf("failed to remove backup file; err: %s\n", err.Error())
	}
	t.Cleanup(func() {
		cleanup()
		err := ioutil.WriteFile(bakdir, nil, os.ModeDir)
		if err != nil {
			t.Logf("failed to recreate %s: %s\n", bakdir, err.Error())
		}
	})
}

func TestRun(t *testing.T) {
	var err error
	var bak1, bak2 []byte
	bak1, err = ioutil.ReadFile("testdata/testroot/rootsub1/blah.txt")
	if err != nil {
		panic(err)
	}
	bak2, err = ioutil.ReadFile("testdata/testroot/bleep.txt")
	if err != nil {
		panic(err)
	}



	// run the test
	err = run()
	if err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		t.Fail()
	}

	t1, err := os.Stat(targetPath1)
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}
	t2, err := os.Stat(replacePath1)
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}
	if t1.Size() != t2.Size() {
		fmt.Printf("bad size: %d != %d\n", t1.Size(), t2.Size())
		t.Fail()
	}

	v1, err := os.Stat(targetPath2)
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}
	v2, err := os.Stat(replacePath2)
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}
	if v1.Size() != v2.Size() {
		fmt.Printf("bad size: %d != %d\n", v1.Size(), v2.Size())
		t.Fail()
	}
	t.Cleanup(cleanup)
}
