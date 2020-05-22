package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	bak1         []byte
	bak2         []byte
	targetPath1  string
	targetPath2  string
	replacePath1 string
	replacePath2 string
	bakdir       string
)

func cleanup() {
	_ = os.Remove(targetPath1)
	_ = os.Remove(targetPath2)
	_ = os.Remove(bakdir)
	_ = ioutil.WriteFile(targetPath1, bak1, 0644)
	_ = ioutil.WriteFile(targetPath2, bak2, 0644)
}

func TestBackup(t *testing.T) {

	// backuproot = "testdata/testbackups"
	targetPath1 = "testdata/testroot/rootsub1/blah.txt"
	targetPath2 = "testdata/testroot/bleep.txt"
	replacePath1 = "testdata/replaceRoot/blah.txt"
	replacePath2 = "testdata/replaceRoot/another/bleep.txt"

	targetRoot = "testdata/testroot"
	replacementRoot = "testdata/replaceRoot"

	abs, err := filepath.Abs("testdata/testbackups")
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
	bak1, err = ioutil.ReadFile("testdata/testroot/rootsub1/blah.txt")
	if err != nil {
		panic(err)
	}
	bak2, err = ioutil.ReadFile("testdata/testroot/bleep.txt")
	if err != nil {
		panic(err)
	}

	targetPath1 = "testdata/testroot/rootsub1/blah.txt"
	targetPath2 = "testdata/testroot/bleep.txt"
	replacePath1 = "testdata/replaceRoot/blah.txt"
	replacePath2 = "testdata/replaceRoot/another/bleep.txt"

	// size1 := len(bak1)

	// size2 := len(bak2)

	cleanup := func() {
		_ = os.Remove(targetPath1)
		_ = os.Remove(targetPath2)
		_ = ioutil.WriteFile(targetPath1, bak1, 0644)
		_ = ioutil.WriteFile(targetPath2, bak2, 0644)
	}
	targetRoot = "testdata/testroot"
	replacementRoot = "testdata/replaceRoot"

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
