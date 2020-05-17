package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	bak1, err := ioutil.ReadFile("testdata/testroot/rootsub1/blah.txt")
	if err != nil {
		panic(err)
	}
	bak2, err := ioutil.ReadFile("testdata/testroot/bleep.txt")
	if err != nil {
		panic(err)
	}

	var targetPath1 = "testdata/testroot/rootsub1/blah.txt"
	var targetPath2 = "testdata/testroot/bleep.txt"
	var replacePath1 = "testdata/replaceRoot/blah.txt"
	var replacePath2 = "testdata/replaceRoot/another/bleep.txt"

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
