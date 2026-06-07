package main

import (
	"testing"

	"github.com/Carnoustie/theinfinitelibrary-backend/algorithms"
)


func TestCode(t *testing.T){

	s1 := "Jane Eyre";
	s2 := "Jane Air"

	dist:=algorithms.LevenshteinDistance(s1, s2, len(s1), len(s2))
	if(dist!=3){
		t.Error("\n\nEditing distance test failed")
	}
}
