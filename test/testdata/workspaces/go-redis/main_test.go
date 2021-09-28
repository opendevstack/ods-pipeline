package main

import "testing"

func TestExample(t *testing.T) {
	err := example()
	if err != nil {
		t.Fatal(err)
	}
}
